package router

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// Mux router engine with websocket upgrader
var (
	R               *mux.Router
	Upgrader        websocket.Upgrader
	LoginFailuresCh = make(chan string)
)

// NewRouter init
func NewRouter() *mux.Router {

	// Router instance
	R = mux.NewRouter()

	// Websocket connection upgrader
	Upgrader = websocket.Upgrader{} // initialize with default options

	return R
}
