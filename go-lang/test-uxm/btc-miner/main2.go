package main

import (
	"github.com/rbrtl/gobtcminer/client"
	"github.com/rbrtl/gobtcminer/config"
	"github.com/rbrtl/gobtcminer/logger"
	"github.com/rbrtl/gobtcminer/mining"
	"time"
)

var monitor logger.Logger

func main() {

	// Read and parse the configuration file
	conf := config.Readconf("config.json")
	monitor = logger.NewLogger(conf.Log)
	diff, err := client.GetDifficulty(conf.User, conf.Password, conf.Host)
	if err != nil {
		monitor.Print("info", "Error getting difficulty: "+err.Error())
	}
	dispatcher := mining.NewDispatcher(monitor)
	dispatcher.Run()
	//Add Chunks on a regular basis
	for {
		//Get a new Chunk and split it accordingly to the machin settings
		for _, chunk := range mining.NewChunkList(2, uint32(time.Now().Unix()), diff) {
			if len(dispatcher.ChunkQueueIn) < cap(dispatcher.ChunkQueueIn) {
				dispatcher.ChunkQueueIn <- chunk
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
}
