package node

import (
	"testing"

	"github.com/dgraph-io/badger/v4"
	"github.com/opctl/opctl/sdks/go/model"
)

func TestBootstrapAuthSnapshot(t *testing.T) {
	database, err := badger.Open(badger.DefaultOptions("").WithInMemory(true).WithLogger(nil))
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	store := &_stateStore{db: database, authsByResourcesKeyPrefix: "auths_"}
	original := model.Auth{Resources: "Example.com/team", Creds: model.Creds{Username: "u", Password: "synthetic-snapshot-token"}}
	if err := store.applyAuthAdded(model.AuthAdded{Auth: original}); err != nil {
		t.Fatal(err)
	}
	snapshot, err := store.AuthSnapshot()
	if err != nil {
		t.Fatal(err)
	}
	if len(snapshot) != 1 || snapshot[0] != original {
		t.Fatal("snapshot did not return original auth records")
	}
	snapshot[0].Password = "changed"
	snapshot, err = store.AuthSnapshot()
	if err != nil {
		t.Fatal(err)
	}
	if snapshot[0] != original {
		t.Fatal("snapshot mutated stored auth")
	}
}
