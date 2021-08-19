package main

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"time"

	"net/http"
	"github.com/gorilla/mux"
	smart "github.com/unix-world/smartgo"
)

const protocol = "tcp"
const nodeVersion = 1
const commandLength = 12

var nodeAddress string
var miningAddress string
var knownNodes = []string{"localhost:17770"}
var blocksInTransit = [][]byte{}
var mempool = make(map[string]Transaction)

type addr struct {
	AddrList []string
}

type block struct {
	AddrFrom string
	Block    []byte
}

type getblocks struct {
	AddrFrom string
}

type getdata struct {
	AddrFrom string
	Type     string
	ID       []byte
}

type inv struct {
	AddrFrom string
	Type     string
	Items    [][]byte
}

type tx struct {
	AddFrom     string
	Transaction []byte
}

type verzion struct {
	Version    int64
	BestHeight int64
	AddrFrom   string
}

func commandToBytes(command string) []byte {
	var bytes [commandLength]byte

	for i, c := range command {
		bytes[i] = byte(c)
	}

	return bytes[:]
}

func bytesToCommand(bytes []byte) string {
	var command []byte

	for _, b := range bytes {
		if b != 0x0 {
			command = append(command, b)
		}
	}

	return fmt.Sprintf("%s", command)
}

func extractCommand(request []byte) []byte {
	return request[:commandLength]
}

func requestBlocks() {
	for _, node := range knownNodes {
		sendGetBlocks(node)
	}
}

func sendAddr(address string) {
	nodes := addr{knownNodes}
	nodes.AddrList = append(nodes.AddrList, nodeAddress)
	payload := gobEncode(nodes)
	request := append(commandToBytes("addr"), payload...)

	sendData(address, request)
}

func sendBlock(addr string, b *Block) {
	data := block{nodeAddress, b.Serialize()}
	payload := gobEncode(data)
	request := append(commandToBytes("block"), payload...)

	sendData(addr, request)
}

func sendData(addr string, data []byte) {
	conn, err := net.Dial(protocol, addr)
	if err != nil {
		fmt.Printf("%s is not available\n", addr)
		var updatedNodes []string

		for _, node := range knownNodes {
			if node != addr {
				updatedNodes = append(updatedNodes, node)
			}
		}

		knownNodes = updatedNodes

		return
	}
	defer conn.Close()

	_, err = io.Copy(conn, bytes.NewReader(data))
	if err != nil {
		log.Panic(err)
	}
}

func sendInv(address, kind string, items [][]byte) {
	inventory := inv{nodeAddress, kind, items}
	payload := gobEncode(inventory)
	request := append(commandToBytes("inv"), payload...)

	sendData(address, request)
}

func sendGetBlocks(address string) {
	payload := gobEncode(getblocks{nodeAddress})
	request := append(commandToBytes("getblocks"), payload...)

	sendData(address, request)
}

func sendGetData(address, kind string, id []byte) {
	payload := gobEncode(getdata{nodeAddress, kind, id})
	request := append(commandToBytes("getdata"), payload...)

	sendData(address, request)
}

func sendTx(addr string, tnx *Transaction) {
	data := tx{nodeAddress, tnx.Serialize()}
	payload := gobEncode(data)
	request := append(commandToBytes("tx"), payload...)

	sendData(addr, request)
}

func sendVersion(addr string, bc *Blockchain) {
	var bestHeight int64 = int64(bc.GetBestHeight())
	payload := gobEncode(verzion{nodeVersion, bestHeight, nodeAddress})

	request := append(commandToBytes("version"), payload...)

	sendData(addr, request)
}

func handleAddr(request []byte) {
	var buff bytes.Buffer
	var payload addr

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	knownNodes = append(knownNodes, payload.AddrList...)
	fmt.Printf("There are %d known nodes now!\n", len(knownNodes))
	requestBlocks()
}

func handleBlock(request []byte, bc *Blockchain) {
	var buff bytes.Buffer
	var payload block

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	blockData := payload.Block
	block := DeserializeBlock(blockData)

	fmt.Println("Recevied a new block!")
	bc.AddBlock(block)

	fmt.Printf("Added block %x\n", block.Hash)

	if len(blocksInTransit) > 0 {
		blockHash := blocksInTransit[0]
		sendGetData(payload.AddrFrom, "block", blockHash)

		blocksInTransit = blocksInTransit[1:]
	} else {
		UTXOSet := UTXOSet{bc}
		UTXOSet.Reindex()
	}
}

