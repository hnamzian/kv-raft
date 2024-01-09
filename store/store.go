package store

import (
	"io"

	"github.com/hashicorp/raft"
)

type Store interface {
	Apply(log *raft.Log) any
	Restore(rc io.ReadCloser) error
	Snapshot() (raft.FSMSnapshot, error)
}

type snapshotNoop struct{}

func (sn snapshotNoop) Persist(_ raft.SnapshotSink) error {
	return nil
}

func (sn snapshotNoop) Release() {}
