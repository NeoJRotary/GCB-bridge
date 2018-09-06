package gcloud

import (
	"log"
	"time"
)

// EnableMsgListener enable message listener for testing
var EnableMsgListener = false

// MessageListener MessageListener type
type MessageListener struct {
	uid        string
	statusChan chan string
}

// MsgListener global MsgListener
var MsgListener *MessageListener

// InitMessageListener enable MsgListener, give uid as Repo.RepositoryNodeID to valid current building
func InitMessageListener(uid string) {
	EnableMsgListener = true

	MsgListener = &MessageListener{
		statusChan: make(chan string, 10),
		uid:        uid,
	}
}

func (listener *MessageListener) receive(msg *buildMsg) {
	if msg.Substitutions.BridgeUID == listener.uid {
		log.Println("Message Listener Get", msg.Status)
		listener.statusChan <- msg.Status
	} else {
		log.Println("Message Listener Skip", msg.Substitutions.BridgeUID, msg.Status)
	}
}

// ListenStatus listen status until timeout
func (listener *MessageListener) ListenStatus(timeout time.Duration) string {
	select {
	case <-time.After(timeout):
	case status := <-listener.statusChan:
		return status
	}
	return ""
}
