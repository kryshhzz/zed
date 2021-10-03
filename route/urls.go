package route

import (
  "github.com/gofiber/fiber/v2"
  "github.com/kryshhzz/zed/ws"
  "github.com/kryshhzz/zed/user"
  "github.com/kryshhzz/zed/room"
  "github.com/gofiber/websocket/v2"
  )
  
func Route(app *fiber.App){
  
  // landing page view
  app.Get("/",landing_view)
  
  app.Get("/home",home_view)
  
  // web socket handling
  app.Get("/ws/:id",websocket.New(ws.ConnHandle))
  
  //userstuff
  app.Get("/users",user.GetUsers) // DEBUG
  
  
  app.Get("/register",register_view )
  app.Post("/register",user.Signup)
  
  
  app.Post("/login",user.Login)
  app.Get("/login",login_view) 
  
  app.Get("/logout",user.Logout)
  
  
  // rooms
  app.Get("/rooms",rooms_view)
  
  app.Get("/rooms/new",new_room_view)
  app.Post("/rooms/new",room.NewRoom)
  
  app.Get("/rooms/delete/:id",room.DeleteRoom)
  
  
  app.Get("/rooms/:id",room.EnterRoomView)
  
  app.Get("/rooms/auth/:id",RoomAuthView)
  app.Post("rooms/auth/:id",room.AuthenticateToRoom)
  
}
