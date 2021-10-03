package room

import (
  "fmt"
  "gorm.io/gorm"
  "github.com/kryshhzz/zed/message"
  "github.com/kryshhzz/zed/database"
  "github.com/kryshhzz/zed/user"
  
  "github.com/gofiber/fiber/v2"
  )
  
type Room struct {
  gorm.Model
  
  Title string `json:"title" gorm:"unique"`
  SecKey string `json:"-"`
  Msgs []message.Message `gorm:"foreignKey:RoomRefer" `
  Members []user.User `gorm:"many2many:room_members"`
  Private bool `json:"private"`
  
  AdminId uint
  Admin user.User `gorm:"foreignKey:AdminId"`
}


func NewRoom(ctx*fiber.Ctx) error{
  var newroom Room
  var dummy Room
  
  // authentication of user
  err,usr := user.UserInfo(ctx)
  if err != nil {
    return err
  }
  if usr.ID == 0 {
    return ctx.Redirect("/login")
  }
  
  // retrieving data from the user
  var data map[string]string
  
  if err := ctx.BodyParser(&data); err != nil {
    return err
  }

  newroom.Title = data["title"]
  newroom.SecKey = data["seckey"]
  
  if data["private"] == "yes"{
    newroom.Private = true;
  }else{
    newroom.Private = false;
  }

  r := database.DBConn.Where("Title = ?",newroom.Title).First(&dummy)
  
  if r.Error == nil {
    return ctx.JSON(fiber.Map{
      "error" : "Title already taken",
    })
  }
  
  newroom.Members = []user.User{usr}
  newroom.Admin = usr
  newroom.AdminId = usr.ID
  
  //DEBUG
  fmt.Println("Shit",len(newroom.Members))
  
  r = database.DBConn.Create(&newroom)
  if r.Error != nil {
    fmt.Println(r.Error)
  }
  
  r = database.DBConn.Save(&newroom)
  if r.Error != nil {
    fmt.Println(r.Error)
  }
  
  return ctx.JSON(fiber.Map{
      "roomCreated" : true,
    })
}



func DeleteRoom(ctx*fiber.Ctx) error{
  
  id := ctx.Params("id")
  var room Room

  database.DBConn.First(&room,id)
  if room.Title == "" {
    return ctx.Status(500).SendString("No room found with the id ")
  }
  
  database.DBConn.Delete(&room)
  return ctx.SendString("Deleted room :: ")
}