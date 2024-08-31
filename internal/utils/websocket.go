package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/nihal-ramaswamy/GoChat/internal/dto"
)

func CreateNewConnection(
	websocketConnectionMap *dto.WebsocketConnectionMap,
	upgrader *websocket.Upgrader,
	c *gin.Context,
	id string,
) error {
	data, exists := websocketConnectionMap.Get(id)
	if exists == true && data.Active {
		return nil
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return err
	}
	websocketConnectionMap.Add(id, dto.NewWebsocketConnection(conn))
	return nil
}
