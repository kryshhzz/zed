package route

import (
  "github.com/gofiber/fiber/v2"
  "github.com/kryshhzz/zed/user"
  )
  
func landing_view(ctx *fiber.Ctx) error{
  return ctx.Render("index", fiber.Map{
        //"student_name": "Steve",
    })
}

func home_view(ctx *fiber.Ctx) error {
  err,usr := user.UserInfo(ctx)
  if err != nil {
    return err
  }
  
  if usr.ID == 0 {
    return ctx.Redirect("/login")
  }
  
  return ctx.Render("home",fiber.Map{
    "username" : usr.Username,
  })
}

func login_view(ctx *fiber.Ctx )error {
  return ctx.Render("auth/login",fiber.Map{})
}

func register_view(ctx *fiber.Ctx) error {
  return ctx.Render("auth/register",fiber.Map{})
}