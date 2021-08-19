
package main

// message.go
// r.20210819.0346 :: STABLE

import (
	"log"

	"strings"
	"strconv"

	"github.com/gorilla/websocket"

	smart "github.com/unix-world/smartgo"
	uid   "github.com/unix-world/smartgo/uuid"
	b58   "github.com/unix-world/smartgo/base58"
)

type messagePack struct {
	Cmd        string `json:"cmd"`
	Data       string `json:"data"`
	CheckSum   string `json:"checksum"`
}

func ComposePakMessage(cmd string, data string) (msg string, errMsg string) {
	cmd = smart.StrTrimWhitespaces(cmd)
	if(cmd == "") {
		return "", "Command is empty"
	}
	var dataEnc string = smart.BlowfishEncryptCBC(data, cmd)
	sMsg := &messagePack{
		cmd,
		dataEnc,
		smart.Sha512(cmd + "\n" + dataEnc + "\n" + data),
	}
	return smart.DataArchive(smart.JsonEncode(sMsg)), ""
}

func ParsePakMessage(msg string) (msgStruct *messagePack, errMsg string) {
	msg = smart.StrTrimWhitespaces(msg)
	if(msg == "") {
		return nil, "Message is empty"
	}
	msg = smart.DataUnArchive(msg)
	if(msg == "") {
		return nil, "Message Unarchiving FAILED"
	}
	msg = smart.StrTrimWhitespaces(msg)
	if(msg == "") {
		return nil, "Message is empty after Unarchiving"
	}
	D := smart.JsonDecode(msg)
	if(D == nil) {
		return nil, "Message Decoding FAILED"
	}
	sMsg := &messagePack{
		D["cmd"].(string),
		D["data"].(string),
		D["checksum"].(string),
	}
	sMsg.Data = smart.BlowfishDecryptCBC(sMsg.Data, sMsg.Cmd)
	if(sMsg.CheckSum != smart.Sha512(sMsg.Cmd + "\n" + D["data"].(string) + "\n" + sMsg.Data)) {
		return nil, "Invalid Message Checksum"
	}
	return sMsg, ""
}

func HandleMessage(id string, message []byte, conn *websocket.Conn) (ok bool, errMsg string) {
	msg, errMsg := ParsePakMessage(string(message))
	message = nil
	if(errMsg != "") {
		return false, errMsg
	}
	log.Println("[DEBUG] Received Command: `" + msg.Cmd + "`")
	log.Println("[DEBUG] Data Length:", len(msg.Data))
	if(conn == nil) {
		log.Println("[DATA] Data:\n", msg.Data, "\n\n")
		return true, ""
	}
	var answer string = "Got Command from [" + id + "]: " + msg.Cmd + "\n" + strings.Repeat(msg.Data, 50)
	AnswerMsg, errAnswerMsg := ComposePakMessage("OK: " + msg.Cmd, answer + "\n" + "Length=" + strconv.Itoa(len(answer)) + "\n")
	if(errAnswerMsg != "") {
		return false, "Answer Message Error: " + errAnswerMsg
	}
	err := conn.WriteMessage(websocket.TextMessage, []byte(AnswerMsg))
	if err != nil {
		return false, "Error during client answer to websocket: " + err.Error()
	}
	return true, ""
}

func GenerateUUID() string {
	var theTime string = ""
	dtObjUtc := smart.DateTimeStructUtc("")
	if(dtObjUtc.Status != "OK") {
		log.Println("[ERROR] Date Time Failed:", dtObjUtc.ErrMsg)
	} else {
		theTime = smart.ConvertInt64ToStr(dtObjUtc.Time)
	}
	log.Println("Time Seed:", theTime)
	var uuid string = uid.Uuid1013Str(13) + "-" + uid.Uuid1013Str(10) + "-" + uid.Uuid1013Str(13);
	if(theTime != "") {
		uuid += "-" + b58.Encode([]byte(theTime))
	}
	return uuid
}

// #END
