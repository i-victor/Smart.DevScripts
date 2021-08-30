package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"sync"

	"github.com/chain/txvm/protocol/bc"
	"github.com/chain/txvm/protocol/state"
)

// Implements protocol.Store from github.com/chain/txvm
type pstore struct{}

func (s *pstore) Height(_ context.Context) (uint64, error) {
	infos, err := ioutil.ReadDir(blockDir())
	if err != nil {
		return 0, err
	}
	var result int64
	for _, info := range infos {
		n, err := strconv.ParseInt(info.Name(), 10, 64)
		if err != nil {
			continue
		}
		if n > result {
			result = n
		}
	}
	return uint64(result), nil
}

func (s *pstore) GetBlock(_ context.Context, height uint64) (*bc.Block, error) {
	return readBlockFile(path.Join(blockDir(), strconv.FormatUint(height, 10)))
}

func (s *pstore) LatestSnapshot(_ context.Context) (*state.Snapshot, error) {
	infos, err := ioutil.ReadDir(snapshotDir())
	if err != nil {
		return nil, err
	}
	var highest int
	for _, info := range infos {
		n, err := strconv.Atoi(info.Name())
		if err != nil {
			continue
		}
		if n > highest {
			highest = n
		}
	}
	if highest <= 0 {
		return state.Empty(), nil
	}
	filename := path.Join(snapshotDir(), strconv.Itoa(highest))
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var snapshot state.Snapshot
	err = snapshot.FromBytes(b)
	if err != nil {
		return nil, err
	}
	return &snapshot, nil
}

func (s *pstore) SaveBlock(_ context.Context, block *bc.Block) error {
	err := storeBlock(block)
	if err != nil {
		return err
	}
	oldName := blockFilename(int(block.Height), block.Hash()) // xxx uint64->int
	newName := path.Join(blockDir(), strconv.FormatUint(block.Height, 10))
	err = os.Link(oldName, newName)
	if os.IsExist(err) {
		return nil
	}
	return err
}

func (s *pstore) FinalizeHeight(_ context.Context, height uint64) error {
	return nil
}

func (s *pstore) SaveSnapshot(_ context.Context, snapshot *state.Snapshot) error {
	filename := path.Join(snapshotDir(), strconv.FormatUint(snapshot.Height(), 10))
	b, err := snapshot.Bytes()
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, b, 0644)
}

func blockFilename(height int, id bc.Hash) string {
	return path.Join(blockDir(), fmt.Sprintf("%d-%x", height, id.Bytes()))
}

func getBlock(height int, id bc.Hash) (*bc.Block, error) {
	return readBlockFile(blockFilename(height, id))
}

func haveBlock(height int, id bc.Hash) (bool, error) {
	filename := blockFilename(height, id)
	_, err := os.Stat(filename)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func readBlockFile(filename string) (*bc.Block, error) {
	bits, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var block bc.Block
	err = block.FromBytes(bits)
	if err != nil {
		return nil, err
	}
	return &block, nil
}

var storeBlockMu sync.Mutex

func storeBlock(block *bc.Block) error {
	storeBlockMu.Lock()
	defer storeBlockMu.Unlock()

	filename := blockFilename(int(block.Height), block.Hash()) // xxx uint64->int
	_, err := os.Stat(filename)
	if err == nil {
		// File exists already.
		return nil
	}
	if !os.IsNotExist(err) {
		// Problem is other than file-doesn't-exist.
		return err
	}
	bits, err := block.Bytes()
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, bits, 0644)
}

func blockDir() string {
	return path.Join(dir, "blocks")
}

func snapshotDir() string {
	return path.Join(dir, "snapshots")
}
