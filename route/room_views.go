package route

import (
  _"fmt"
  "github.com/gofiber/fiber/v2"
  "github.com/kryshhzz/zed/user"
  "github.com/kryshhzz/zed/room"
  "github.com/kryshhzz/zed/database"
  )

func rooms_view(ctx *fiber.Ctx) error{
  
  // authentication of user
  err,usr := user.UserInfo(ctx)
  if err != nil {
    return err
  }
  if usr.ID == 0 {
    return ctx.Redirect("/login")
  }
  
  // returning rooms 
  var rooms []room.Room
  database.DBConn.Preload("Msgs").Preload("Members").Preload("Admin").Find(&rooms)
  //database.DBConn.Find(&rooms)
  
  return ctx.Render("room/rooms",fiber.Map{
      "rooms" : rooms,
    }) 
}


func new_room_view(ctx *fiber.Ctx) error{
  // authentication of user
  err,usr := user.UserInfo(ctx)
  if err != nil {
    return err
  }
  if usr.ID == 0 {
    return ctx.Redirect("/login")
  }
  
  return ctx.Render("room/new_room",fiber.Map{})
}


func RoomAuthView(ctx *fiber.Ctx) error{
  err,usr := user.UserInfo(ctx)
  if err != nil {
    return err
  }
  if usr.ID == 0 {
    return ctx.Redirect("/login")
  }
  
  return ctx.Render("room/room_auth",fiber.Map{
    "roomId" : ctx.Params("id"),
  })
}