func handleInv(request []byte, bc *Blockchain) {
	var buff bytes.Buffer
	var payload inv

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	fmt.Printf("Recevied inventory with %d %s\n", len(payload.Items), payload.Type)

	if payload.Type == "block" {
		blocksInTransit = payload.Items

		blockHash := payload.Items[0]
		sendGetData(payload.AddrFrom, "block", blockHash)

		newInTransit := [][]byte{}
		for _, b := range blocksInTransit {
			if bytes.Compare(b, blockHash) != 0 {
				newInTransit = append(newInTransit, b)
			}
		}
		blocksInTransit = newInTransit
	}

	if payload.Type == "tx" {
		txID := payload.Items[0]

		if mempool[hex.EncodeToString(txID)].ID == nil {
			sendGetData(payload.AddrFrom, "tx", txID)
		}
	}
}

func handleGetBlocks(request []byte, bc *Blockchain) {
	var buff bytes.Buffer
	var payload getblocks

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	blocks := bc.GetBlockHashes()
	sendInv(payload.AddrFrom, "block", blocks)
}

func handleGetData(request []byte, bc *Blockchain) {
	var buff bytes.Buffer
	var payload getdata

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	if payload.Type == "block" {
		block, err := bc.GetBlock([]byte(payload.ID))
		if err != nil {
			return
		}

		sendBlock(payload.AddrFrom, &block)
	}

	if payload.Type == "tx" {
		txID := hex.EncodeToString(payload.ID)
		tx := mempool[txID]

		sendTx(payload.AddrFrom, &tx)
		// delete(mempool, txID)
	}
}

func handleTx(request []byte, bc *Blockchain) {
	var buff bytes.Buffer
	var payload tx

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	txData := payload.Transaction
	tx := DeserializeTransaction(txData)
	mempool[hex.EncodeToString(tx.ID)] = tx

	if nodeAddress == knownNodes[0] {
		for _, node := range knownNodes {
			if node != nodeAddress && node != payload.AddFrom {
				sendInv(node, "tx", [][]byte{tx.ID})
			}
		}
	} else {
	//	if len(mempool) >= 2 && len(miningAddress) > 0 {
		if len(mempool) >= 1 && len(miningAddress) > 0 {
		fmt.Println("Start Mining ...")
		MineTransactions:
			var txs []*Transaction

			for id := range mempool {
				tx := mempool[id]
				if bc.VerifyTransaction(&tx) {
					txs = append(txs, &tx)
				}
			}

			if len(txs) == 0 {
				fmt.Println("All transactions are invalid! Waiting for new ones...")
				return
			}

			cbTx := NewCoinbaseTX(miningAddress, "")
			txs = append(txs, cbTx)

			newBlock := bc.MineBlock(txs)
			UTXOSet := UTXOSet{bc}
			UTXOSet.Reindex()

			fmt.Println("New block is mined!")

			for _, tx := range txs {
				txID := hex.EncodeToString(tx.ID)
				delete(mempool, txID)
			}

			for _, node := range knownNodes {
				if node != nodeAddress {
					sendInv(node, "block", [][]byte{newBlock.Hash})
				}
			}

			if len(mempool) > 0 {
				goto MineTransactions
			}
		}
	}
}

func handleVersion(request []byte, bc *Blockchain) {
	var buff bytes.Buffer
	var payload verzion

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	myBestHeight := bc.GetBestHeight()
	foreignerBestHeight := payload.BestHeight

	if myBestHeight < foreignerBestHeight {
		sendGetBlocks(payload.AddrFrom)
	} else if myBestHeight > foreignerBestHeight {
		sendVersion(payload.AddrFrom, bc)
	}

	// sendAddr(payload.AddrFrom)
	if !nodeIsKnown(payload.AddrFrom) {
		knownNodes = append(knownNodes, payload.AddrFrom)
	}
}

func handleConnection(conn net.Conn, bc *Blockchain) {
	request, err := ioutil.ReadAll(conn)
	if err != nil {
		log.Panic(err)
	}
	command := bytesToCommand(request[:commandLength])
	fmt.Printf("Received %s command\n", command)
//handleGetBlocks(request, bc)
	switch command {
		case "addr":
			handleAddr(request)
		case "block":
			handleBlock(request, bc)
		case "inv":
			handleInv(request, bc)
		case "getblocks":
			handleGetBlocks(request, bc)
		case "getdata":
			handleGetData(request, bc)
		case "tx":
			handleTx(request, bc)
		case "version":
			handleVersion(request, bc)
		default:
			fmt.Println("Unknown command!")
	}

	conn.Close()
}

