package node

import (
	"context"
	"path/filepath"
	"slices"
	"sync"
)

// One node process owns a data dir, so in-process locks suffice for its cache dirs.
type cacheDirLocks struct {
	mu         sync.Mutex
	byHostPath map[string]*cacheDirLock
}

type cacheDirLock struct {
	held  chan struct{}
	users int // includes holders and waiters; discard unused locks, including cold-run paths
}

func (locks *cacheDirLocks) acquire(ctx context.Context, hostPaths []string, onWait func(string)) (func(), error) {
	paths := make([]string, len(hostPaths))
	for i, hostPath := range hostPaths {
		paths[i] = filepath.Clean(hostPath)
	}
	slices.Sort(paths)
	paths = slices.Compact(paths)
	if len(paths) == 0 {
		return func() {}, nil
	}

	// Every call takes distinct paths in the same order, avoiding deadlocks.
	var acquired []*cacheDirLock
	release := func() {
		for i := len(acquired) - 1; i >= 0; i-- {
			locks.release(paths[i], acquired[i], true)
		}
	}
	for _, hostPath := range paths {
		if err := ctx.Err(); err != nil {
			release()
			return nil, err
		}
		locks.mu.Lock()
		if locks.byHostPath == nil {
			locks.byHostPath = map[string]*cacheDirLock{}
		}
		lock := locks.byHostPath[hostPath]
		if lock == nil {
			lock = &cacheDirLock{held: make(chan struct{}, 1)}
			locks.byHostPath[hostPath] = lock
		}
		lock.users++
		locks.mu.Unlock()

		select {
		case lock.held <- struct{}{}:
		default:
			onWait(hostPath)
			select {
			case lock.held <- struct{}{}:
			case <-ctx.Done():
				locks.release(hostPath, lock, false)
				release()
				return nil, ctx.Err()
			}
		}
		acquired = append(acquired, lock)
	}
	if err := ctx.Err(); err != nil {
		release()
		return nil, err
	}
	// a second release would block on an empty channel while holding mu
	return sync.OnceFunc(release), nil
}

func (locks *cacheDirLocks) release(hostPath string, lock *cacheDirLock, held bool) {
	locks.mu.Lock()
	defer locks.mu.Unlock()
	if held {
		<-lock.held
	}
	lock.users--
	if lock.users == 0 {
		delete(locks.byHostPath, hostPath)
	}
}
