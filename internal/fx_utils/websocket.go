package fx_utils

import (
	"github.com/gorilla/websocket"
	"github.com/nihal-ramaswamy/GoChat/internal/dto"
	"go.uber.org/fx"
)

var WebsocketModule = fx.Module(
	"WebsocketModule",
	fx.Provide(NewWebsocketUpgrader),
	fx.Provide(dto.NewWebsocketConnectionMap),
)

func NewWebsocketUpgrader() *websocket.Upgrader {
	return &websocket.Upgrader{}
}
