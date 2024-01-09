package config

import (
	"log"
	"os"
)

type Config struct {
	ID       string
	HttpPort string
	RaftPort string
}

func GetConfig() Config {
	cfg := Config{}
	for i, arg := range os.Args[1:] {
		if arg == "--node-id" {
			cfg.ID = os.Args[i+2]
			i++
			continue
		}

		if arg == "--http-port" {
			cfg.HttpPort = os.Args[i+2]
			i++
			continue
		}

		if arg == "--raft-port" {
			cfg.RaftPort = os.Args[i+2]
			i++
			continue
		}
	}

	if cfg.ID == "" {
		log.Fatal("Missing required parameter: --node-id")
	}

	if cfg.RaftPort == "" {
		log.Fatal("Missing required parameter: --raft-port")
	}

	if cfg.HttpPort == "" {
		log.Fatal("Missing required parameter: --http-port")
	}

	return cfg
}
