package message

import (
  "github.com/jinzhu/gorm"
  "github.com/kryshhzz/zed/user"
  )
  
type Message struct {
  gorm.Model
  
  Text string `json:text"`
  SenderID uint
  Sender user.User `gorm:"foreignKey:SenderID"`
  
  RepToID uint 
  ReplyTo *Message `gorm:"foreignKey:RepToID"`
  
  RoomRefer uint
}