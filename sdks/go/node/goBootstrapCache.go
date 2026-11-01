package node

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/opctl/opctl/sdks/go/internal/uniquestring"
	"github.com/opctl/opctl/sdks/go/internal/unsudo"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/containerruntime"
	"github.com/opctl/opctl/sdks/go/node/pubsub"
	"github.com/opctl/opctl/sdks/go/opspec/goop"
)

// goBootstrapCache is shared by API describes and all op callers in this node.
type goBootstrapCache struct {
	// schema is immutable for the lifetime of a node. Zero keeps test fixtures at the production default.
	schema        goop.SchemaRange
	describes     map[string]describeMemo
	describeBytes int
	describeClock uint64
	runtime       containerruntime.ContainerRuntime
	stateStore    stateStore
	pubSub        pubsub.PubSub
	once          sync.Once
	initErr       error
	directory     string
	mu            sync.Mutex
	locks         map[string]*bootstrapLock
	leases        map[string]int
	budget        int64
}
type bootstrapLock struct {
	token chan struct{}
	users int
}
type goProgram struct {
	key, path string
	source    *goop.Source
	image     *containerruntime.ImageIdentity
	metadata  programMetadata
	ephemeral bool
}
type programMetadata struct {
	Recipe, Key, ImageID, OS, Arch, Variant, GoVersion, Checksum string
	Size, LauncherSize                                           int64
	LauncherChecksum                                             string
	Created, LastUsed                                            time.Time
}

type discardBootstrapEvents struct{}

func (discardBootstrapEvents) Publish(model.Event) {}

func (cache *goBootstrapCache) initialize(dataDir string) error {
	cache.once.Do(func() {
		cache.directory = filepath.Join(dataDir, "go", "programs", goop.CacheFormatVersion)
		cache.initErr = unsudo.CreateDir(cache.directory)
		if cache.initErr != nil {
			return
		}
		entries, err := os.ReadDir(cache.directory)
		if err != nil {
			cache.initErr = err
			return
		}
		for _, entry := range entries {
			if strings.HasPrefix(entry.Name(), ".tmp-") {
				if err := os.RemoveAll(filepath.Join(cache.directory, entry.Name())); err != nil {
					cache.initErr = err
					return
				}
			}
		}
	})
	return cache.initErr
}
func (cache *goBootstrapCache) acquire(ctx context.Context, key string) (func(), bool, error) {
	if err := ctx.Err(); err != nil {
		return nil, false, err
	}
	cache.mu.Lock()
	if cache.locks == nil {
		cache.locks = map[string]*bootstrapLock{}
	}
	lock := cache.locks[key]
	if lock == nil {
		lock = &bootstrapLock{token: make(chan struct{}, 1)}
		cache.locks[key] = lock
	}
	lock.users++
	cache.mu.Unlock()
	drop := func() {
		cache.mu.Lock()
		defer cache.mu.Unlock()
		lock.users--
		if lock.users == 0 {
			delete(cache.locks, key)
		}
	}
	select {
	case lock.token <- struct{}{}:
		return func() { <-lock.token; drop() }, false, nil
	default:
	}
	select {
	case lock.token <- struct{}{}:
		return func() { <-lock.token; drop() }, true, nil
	case <-ctx.Done():
		drop()
		return nil, true, ctx.Err()
	}
}
func (cache *goBootstrapCache) limit() int64 {
	if cache.budget > 0 {
		return cache.budget
	}
	return 1 << 30
}
func (cache *goBootstrapCache) lease(program *goProgram) func() {
	cache.mu.Lock()
	if cache.leases == nil {
		cache.leases = map[string]int{}
	}
	cache.leases[program.path]++
	cache.mu.Unlock()
	return func() {
		cache.mu.Lock()
		defer cache.mu.Unlock()
		cache.leases[program.path]--
		if cache.leases[program.path] == 0 {
			delete(cache.leases, program.path)
			if program.ephemeral {
				os.RemoveAll(program.path)
			}
		}
	}
}
func executableChecksum(path string) (string, int64, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return "", 0, err
	}
	if !info.Mode().IsRegular() || info.Mode().Perm()&0111 == 0 {
		return "", 0, errInvalidProgram
	}
	file, err := os.Open(path)
	if err != nil {
		return "", 0, err
	}
	defer file.Close()
	sum := sha256.New()
	if _, err := io.Copy(sum, file); err != nil {
		return "", 0, err
	}
	return hex.EncodeToString(sum.Sum(nil)), info.Size(), nil
}

var errInvalidProgram = errors.New("invalid cached Go program")

