package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"sync"

	"github.com/hnamzian/kv-raft/config"
	"github.com/hnamzian/kv-raft/raft"
	"github.com/hnamzian/kv-raft/rest"
	"github.com/hnamzian/kv-raft/store"
)

func main() {
	cfg := config.GetConfig()

	db := &sync.Map{}
	kvs := store.NewKeyValueStore(db)

	dataDir := "data"
	err := os.MkdirAll(dataDir, os.ModePerm)
	if err != nil {
		log.Fatalf("Could not create data directory: %s", err)
	}

	r, err := raft.SetupRaft(path.Join(dataDir, "raft"+cfg.ID), cfg.ID, "localhost:"+cfg.RaftPort, kvs)
	if err != nil {
		log.Fatal(err)
	}

	lnAddr := fmt.Sprintf(":%s", cfg.HttpPort)
	hs := rest.New(r, db, lnAddr)
	hs.Start()
}
