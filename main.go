package main

import (
  "log"
  "os"
  
	"github.com/gofiber/template/django"
  "github.com/gofiber/fiber/v2"	
  "github.com/gofiber/fiber/v2/middleware/cors"	
  "github.com/gofiber/fiber/v2/middleware/logger"	
  "gorm.io/gorm"
  "gorm.io/driver/sqlite"
  
	"github.com/kryshhzz/zed/ws"
	"github.com/kryshhzz/zed/route"
	"github.com/kryshhzz/zed/room"
	"github.com/kryshhzz/zed/user"
	"github.com/kryshhzz/zed/database"
	"github.com/kryshhzz/zed/message"
)

func initDb (){
  var err error

  database.DBConn, err = gorm.Open(sqlite.Open("db.sqlite3"),&gorm.Config{})
  
  if err != nil {
    log.Println("😭 Database didn't connect")
    panic(err)
  }
  log.Println("👍 Database successfully connected ")
  
  database.DBConn.AutoMigrate(&user.User{})
  database.DBConn.AutoMigrate(&message.Message{})
  database.DBConn.AutoMigrate(&room.Room{})
  log.Println("👍 Database successfully migrated ")
}

func main() {
  
  // initialising database
  initDb()

  engine := django.New("./templates", ".html")

  app := fiber.New(fiber.Config{
      Views: engine,
  })  
  
  // using logger 
  app.Use(logger.New())
	
	// using cors
  app.Use(cors.New(cors.Config{
    AllowCredentials : true,
  }))
	
	// sending the app instance for websocket
	ws.Upgrader(app)
	
	// sending the app instance for routing
	route.Route(app)
	
  // listening
	//log.Fatal(app.Listen(":6969"))
	log.Fatal(app.Listen(":"+os.Getenv("PORT")))
}

