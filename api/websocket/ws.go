package websocket

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"im/db"
	"im/schema"
	"log"
	"net/http"
)

type MessageBasic2 struct {
	UserIdentity string `json:"user_identity" bson:"user_identity"`
	RoomIdentity string `json:"room_identity" bson:"room_identity"`
	Data         string `json:"data" bson:"data"`
	CreatedAt    string `json:"created_at" bson:"created_at"`
	UpdatedAt    string `json:"updated_at" bson:"updated_at"`
}

// 处理WS跨域
var upgrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
var ws = make(map[string]*websocket.Conn)

func Ws(c *gin.Context) {
	// 升级为WS协议 s
	conn, err := upgrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "系统异常：" + err.Error(),
		})
		return
	}
	// 升级为WS协议 e
	for {
		// 创建消息JSON结构体, 保存消息与额外的信息
		message := new(schema.MessageBasic)

		fmt.Println("ws:", ws)
		// ReadJSON
		// message: 需要读取的消息对象(一条消息一般包含多个属性用于其他用途)
		err := conn.ReadJSON(message)
		if err != nil {
			log.Println("read:", err)
			break
		}

		fmt.Println("message:", message)
		var savaMessage = schema.MessageBasic{
			UserIdentity: message.UserIdentity,
			RoomIdentity: message.RoomIdentity,
			Data:         message.Data,
			CreatedAt:    message.CreatedAt,
			UpdatedAt:    message.UpdatedAt,
		}
		jsonMessage, err := json.Marshal(&savaMessage)
		if err != nil {
			log.Println("转换JSON失败:" + err.Error())
		}
		result, err := db.Redis.RPush(message.UserIdentity, jsonMessage).Result()
		if err != nil {
			log.Println("存储消息失败:" + err.Error())
		}
		log.Println("result:", result)
		log.Println("jsonMessage:", jsonMessage)
		ws[savaMessage.UserIdentity] = conn // 根据连接绑定用户id
		log.Println("savaMessage.UserIdentity:", savaMessage.UserIdentity)
		log.Println("savaMessage:", savaMessage)
		for _, cc := range ws {
			// WriteMessage
			// 1 消息类型: websocket.TextMessage文本
			// 2 传输类型: []byte二进制
			err = cc.WriteJSON(savaMessage)
			if err != nil {
				log.Println("write:", err)
				break
			}
		}
	}
}
