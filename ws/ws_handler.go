package ws

import (
  "log"
  "fmt"

  "github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/kryshhzz/zed/user"
	"github.com/kryshhzz/zed/database"
	"github.com/dgrijalva/jwt-go"
)

var app *fiber.App

type WebSocketConn struct {
  conn *websocket.Conn
  usr user.User
}

var TotalCurrentConns []WebSocketConn

func Upgrader(app1 *fiber.App){
  app = app1
  app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})
}

	
	
func ConnHandle(c *websocket.Conn) {
  
  var usr user.User
  var err error

  cookie := c.Cookies("jwt")
  
  token,err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func (token *jwt.Token) (interface{}, error){
    return []byte("zedkryshhzzneon"),nil
  })
  
  if err != nil {
    return 
  }
  
  claims := token.Claims.(*jwt.StandardClaims)
  
  err = database.DBConn.Where("id = ?", claims.Issuer).First(&usr).Error

  if err != nil {
    return
  }
  
  curCon := WebSocketConn {
    conn : c,
    usr : usr,
  }
  
  TotalCurrentConns = append(TotalCurrentConns,curCon)
  
  fmt.Printf("%v \n", TotalCurrentConns)

  log.Println("Connection established")

	var (
		mt  int
		msg []byte
	)
	for {
		if mt, msg, err = c.ReadMessage(); err != nil {
			log.Println("read:", err)
			RemoveFromTotalCons(c)
			break
		}
		log.Printf("recv: %s", msg)

    WriteToRoom(c,mt,msg,usr)

	}
}



func RemoveFromTotalCons(c *websocket.Conn){
  
  var newTotalConns []WebSocketConn
  
  for _,xc := range TotalCurrentConns {
    if xc.conn != c {
      newTotalConns = append(newTotalConns,xc)
    }
  }
  
  TotalCurrentConns = newTotalConns
  
  fmt.Println("Removed ",c)
}


