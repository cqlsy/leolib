package leowebecho

import (
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/cqlsy/leolib/leoutil"
	"net/http"
)

func AddSocketClient(path string, manager *SocketManager) {
	Echo.GET(path, manager.WsHandler)
	go manager.Start()
}

// Client is a websocket client
type Client struct {
	ID          string          // 区分每一次的消息
	socket      *websocket.Conn // socket的链接实例
	send        chan []byte     // 发送消息的通信
	isSendClose bool            // 是否关闭
}

// SocketManager is a websocket manager
type SocketManager struct {
	Clients      map[string]*Client // 所有的链接信息存储在这里
	register     chan *Client
	unregister   chan *Client
	generateID   func(c echo.Context) string      // 生成ID
	onGetMessage func(client *Client, msg []byte) // 当收到客户端的消息时调用该函数
	log          func(errStr string)              // 错误信息传出
}

func NewManager(geId func(c echo.Context) string,
	onGetMessage func(client *Client, msg []byte)) *SocketManager {
	manage := &SocketManager{
		register:   make(chan *Client),       // 这里定义用户链接注册
		unregister: make(chan *Client),       // 用户离开了，现在需要保存
		Clients:    make(map[string]*Client), // 链接信息
	}
	if geId == nil {
		geId = func(c echo.Context) string {
			return leoutil.RandString(32, "socket")
		}
	}
	manage.generateID = geId
	manage.onGetMessage = onGetMessage
	return manage
}

// log打印
func (manager *SocketManager) InitLog(f func(str string)) {
	manager.log = f
}

// start is  项目运行前, 协程开启start -> go Manager.start()
func (manager *SocketManager) Start() {
	for {
		select {
		case conn := <-manager.register:
			// 链接成功，将当前的链接存入缓存
			manager.Clients[conn.ID] = conn
			//conn.SendMessage([]byte("client "))
		case conn := <-manager.unregister:
			// 当链接失败了，我们需要将链接移除
			if _, ok := manager.Clients[conn.ID]; ok {
				conn.isSendClose = true
				close(conn.send)
				delete(manager.Clients, conn.ID)
			}
		}
	}
}

// 发送消息,发送之前判断通道是否关闭了
func (c *Client) SendMessage(message []byte) {
	if c.isSendClose {
		return
	}
	c.send <- message
}

func (c *Client) read(manager *SocketManager) {
	defer func() {
		manager.unregister <- c
		_ = c.socket.Close()
	}()
	for {
		_, message, err := c.socket.ReadMessage()
		if err != nil {
			manager.unregister <- c
			_ = c.socket.Close()
			break
		}
		// 把接收到的消息移交出去
		if manager.onGetMessage != nil {
			manager.onGetMessage(c, message)
		}
	}
}

// 开通写的通信操作
func (c *Client) write(manager *SocketManager) {
	defer func() {
		_ = c.socket.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				err := c.socket.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil && manager.log != nil {
					manager.log(err.Error())
				}
				return
			}
			err := c.socket.WriteMessage(websocket.TextMessage, message)
			if err != nil && manager.log != nil {
				manager.log(err.Error())
			}
		}
	}
}

//TestHandler socket 连接 中间件 作用:升级协议,用户验证,自定义信息等
func (manager *SocketManager) WsHandler(c echo.Context) error {
	conn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(c.Response().Writer, c.Request(), nil)
	if err != nil {
		http.NotFound(c.Response().Writer, c.Request())
		return err
	}
	// 可以添加用户信息验证
	client := &Client{
		ID:          manager.generateID(c),
		socket:      conn,
		send:        make(chan []byte),
		isSendClose: false,
	}
	manager.register <- client
	go client.read(manager)
	go client.write(manager)
	return nil
}
