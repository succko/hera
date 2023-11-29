package ws

import (
	"github.com/golang/protobuf/proto"
	"github.com/succko/hera/pb"
	"github.com/succko/hera/utils"
	"go.uber.org/zap"
	"strings"
	"sync"
)

// Hub Hub维护活跃客户端的集合，并向客户端广播消息。
type Hub struct {
	// 注册的客户端。
	clients map[*Client]bool

	// 客户端的ID集合。
	ids map[string][]*Client

	// 男生的客户端集合。
	boys map[string]*Client

	// 女生的客户端集合。
	girls map[string]*Client

	// 客户端发送的入站消息。
	broadcast chan []byte

	// 客户端的注册请求。
	register chan *Client

	// 客户端的注销请求。
	unregister chan *Client

	// 运行状态。
	running bool
}

var (
	hub  *Hub
	once sync.Once
)

// SingletonHub 获取单例Hub实例。
func SingletonHub() *Hub {
	once.Do(func() {
		hub = &Hub{
			clients:    make(map[*Client]bool),
			ids:        make(map[string][]*Client),
			boys:       make(map[string]*Client),
			girls:      make(map[string]*Client),
			broadcast:  make(chan []byte),
			register:   make(chan *Client),
			unregister: make(chan *Client),
			running:    false,
		}
		zap.L().Info("ws hub init")
	})
	return hub
}

// Run 运行Hub。
func (h *Hub) Run() {
	var mu sync.Mutex
	// 如果hub已经正在运行，则直接返回
	mu.Lock()
	defer mu.Unlock()
	if hub.running {
		zap.L().Info("ws hub run already")
		mu.Unlock()
		return
	}
	zap.L().Info("ws hub run start")
	hub.running = true // 将hub的running标志设置为true，表示正在运行
	mu.Unlock()
	for {
		select {
		case client := <-h.register: // 从h.register通道接收新连接的客户端
			register(client)
		case client := <-h.unregister: // 从h.unregister通道接收断开连接的客户端
			unregister(client)
		case message := <-h.broadcast: // 从h.broadcast通道接收广播消息
			broadcast(message)
		}
	}
}

// Broadcast 广播消息给所有客户端
func (h *Hub) Broadcast(innoPacket *pb.InnoPacket) {
	message, _ := proto.Marshal(innoPacket)
	h.broadcast <- message
}

// 注册客户端。
func register(c *Client) {
	hub.clients[c] = true // 将该客户端添加到h.clients映射中
	writeLog("register", c)
}

// 注销客户端。
func unregister(c *Client) {
	if _, ok := hub.clients[c]; ok { // 判断该客户端是否存在于h.clients映射中
		writeLog("unregister", c)
		delete(hub.clients, c) // 从h.clients映射中删除该客户端
		close(c.send)          // 关闭该客户端的send通道
	}
}

// 广播消息给所有客户端。
func broadcast(message []byte) {
	writeLog("broadcast", nil)
	for c := range hub.clients { // 遍历h.clients映射中的所有客户端
		select {
		case c.send <- message: // 将广播消息发送给该客户端
		default:
			unregister(c)
		}
	}
}

func writeLog(msg string, c *Client) {
	if c == nil {
		zap.L().Info("hub client "+msg,
			zap.Int("clients", len(hub.clients)),
			zap.String("all", strings.Join(utils.MapKeys(hub.ids), ",")),
			zap.String("boys", strings.Join(utils.MapKeys(hub.boys), ",")),
			zap.String("girls", strings.Join(utils.MapKeys(hub.girls), ",")),
		)
	} else {
		zap.L().Info("hub client "+msg,
			zap.Int("clients", len(hub.clients)),
			zap.String("uuid", c.uuid),
			zap.Int("ids", len(hub.ids[c.uuid])),
			zap.String("all", strings.Join(utils.MapKeys(hub.ids), ",")),
			zap.String("boys", strings.Join(utils.MapKeys(hub.boys), ",")),
			zap.String("girls", strings.Join(utils.MapKeys(hub.girls), ",")),
		)
	}
}