// StartServer starts a node
func StartServer(nodeID, minerAddress string) {

	bc := NewBlockchain(nodeID)
	defer bc.db.Close()

	go func() {
		//--
		fmt.Println("Starting HTTP Mux ...")
		r := mux.NewRouter()
		//--
		r.HandleFunc("/reindexutxo", func(w http.ResponseWriter, r *http.Request) {
			//--
			UTXOSet := UTXOSet{bc}
			UTXOSet.Reindex()
			count := UTXOSet.CountTransactions()
			//--
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "Done! There are %d transactions in the UTXO set.\n", count)
			//--
		})
		//--
		r.HandleFunc("/send/from/{from}/to/{to}/amount/{amount}", func(w http.ResponseWriter, r *http.Request) {
			//--
			vars := mux.Vars(r)
			var from string = vars["from"]
			var to string = vars["to"]
			var amount int64 = smart.ParseIntegerStrAsInt64(vars["amount"])
			if(!ValidateAddress(from)) {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "Bad Request, Empty or Invalid FROM Address: %v\n", from)
				return
			}
			if(!ValidateAddress(to)) {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "Bad Request, Empty or Invalid TO Address: %v\n", to)
				return
			}
			if(amount <= 0) {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "Bad Request, Empty or Invalid AMOUNT: %d\n", amount)
				return
			}
			//--
			UTXOSet := UTXOSet{bc}
			wallets, err := NewWallets(nodeID)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "Error: %v\n", err)
				return
			}
			addresses := wallets.GetAddresses()
			var addrFound bool = false
			for _, addrs := range addresses {
				if(addrs == from) {
					addrFound = true
				}
			}
			if(addrFound != true) {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "Bad Request, Invalid Wallet Address: %v\n", from)
				return
			}
			wallet := wallets.GetWallet(from)
			tx := NewUTXOTransaction(&wallet, to, amount, &UTXOSet)
		//	if mineNow {
		//		cbTx := NewCoinbaseTX(from, "")
		//		txs := []*Transaction{cbTx, tx}
		//		newBlock := bc.MineBlock(txs)
		//		UTXOSet.Update(newBlock)
		//	} else {
				sendTx(knownNodes[0], tx)
		//	}
			//--
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "Success !\n")
				fmt.Fprintf(w, "[DEBUG] Sending to nodes ...", knownNodes[0])
			//--
		})
		//--
		r.HandleFunc("/getbalance/{address}", func(w http.ResponseWriter, r *http.Request) {
			vars := mux.Vars(r)
			var address string = vars["address"]
		//	if(len(address) < 34) {
			if(!ValidateAddress(address)) {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "Bad Request, Empty or Invalid Address: %v\n", address)
				return
			}
			//--
			UTXOSet := UTXOSet{bc}
			var balance int64 = 0
			pubKeyHash := Base58Decode([]byte(address))
			pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
			UTXOs := UTXOSet.FindUTXO(pubKeyHash)
			if(UTXOs != nil) {
				for _, out := range UTXOs {
					balance += out.Value
				}
			}
			//--
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "Balance of '%s': %d\n", address, balance)
		})
		//--
		var muxAddr string = fmt.Sprintf("localhost:%d", smart.ParseIntegerStrAsInt(nodeID)+10000)
		fmt.Println("HTTP Mux Addr:", muxAddr)
		http.ListenAndServe(muxAddr, r)
		//--
	}()

	nodeAddress = fmt.Sprintf("localhost:%s", nodeID)
	miningAddress = minerAddress
	ln, err := net.Listen(protocol, nodeAddress)
	if err != nil {
		log.Panic(err)
	}
	defer ln.Close()

	if nodeAddress != knownNodes[0] {
		sendVersion(knownNodes[0], bc)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Panic(err)
		}
		go handleConnection(conn, bc)
		time.Sleep(time.Duration(1) * time.Second)
	}
}

func gobEncode(data interface{}) []byte {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(data)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

func nodeIsKnown(addr string) bool {
	for _, node := range knownNodes {
		if node == addr {
			return true
		}
	}

	return false
}
