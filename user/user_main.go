package user

import (
  "errors"
  "time"
  "strconv"
  
  "github.com/gofiber/fiber/v2"
  "gorm.io/gorm"
  "github.com/kryshhzz/zed/database"
  "golang.org/x/crypto/bcrypt"
  
	"github.com/dgrijalva/jwt-go"
  )


type User struct {
  gorm.Model
  Username string `json:"username" gorm:"unique"`
  Password []byte `json:"-"`
  Email string `json:"email" gorm:"unique"`
}

func Hash(password string) ([]byte, error) {
  return bcrypt.GenerateFromPassword([]byte(password), 14)
}


func GetUsers(ctx*fiber.Ctx) error{
  var usrs []User
  database.DBConn.Find(&usrs)

  return ctx.JSON(usrs)
}


func Signup(ctx *fiber.Ctx) error{
  var usr User
  var dummy User
  
  var data map[string]string
  
  if err := ctx.BodyParser(&data); err != nil {
    return err
  }
  
  hashed_pwd,err := Hash(data["password"])
  
  if err != nil {
    return ctx.Status(500).SendString("Hashing error")
  }
  
  usr.Username = data["username"]
  usr.Password = hashed_pwd
  usr.Email = data["email"]
  
  
  r := database.DBConn.Where("username = ?", data["username"]).First(&dummy)
  
  if r.Error == nil {
    return ctx.JSON(fiber.Map{
      "error" : "Username taken",
    })
  }else{
    r = database.DBConn.Where("email = ?", data["email"]).First(&dummy)
    if r.Error == nil {
      return ctx.JSON(fiber.Map{
        "error" : "Email already registered",
      })
    }
  }
  
  database.DBConn.Create(&usr)
  
  return ctx.JSON(fiber.Map{
      "userCreated" : true,
    })
}


func Login(ctx *fiber.Ctx) error{
  
  var data map[string]string
  
  if err := ctx.BodyParser(&data); err != nil {
    return err
  }
  
  usrname := data["username"]
  pwd := data["password"]
 
  var usr User
  
  database.DBConn.Where("username = ?", usrname).First(&usr)
  
  if usr.ID == 0 {
    ctx.Status(fiber.StatusNotFound)
    return ctx.JSON(fiber.Map{
        "error" : "No such user",
      })
  }
  
  if err := bcrypt.CompareHashAndPassword(usr.Password,[]byte(pwd)); err != nil {
    ctx.Status(fiber.StatusBadRequest)
    return ctx.JSON(fiber.Map{
        "error" : "Invalid Credentials",
      })
  }
	
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256,jwt.StandardClaims{
	    
	    Issuer : strconv.Itoa(int(usr.ID)),
	    ExpiresAt : time.Now().Add(time.Hour * 24).Unix(), // one day
	    
	  })
	  
	token,err := claims.SignedString([]byte("zedkryshhzzneon"))
  
  if err != nil {
    ctx.Status(fiber.StatusInternalServerError)
    return ctx.JSON(fiber.Map{
        "error" : "could not login",
      })
  }
  
  cookie := fiber.Cookie {
    Name : "jwt",
    Value : token,
    Expires : time.Now().Add(time.Hour * 24),
    HTTPOnly : true,
  }
  
  ctx.Cookie(&cookie)
  
	return ctx.JSON(fiber.Map{
	  "loggedIn" : true,
	})
}

func UserInfo(ctx *fiber.Ctx) (error,User) {
  
  var usr User
  var err error

  cookie := ctx.Cookies("jwt")
  
  token,err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func (token *jwt.Token) (interface{}, error){
    return []byte("zedkryshhzzneon"),nil
  })
  
  if err != nil {
    return ctx.Redirect("/login"),usr
  }
  
  claims := token.Claims.(*jwt.StandardClaims)
  
  err = database.DBConn.Where("id = ?", claims.Issuer).First(&usr).Error
  
  return err,usr
}


func Logout(ctx *fiber.Ctx) error{
  
  cookie := fiber.Cookie{
    Name : "jwt",
    Value : "",
    Expires : time.Now().Add(-time.Hour),
    HTTPOnly : true,
  }
  
  ctx.Cookie(&cookie)
  
  /*
  return ctx.JSON(fiber.Map{
      "loggedIn" : false,
    })*/
  
  return ctx.Redirect("/login")
}


func DeleteUser(id int) error {
  var usr User
  
  database.DBConn.First(&usr,id)
  
  if usr.Username == "" {
    return errors.New("No user found with the id ")
  }
  database.DBConn.Delete(&usr)
  return errors.New("")
}
