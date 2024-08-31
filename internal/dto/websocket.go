package dto

import (
	"sync"

	"github.com/gorilla/websocket"
)

type WebsocketConnection struct {
	Conn   *websocket.Conn
	Active bool
}

func NewWebsocketConnection(conn *websocket.Conn) *WebsocketConnection {
	return &WebsocketConnection{
		Conn:   conn,
		Active: true,
	}
}

func (wc *WebsocketConnection) Close() {
	wc.Active = false
	wc.Conn.Close()
}

// ------------------------------------------------------------------------------------------------

type WebsocketConnectionMap struct {
	mp   map[string]*WebsocketConnection
	lock sync.RWMutex
}

func NewWebsocketConnectionMap() *WebsocketConnectionMap {
	return &WebsocketConnectionMap{
		mp: make(map[string]*WebsocketConnection),
	}
}

func (wm *WebsocketConnectionMap) Add(id string, conn *WebsocketConnection) {
	wm.lock.Lock()
	defer wm.lock.Unlock()
	wm.mp[id] = conn
}

func (wm *WebsocketConnectionMap) Get(id string) (*WebsocketConnection, bool) {
	wm.lock.RLock()
	defer wm.lock.RUnlock()
	val, ok := wm.mp[id]

	return val, ok
}

func (wm *WebsocketConnectionMap) Delete(id string) {
	wm.lock.Lock()
	defer wm.lock.Unlock()
	delete(wm.mp, id)
}
