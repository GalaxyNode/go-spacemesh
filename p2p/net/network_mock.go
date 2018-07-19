package net

import (
	"github.com/spacemeshos/go-spacemesh/crypto"
	"github.com/spacemeshos/go-spacemesh/log"
	"github.com/spacemeshos/go-spacemesh/p2p/net/wire"
	"gopkg.in/op/go-logging.v1"
	"time"
	"net"
	"sync/atomic"
)

type ReadWriteCloserMock struct {
}

func (m ReadWriteCloserMock) Read(p []byte) (n int, err error) {
	return 0, nil
}

func (m ReadWriteCloserMock) Write(p []byte) (n int, err error) {
	return 0, nil
}

func (m ReadWriteCloserMock) Close() error {
	return nil
}

func (m ReadWriteCloserMock) RemoteAddr() net.Addr {
	r, err := net.ResolveTCPAddr("tcp", "127.0.0.0")
	if err != nil {
		panic(err)
	}
	return r
}

func getTestLogger(name string) *logging.Logger {
	return log.New(name, "", "").Logger
}

type NetworkMock struct {
	dialErr          error
	dialDelayMs      int8
	dialCount			int32
	regNewRemoteConn []chan Connectioner
	networkId        int8
	closingConn      chan Connectioner
	incomingMessages      chan IncomingMessageEvent
	logger           *logging.Logger
}

func NewNetworkMock() *NetworkMock {
	return &NetworkMock{
		regNewRemoteConn: make([]chan Connectioner, 0),
		closingConn:      make(chan Connectioner),
		logger:           getTestLogger("network mock"),
	}
}

func (n * NetworkMock) reset() {
	n.dialCount = 0
	n.dialDelayMs = 0
	n.dialErr = nil
}

func (n *NetworkMock) SetDialResult(err error) {
	n.dialErr = err
}

func (n *NetworkMock) SetDialDelayMs(delay int8) {
	n.dialDelayMs = delay
}



func (n *NetworkMock) Dial(address string, remotePublicKey crypto.PublicKey, networkId int8) (Connectioner, error) {
	n.networkId = networkId
	atomic.AddInt32(&n.dialCount, 1)
	time.Sleep(time.Duration(n.dialDelayMs) * time.Millisecond)
	conn := NewConnection(ReadWriteCloserMock{}, n, Local, remotePublicKey, n.logger)
	return conn, n.dialErr
}

func (n *NetworkMock) GetDialCount() int32 {
	return n.dialCount
}

func (n *NetworkMock) SubscribeOnNewRemoteConnections() chan Connectioner {
	ch := make(chan Connectioner, 20)
	n.regNewRemoteConn = append(n.regNewRemoteConn, ch)
	return ch
}

func (n NetworkMock) PublishNewRemoteConnection(conn Connectioner) {
	for _, ch := range n.regNewRemoteConn {
		ch <- conn
	}
}

func (n *NetworkMock) setNetworkId(id int8) {
	n.networkId = id
}

func (n *NetworkMock) GetNetworkId() int8 {
	return n.networkId
}

func (n *NetworkMock) ClosingConnections() chan Connectioner {
	return n.closingConn
}

func (n* NetworkMock) IncomingMessages() chan IncomingMessageEvent {
	return n.incomingMessages
}

func (n NetworkMock) PublishClosingConnection(conn Connectioner) {
	go func() {
		n.closingConn <- conn
	}()
}

func (n *NetworkMock) HandlePreSessionIncomingMessage(c Connectioner, msg wire.InMessage) error {
	return nil
}

func (n *NetworkMock) GetLogger() *logging.Logger {
	return n.logger
}
