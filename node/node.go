package node

import (
	"fmt"
	"github.com/appdataspec/sdk-golang/appdatapath"
	"github.com/golang-utils/lockfile"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/opctl/opctl/node/statik"
	"github.com/opspec-io/sdk-golang/node/api/handler"
	"github.com/opspec-io/sdk-golang/node/core"
	"github.com/opspec-io/sdk-golang/util/containerprovider/docker"
	"github.com/opspec-io/sdk-golang/util/pubsub"
	"github.com/rakyll/statik/fs"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func StripPrefix(prefix string, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.RequestURI = strings.Replace(r.RequestURI, prefix, "", 1)
		h.ServeHTTP(w, r)
	})
}

func New() {
	rootFSPath, err := rootFSPath()
	if nil != err {
		panic(err)
	}

	lockFile := lockfile.New()
	// ensure we're the only node around these parts
	err = lockFile.Lock(lockFilePath(rootFSPath))
	if nil != err {
		pIdOExistingNode := lockFile.PIdOfOwner(lockFilePath(rootFSPath))
		panic(fmt.Errorf("node already running w/ PId: %v\n", pIdOExistingNode))
	}

	containerProvider, err := docker.New()
	if nil != err {
		panic(err)
	}

	// cleanup [legacy] opspec.engine container if exists; ignore errors
	containerProvider.DeleteContainerIfExists("opspec.engine")

	// cleanup existing pkg cache
	pkgCachePath := filepath.Join(rootFSPath, "pkgs")
	if err := os.RemoveAll(pkgCachePath); nil != err {
		panic(fmt.Errorf("unable to cleanup pkg cache at path: %v\n", pkgCachePath))
	}

	// cleanup existing DCG (dynamic call graph) data
	dcgDataDirPath := dcgDataDirPath(rootFSPath)
	if err := os.RemoveAll(dcgDataDirPath); nil != err {
		panic(fmt.Errorf("unable to cleanup DCG (dynamic call graph) data at path: %v\n", dcgDataDirPath))
	}

	router := mux.NewRouter()
	router.UseEncodedPath()

	router.PathPrefix("/api/").Handler(
		StripPrefix("/api",
			handler.New(
				core.New(
					pubsub.New(pubsub.NewEventRepo(eventDbPath(dcgDataDirPath))),
					containerProvider,
					rootFSPath,
				),
			),
		),
	)

	statikFS, err := fs.New()
	if nil != err {
		panic(err)
	}
	router.PathPrefix("/").Handler(http.FileServer(statikFS))

	err = http.ListenAndServe(":42224", handlers.CORS()(router))
	if nil != err {
		panic(err)
	}

}

// fsRootPath returns the root fs path for the node
func rootFSPath() (string, error) {
	perUserAppDataPath, err := appdatapath.New().PerUser()
	if nil != err {
		return "", err
	}

	return path.Join(
		perUserAppDataPath,
		"opctl",
	), nil
}

func dcgDataDirPath(rootFSPath string) string {
	return path.Join(
		rootFSPath,
		"dcg",
	)
}

func eventDbPath(dcgDataDirPath string) string {
	return path.Join(
		dcgDataDirPath,
		"event.db",
	)
}

func lockFilePath(rootFSPath string) string {
	return path.Join(
		rootFSPath,
		"pid.lock",
	)
}
