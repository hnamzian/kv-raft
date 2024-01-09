package store

import (
	"encoding/json"
	"fmt"
	"io"
	"sync"

	"github.com/hashicorp/raft"
)

type KeyValueStore struct {
	db *sync.Map
}

type setPayload struct {
	Key   string
	Value string
}

func NewKeyValueStore(db *sync.Map) *KeyValueStore {
	return &KeyValueStore{db}
}

func (s *KeyValueStore) Apply(log *raft.Log) any {
	switch log.Type {
	case raft.LogCommand:
		var sp setPayload
		if err := json.Unmarshal(log.Data, &sp); err != nil {
			return err
		}
		s.db.Store(sp.Key, sp.Value)
	default:
		return fmt.Errorf("invalid raft log type: %v", log.Type)
	}

	return nil
}

func (s *KeyValueStore) Restore(rc io.ReadCloser) error {
	s.db.Range(func(key, value any) bool {
		s.db.Delete(key)
		return true
	})

	decoder := json.NewDecoder(rc)

	for decoder.More() {
		var sp setPayload
		if err := decoder.Decode(&sp); err != nil {
			return err
		}
		s.db.Store(sp.Key, sp.Value)
	}

	return nil
}

func (kf *KeyValueStore) Snapshot() (raft.FSMSnapshot, error) {
	return snapshotNoop{}, nil
}