func readProgram(path, key string, image *containerruntime.ImageIdentity) (programMetadata, error) {
	var metadata programMetadata
	info, err := os.Lstat(filepath.Join(path, "metadata.json"))
	if err != nil {
		return metadata, err
	}
	if !info.Mode().IsRegular() || info.Size() > 64*1024 {
		return metadata, errInvalidProgram
	}
	data, err := os.ReadFile(filepath.Join(path, "metadata.json"))
	if err != nil {
		return metadata, err
	}
	if json.Unmarshal(data, &metadata) != nil || metadata.Key != key || metadata.Recipe != goop.CacheFormatVersion || metadata.ImageID != image.ID || metadata.OS != image.OS || metadata.Arch != image.Arch || metadata.Variant != image.Variant {
		return metadata, errInvalidProgram
	}
	checksum, size, err := executableChecksum(filepath.Join(path, "program"))
	if err != nil {
		return metadata, err
	}
	launcherChecksum, launcherSize, launcherErr := executableChecksum(filepath.Join(path, "launcher"))
	if checksum != metadata.Checksum || size != metadata.Size || launcherErr != nil || launcherChecksum != metadata.LauncherChecksum || launcherSize != metadata.LauncherSize {
		return metadata, errInvalidProgram
	}
	return metadata, nil
}
func writeProgramMetadata(path string, metadata programMetadata) error {
	data, err := json.Marshal(metadata)
	if err != nil {
		return err
	}
	temporary := filepath.Join(path, ".metadata.tmp")
	if err := unsudo.CreateFile(temporary, data); err != nil {
		return err
	}
	if err := os.Chmod(temporary, 0600); err != nil {
		return err
	}
	return os.Rename(temporary, filepath.Join(path, "metadata.json"))
}
func (cache *goBootstrapCache) prune() error {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	entries, err := os.ReadDir(cache.directory)
	if err != nil {
		return err
	}
	type candidate struct {
		path string
		size int64
		used time.Time
	}
	candidates := []candidate{}
	var total int64
	for _, entry := range entries {
		if len(entry.Name()) != 64 || !entry.IsDir() {
			continue
		}
		path := filepath.Join(cache.directory, entry.Name())
		data, err := os.ReadFile(filepath.Join(path, "metadata.json"))
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return err
		}
		var metadata programMetadata
		if json.Unmarshal(data, &metadata) != nil {
			continue
		}
		total += metadata.Size + metadata.LauncherSize
		if cache.leases[path] == 0 {
			candidates = append(candidates, candidate{path, metadata.Size + metadata.LauncherSize, metadata.LastUsed})
		}
	}
	sort.Slice(candidates, func(i, j int) bool { return candidates[i].used.Before(candidates[j].used) })
	for _, candidate := range candidates {
		if total <= cache.limit() {
			break
		}
		if err := os.RemoveAll(candidate.path); err != nil {
			return err
		}
		total -= candidate.size
	}
	return nil
}

