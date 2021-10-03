package room

import (
  "fmt"
  "strconv"
  "github.com/gofiber/fiber/v2"
  "github.com/kryshhzz/zed/database"
  "github.com/kryshhzz/zed/user"
  
 )


func EnterRoomView(ctx *fiber.Ctx) error{
  
  err,usr := user.UserInfo(ctx)
  if err != nil {
    return err
  }
  if usr.ID == 0 {
    return ctx.Redirect("/login")
  }
  
  var room_to_enter Room
  var inRoom bool = false

  roomId := ctx.Params("id")
  
  _,err = strconv.Atoi(roomId)
  
  if err != nil {
    return ctx.Status(500).SendString("Got string as ID expected an integer")
  }
  
  r := database.DBConn.Preload("Members").Preload("Msgs").Preload("Msgs.Sender").Preload("Msgs.ReplyTo").Preload("Msgs.ReplyTo.Sender").Preload("Admin").First(&room_to_enter,roomId)
  if r.Error != nil {
    return ctx.Status(500).SendString(r.Error.Error())
  }else{
    fmt.Println("Rooom to enter ",room_to_enter.ID)
  }
  
  if room_to_enter.ID == 0 {
    return ctx.JSON("Room doesn't exist")
  }
  
  
  fmt.Println("Private room :: ",room_to_enter.Private)
  
  // checking if the person is a member of the room
  for _,mem := range room_to_enter.Members{
    fmt.Println(mem.Username)
    if mem.ID == usr.ID {
      inRoom = true;
    }
  }
  fmt.Println("inroom : ",inRoom)
  
  // room is public
  if room_to_enter.Private == false {
    
    if inRoom == false {
      database.DBConn.Model(&room_to_enter).Association("Members").Append([]user.User{usr})
    }
    
    return ctx.Render("room/index",fiber.Map{
      "roomId" : room_to_enter.ID,
      "msgs" : room_to_enter.Msgs,
      "usr" : usr,
    })
    
  }else{
    // private room 
    if inRoom == false {
      authUrl := "/rooms/auth/"
      return ctx.Redirect(authUrl+roomId)
    }else{
      return ctx.Render("room/index",fiber.Map{
        "roomId" : room_to_enter.ID,
        "msgs" : room_to_enter.Msgs,
        "usr" : usr,
      })
    }
  }
}


func AuthenticateToRoom(ctx *fiber.Ctx) error{
  
  err,usr := user.UserInfo(ctx)
  if err != nil {
    return err
  }
  if usr.ID == 0 {
    return ctx.Redirect("/login")
  }
  
  var data map[string]string

  var roomId = ctx.Params("id")
  var room_to_enter Room

  if err = ctx.BodyParser(&data); err != nil {
    return err
  }
  
  // checking if the id is string
  _,err = strconv.Atoi(roomId)
  if err != nil {
    return ctx.Status(500).SendString("Got string as ID expected an integer")
  }
  
  r := database.DBConn.Preload("Members").Preload("Msgs").Preload("Msgs.Sender").Preload("Msgs.ReplyTo").Preload("Msgs.ReplyTo.Sender").Preload("Admin").First(&room_to_enter,roomId)
  if r.Error != nil {
    return ctx.Status(500).SendString(r.Error.Error())
  }
  
  if room_to_enter.ID == 0 {
    return ctx.JSON(fiber.Map{
      "error" : "Room doesn't exist!",
    })
  }
  
  if data["seckey"] == room_to_enter.SecKey{
    database.DBConn.Model(&room_to_enter).Association("Members").Append([]user.User{usr})
    return ctx.JSON(fiber.Map{
      "joinedRoom" : true,
    })
  }else{
    fmt.Println("Expected ",room_to_enter.SecKey," Received ",data["seckey"])
    return ctx.JSON(fiber.Map{
      "error" : "Wrong key !",
    })
  }
}