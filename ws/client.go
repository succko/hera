// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ws

import (
	"bytes"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/nacos-group/nacos-sdk-go/v2/util"
	"github.com/succko/hera/mq"
	"github.com/succko/hera/pb"
	"github.com/succko/hera/utils"
	"go.uber.org/zap"
	"math"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	// writeWait 定义了向对等端写入消息的时间限制。
	writeWait = 10 * time.Second

	// pongWait 定义了从对等端读取下一个pong消息的时间限制。
	pongWait = 60 * time.Second

	// pingPeriod 定义了向对等端发送ping消息的周期。必须小于pongWait。
	pingPeriod = (pongWait * 9) / 10

	// maxMessageSize 定义了允许从对等端接收的最大消息大小。
	maxMessageSize = 512

	// secretKey 定义了用于计算签名的密钥。
	secretKey = "%YhM=ULje*eX"

	// timeout 定义了客户端超时时间。
	timeout = float64(60 * time.Second * 3600)
)

var (
	newline = []byte{'\n'} // newline是一个包含换行符的字节切片。
	space   = []byte{' '}  // space是一个包含空格的字节切片。
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// 解决跨域问题
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Client 是websocket连接和hub之间的中间件。
type Client struct {
	hub  *Hub            // hub是该Client所属的Hub实例。
	conn *websocket.Conn // conn是与websocket连接相关的网络连接。
	send chan []byte     // send是用于向hub发送消息的缓冲通道。
	uuid string          // uuid是客户端的唯一标识符。
	hbts int             // 最后一次心跳时间
}

// readPump 将websocket连接上的消息泵送到hub。
//
// 应用程序在每个连接上运行readPump goroutine。通过在此goroutine中执行所有读取操作，确保连接上最多只有一个读取器。
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		err := c.conn.Close()
		if err != nil {
			zap.L().Error("readPump close error", zap.Error(err))
			return
		}
	}()
	c.conn.SetReadLimit(maxMessageSize)
	err := c.conn.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		zap.L().Error("readPump SetReadDeadline error", zap.Error(err))
		return
	}
	c.conn.SetPongHandler(func(string) error {
		err := c.conn.SetReadDeadline(time.Now().Add(pongWait))
		if err != nil {
			zap.L().Error("readPump SetReadDeadline error", zap.Error(err))
			return err
		}
		return nil
	})
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				zap.L().Error("readPump ReadMessage error", zap.Error(err))
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		zap.L().Info("readPump ReadMessage message:" + string(message))
		c.handleC2S(message)
	}
}

// writePump 将hub上的消息泵送到websocket连接。
//
// 为每个连接启动一个写入泵goroutine。通过在此goroutine中执行所有写入操作，确保连接上最多只有一个写入器。
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		err := c.conn.Close()
		if err != nil {
			zap.L().Error("writePump close error", zap.Error(err))
			return
		}
	}()
	for {
		select {
		case message, ok := <-c.send:
			err := c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				zap.L().Error("writePump SetWriteDeadline error", zap.Error(err))
				return
			}
			if !ok {
				// hub关闭了通道。
				err := c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					zap.L().Error("writePump conn.WriteMessage error", zap.Error(err))
					return
				}
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				zap.L().Error("writePump conn.NextWriter error", zap.Error(err))
				return
			}
			_, err = w.Write(message)
			if err != nil {
				zap.L().Error("writePump conn.Write error", zap.Error(err))
				return
			}

			// 将队列中的聊天消息添加到当前websocket消息中。
			n := len(c.send)
			for i := 0; i < n; i++ {
				_, err := w.Write(newline)
				if err != nil {
					zap.L().Error("writePump w.Write error", zap.Error(err))
					return
				}
				_, err = w.Write(<-c.send)
				if err != nil {
					zap.L().Error("writePump c.send <- c.send error", zap.Error(err))
					return
				}
			}

			if err := w.Close(); err != nil {
				zap.L().Error("writePump w.Close error", zap.Error(err))
				return
			}
		case <-ticker.C:
			err := c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				zap.L().Error("writePump SetWriteDeadline error", zap.Error(err))
				return
			}
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				zap.L().Error("writePump conn.WriteMessage error", zap.Error(err))
				return
			}
		}
	}
}

// ServeWs 处理来自对等端的websocket请求。
func ServeWs(w http.ResponseWriter, r *http.Request) {
	// 将HTTP响应写入和请求作为参数传递给WebSocket升级函数
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		zap.L().Error("ServeWs upgrader.Upgrade err", zap.Error(err))
		return
	}
	// 创建一个新的客户端实例并将其连接到WebSocket Hub
	c := &Client{hub: SingletonHub(), conn: conn, send: make(chan []byte, 256)}
	c.hub.register <- c

	// 使用新的goroutine执行读取和写入循环，允许调用者引用与客户端相关联的内存
	go c.writePump()
	go c.readPump()
}

