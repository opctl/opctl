package core

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/dgraph-io/badger/v2"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/pubsub"
)

//counterfeiter:generate -o internal/fakes/stateStore.go . stateStore
// stateStore allows efficiently querying the current state of opctl.
//
// State is materialized by applying events in the order in which they are/were received.
//
// efficient startup:
// A lastAppliedEventTimestamp is maintained and used at startup to pickup applying events
// from where we left off.
type stateStore interface {
	// lists all calls w/ parentID
	ListWithParentID(parentID string) []*model.Call

	TryGet(id string) *model.Call

	// TryGetCreds returns creds for a ref if any exist
	TryGetAuth(resource string) *model.Auth
}

func newStateStore(
	db *badger.DB,
	pubSub pubsub.PubSub,
) stateStore {

	stateStore := &_stateStore{
		authsByResourcesKeyPrefix:    "authsByResources_",
		callsByID:                    make(map[string]*model.Call),
		db:                           db,
		lastAppliedEventTimestampKey: "lastAppliedEventTimestamp",
	}

	go func() {
		// apply events in background

		// make best effort to get lastAppliedEventTimestamp
		lastAppliedEventTimestamp, _ := stateStore.getLastAppliedEventTimestamp()

		// replay from second before last applied event to ensure we see events
		// at least once (multiple events w/ same timestamp could exist)
		since := lastAppliedEventTimestamp.Add(-time.Second)

		eventChannel, _ := pubSub.Subscribe(
			context.Background(),
			model.EventFilter{
				Since: &since,
			},
		)

		for event := range eventChannel {
			switch {
			case nil != event.AuthAdded:
				stateStore.applyAuthAdded(*event.AuthAdded)
			case nil != event.CallEnded:
				stateStore.applyCallEnded(*event.CallEnded)
			case nil != event.CallStarted:
				stateStore.applyCallStarted(*event.CallStarted)
			}

			stateStore.updateLastAppliedEventTimestamp(event.Timestamp)
		}
	}()

	return stateStore

}

type _stateStore struct {
	lastAppliedEventTimestampKey string
	authsByResourcesKeyPrefix    string
	callsByID                    map[string]*model.Call
	db                           *badger.DB
	// synchronize access via mutex
	mux sync.RWMutex
}

func (ss *_stateStore) getLastAppliedEventTimestamp() (time.Time, error) {
	var lastAppliedEventTimestamp time.Time
	err := ss.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(ss.lastAppliedEventTimestampKey))
		if nil != err {
			return err
		}

		return item.Value(func(val []byte) error {
			seconds, err := strconv.ParseInt(string(val), 10, 64)
			if nil != err {
				return err
			}
			lastAppliedEventTimestamp = time.Unix(seconds-1, 0)
			return nil
		})
	})

	return lastAppliedEventTimestamp, err
}

func (ss *_stateStore) updateLastAppliedEventTimestamp(lastAppliedEventTimestamp time.Time) error {
	return ss.db.Update(func(txn *badger.Txn) error {
		return txn.Set(
			[]byte(ss.lastAppliedEventTimestampKey),
			[]byte(
				strconv.FormatInt(lastAppliedEventTimestamp.Unix(), 10),
			),
		)
	})
}

func (ss *_stateStore) applyAuthAdded(authAdded model.AuthAdded) error {
	return ss.db.Update(func(txn *badger.Txn) error {
		auth := authAdded.Auth
		encodedAuth, err := json.Marshal(auth)
		if nil != err {
			return err
		}

		return txn.Set(
			[]byte(ss.authsByResourcesKeyPrefix+strings.ToLower(auth.Resources)),
			encodedAuth,
		)
	})
}

func (ss *_stateStore) applyCallEnded(callEnded model.CallEnded) {
	if callEnded.Outcome != model.OpOutcomeFailed {
		return
	}

	ss.mux.RLock()
	defer ss.mux.RUnlock()

	callID := callEnded.Call.ID
	if _, ok := ss.callsByID[callID]; ok {
		ss.callsByID[callID].IsKilled = true
	}
}

// O(1) complexity
func (ss *_stateStore) applyCallStarted(callStarted model.CallStarted) {
	ss.mux.Lock()
	defer ss.mux.Unlock()

	call := callStarted.Call
	ss.callsByID[call.ID] = &call
}

// O(n) complexity (n being active call count)
func (ss *_stateStore) ListWithParentID(parentID string) []*model.Call {
	ss.mux.RLock()
	defer ss.mux.RUnlock()

	results := []*model.Call{}
	for _, call := range ss.callsByID {
		if nil != call.ParentID && *call.ParentID == parentID {
			results = append(results, call)
		}
	}
	return results
}

func (ss *_stateStore) TryGet(
	id string,
) *model.Call {
	ss.mux.RLock()
	defer ss.mux.RUnlock()

	if call, ok := ss.callsByID[id]; ok {
		return call
	}

	return nil
}

func (ss *_stateStore) TryGetAuth(
	ref string,
) *model.Auth {
	ref = strings.ToLower(ref)
	var auth *model.Auth
	ss.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()
		prefixBytes := []byte(ss.authsByResourcesKeyPrefix)
		for it.Seek(prefixBytes); it.ValidForPrefix(prefixBytes); it.Next() {
			item := it.Item()
			key := string(item.Key())
			prefix := strings.TrimPrefix(key, ss.authsByResourcesKeyPrefix)

			if strings.HasPrefix(ref, prefix) {
				item.Value(func(value []byte) error {
					auth = &model.Auth{}
					return json.Unmarshal(value, auth)
				})
			}
		}
		return nil
	})

	return auth
}