func (cache *goBootstrapCache) prepare(ctx context.Context, bootstrap _bootstrapper, op *model.OpCall, root string) (*goProgram, func(), error) {
	if err := ctx.Err(); err != nil {
		return nil, nil, err
	}
	if err := cache.initialize(bootstrap.dataDirPath); err != nil {
		return nil, nil, err
	}
	for attempt := 0; attempt < 2; attempt++ {
		source, err := goop.Discover(op.OpPath, op.RepositoryRoot)
		if err != nil {
			return nil, nil, &OpDeclarationError{Err: err}
		}
		stubs, err := bootstrap.stubs(ctx, op, root, source)
		if err != nil {
			return nil, nil, &OpDeclarationError{Err: err}
		}
		ref := "golang:" + source.GoVersion
		imageRequest := &model.ContainerCallImage{Ref: &ref}
		if cache.stateStore != nil {
			if auth := cache.stateStore.TryGetAuth(ref); auth != nil {
				imageRequest.PullCreds = &auth.Creds
			}
		}
		image, err := cache.runtime.ResolveImage(ctx, imageRequest, root, discardBootstrapEvents{})
		if err != nil {
			return nil, nil, err
		}
		toolchain := goop.Toolchain{ID: image.ID, OS: image.OS, Arch: image.Arch, Variant: image.Variant}
		buildSpec, auth, err := bootstrap.prepareBuild(op, source, image, stubs)
		if err != nil {
			return nil, nil, err
		}
		recipe := goop.BuildRecipe{Launcher: goBootstrapLauncher, Command: buildSpec.Cmd.([]interface{}), Env: buildSpec.EnvVars.(map[string]interface{}), Stubs: stubs.digests()}
		key, err := source.Fingerprint(ctx, bootstrapWritablePaths(bootstrap.dataDirPath, source.ModuleRoot), toolchain, recipe)
		if ctx.Err() != nil {
			return nil, nil, ctx.Err()
		}
		if err != nil {
			return nil, nil, &OpDeclarationError{Err: err}
		}
		unlock, waited, err := cache.acquire(ctx, key)
		if err != nil {
			return nil, nil, err
		}
		program, release, retry, err := cache.prepareLocked(ctx, bootstrap, op, root, source, image, key, toolchain, recipe, buildSpec, auth, stubs, waited)
		unlock()
		if err != nil {
			return nil, nil, err
		}
		if retry {
			continue
		}
		if err := cache.prune(); err != nil {
			release()
			return nil, nil, err
		}
		return program, release, nil
	}
	return nil, nil, &OpDeclarationError{Err: fmt.Errorf("Go op source changed repeatedly during bootstrap; retry with stable sources")}
}
func (cache *goBootstrapCache) prepareLocked(ctx context.Context, bootstrap _bootstrapper, op *model.OpCall, root string, source *goop.Source, image *containerruntime.ImageIdentity, key string, toolchain goop.Toolchain, recipe goop.BuildRecipe, buildSpec *model.ContainerCallSpec, auth bootstrapAuth, stubs goStubSet, waited bool) (*goProgram, func(), bool, error) {
	if waited {
		fresh, err := source.Fingerprint(ctx, bootstrapWritablePaths(bootstrap.dataDirPath, source.ModuleRoot), toolchain, recipe)
		if err != nil {
			return nil, nil, false, err
		}
		if fresh != key {
			return nil, nil, true, nil
		}
	}
	path := filepath.Join(cache.directory, key)
	guard := cache.lease(&goProgram{path: path})
	metadata, err := readProgram(path, key, image)
	if err == nil {
		program := &goProgram{key: key, path: path, source: source, image: image, metadata: metadata}
		release := guard
		metadata.LastUsed = time.Now().UTC()
		if err := writeProgramMetadata(path, metadata); err != nil {
			release()
			return nil, nil, false, err
		}
		return program, release, false, nil
	}
	guard()
	if !os.IsNotExist(err) && !errors.Is(err, errInvalidProgram) {
		return nil, nil, false, err
	}
	cache.mu.Lock()
	leased := cache.leases[path] > 0
	cache.mu.Unlock()
	if leased {
		return nil, nil, false, fmt.Errorf("cached Go program became corrupt during an active call")
	}
	if err := os.RemoveAll(path); err != nil {
		return nil, nil, false, err
	}
	id, err := uniquestring.Construct()
	if err != nil {
		return nil, nil, false, err
	}
	staging := filepath.Join(cache.directory, ".tmp-"+key+"-"+id)
	if err := unsudo.CreateDir(staging); err != nil {
		return nil, nil, false, err
	}
	keep := false
	defer func() {
		if !keep {
			os.RemoveAll(staging)
		}
	}()
	if err := bootstrap.buildProgram(ctx, op, root, source, buildSpec, auth, staging, stubs); err != nil {
		return nil, nil, false, err
	}
	if err := ctx.Err(); err != nil {
		return nil, nil, false, err
	}
	fresh, err := source.Fingerprint(ctx, bootstrapWritablePaths(bootstrap.dataDirPath, source.ModuleRoot), toolchain, recipe)
	if err != nil {
		return nil, nil, false, err
	}
	if fresh != key {
		return nil, nil, true, nil
	}
	checksum, size, err := executableChecksum(filepath.Join(staging, "program"))
	if err != nil {
		return nil, nil, false, fmt.Errorf("validate built Go program: %w", err)
	}
	launcherChecksum, launcherSize, err := executableChecksum(filepath.Join(staging, "launcher"))
	if err != nil {
		return nil, nil, false, fmt.Errorf("validate built Go launcher: %w", err)
	}
	version, err := os.ReadFile(filepath.Join(staging, "go-version"))
	if err != nil {
		return nil, nil, false, err
	}
	if len(version) > 1024 {
		return nil, nil, false, fmt.Errorf("invalid Go version diagnostic")
	}
	if err := os.Remove(filepath.Join(staging, "go-version")); err != nil {
		return nil, nil, false, err
	}
	now := time.Now().UTC()
	metadata = programMetadata{Recipe: goop.CacheFormatVersion, Key: key, ImageID: image.ID, OS: image.OS, Arch: image.Arch, Variant: image.Variant, GoVersion: strings.TrimSpace(string(version)), Checksum: checksum, Size: size, LauncherChecksum: launcherChecksum, LauncherSize: launcherSize, Created: now, LastUsed: now}
	if err := writeProgramMetadata(staging, metadata); err != nil {
		return nil, nil, false, err
	}
	program := &goProgram{key: key, path: staging, source: source, image: image, metadata: metadata, ephemeral: size+launcherSize > cache.limit()}
	if !program.ephemeral {
		program.path = path
	}
	release := cache.lease(program)
	if !program.ephemeral {
		if err := os.Rename(staging, path); err != nil {
			release()
			return nil, nil, false, err
		}
	}
	keep = true

	return program, release, false, nil
}

func bootstrapWritablePaths(dataDir, sourceRoot string) []string {
	relative, err := filepath.Rel(filepath.Join(dataDir, "ops"), sourceRoot)
	if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return []string{dataDir, filepath.Join(sourceRoot, "opctlgen")}
	}
	paths := []string{filepath.Join(sourceRoot, "opctlgen")}
	for _, name := range []string{"go", "dcg", "call", "vms"} {
		paths = append(paths, filepath.Join(dataDir, name))
	}
	return paths
}

func (cache *goBootstrapCache) supportedSchemas() goop.SchemaRange {
	if cache.schema == (goop.SchemaRange{}) {
		return goop.SupportedSchemas()
	}
	return cache.schema
}