// 发送消息给客户端
func (c *Client) sendMessage(message []byte) {
	c.send <- message
	zap.L().Info("sendMessage message:"+string(message), zap.String("uuid", c.uuid))
}

// SendMessage 发送消息给客户端
func sendMessage(uuid string, innoPacket *pb.InnoPacket) {
	clients, ok := SingletonHub().ids[uuid]
	if ok {
		message, _ := proto.Marshal(innoPacket)
		for _, c := range clients {
			c.sendMessage(message)
		}
	}
}

func checkSign(heartBeat *pb.HeartBeatPacket) bool {
	if math.Abs(float64(heartBeat.GetTs()-time.Now().Unix())) > timeout {
		zap.L().Error("checkSign heartBeat timeout")
		return false
	}
	sign := util.Md5(heartBeat.GetId() + strconv.Itoa(int(heartBeat.GetTs())) + secretKey)
	return sign == heartBeat.GetSign()
}

func (c *Client) auth(innoPacket *pb.InnoPacket) {
	var mu sync.Mutex
	if !checkSign(innoPacket.GetHeartBeat()) {
		zap.L().Error("handleC2S heartBeat sign error, message:"+innoPacket.GetHeartBeat().String(), zap.String("uuid", c.uuid))
		return
	}
	// 更新客户端uuid
	if c.uuid != "" && c.uuid != innoPacket.GetHeartBeat().GetId() {
		zap.L().Error("handleC2S heartBeat uuid error, message:"+innoPacket.GetHeartBeat().String(), zap.String("uuid", c.uuid))
		return
	}
	c.uuid = innoPacket.GetHeartBeat().GetId()
	c.hbts = int(innoPacket.GetHeartBeat().GetTs())

	// 判断是否已经存在该客户端
	func() {
		mu.Lock()
		defer mu.Unlock()
		if c.hub.ids == nil {
			c.hub.ids = make(map[string][]*Client)
		}
		flag := false
		for _, v := range c.hub.ids[c.uuid] {
			if v == c {
				flag = true
				break
			}
		}
		if !flag {
			c.hub.ids[c.uuid] = append(c.hub.ids[c.uuid], c)
		}
	}()

	//c.hub.boys[c.uuid] = c
	//c.hub.girls[c.uuid] = c
	zap.L().Info("hub client heartBeat auth", zap.Int("clients", len(c.hub.clients)), zap.String("uuid", c.uuid), zap.Int("ids", len(c.hub.ids[c.uuid])), zap.String("all", strings.Join(utils.MapKeys(c.hub.ids), ",")))
}

func (c *Client) handleC2S(message []byte) {
	if strings.ToUpper(string(message)) == "PING" {
		c.sendMessage([]byte("PONG"))
		return
	}
	// 根据不同的指令，执行不同的操作
	var innoPacket = &pb.InnoPacket{}
	err := proto.Unmarshal(message, innoPacket)
	if err != nil {
		err := jsonpb.UnmarshalString(string(message), innoPacket)
		if err != nil {
			zap.L().Error("handleC2S message:"+string(message)+" error="+err.Error(), zap.String("uuid", c.uuid))
			return
		}
	}
	zap.L().Info("handleC2S innoPacket:"+innoPacket.String(), zap.String("uuid", c.uuid))
	if innoPacket.Type == pb.InnoPacket_TYPE_HEARTBEAT {
		c.auth(innoPacket)
	} else if innoPacket.Type == pb.InnoPacket_TYPE_INSTRUCTION {
		if innoPacket.GetInstruction().GetToId() != "" {
			sendMessage(innoPacket.GetInstruction().GetToId(), innoPacket)
		} else {
			// 生产MQ消息
			mq.Producer.SendSync("instruct_c2s", innoPacket)
		}
	}
}
func InstructS2cHandle(message []byte) {
	var innoPacket = &pb.InnoPacket{}
	err := jsonpb.UnmarshalString(string(message), innoPacket)
	if err != nil {
		err = proto.Unmarshal(message, innoPacket)
		if err != nil {
			zap.L().Error("HandleS2C message:" + string(message) + " error=" + err.Error())
			return
		}
	}
	HandleS2C(innoPacket)
}

func HandleS2C(innoPacket *pb.InnoPacket) {
	// 根据不同的指令，执行不同的操作
	zap.L().Info("HandleS2C innoPacket:" + innoPacket.String())
	if innoPacket.Type == pb.InnoPacket_TYPE_HEARTBEAT {
		return
	} else if innoPacket.Type == pb.InnoPacket_TYPE_INSTRUCTION {
		if innoPacket.GetInstruction().GetToId() != "" {
			if strings.ToUpper(innoPacket.GetInstruction().GetToId()) == "ALL" {
				SingletonHub().Broadcast(innoPacket)
			} else {
				sendMessage(innoPacket.GetInstruction().GetToId(), innoPacket)
			}
		}
	}
}
