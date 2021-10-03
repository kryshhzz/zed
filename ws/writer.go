package ws

import (
  "log"
  "time"
  "strconv"
  "strings"
  "encoding/json"
  "github.com/gofiber/websocket/v2"
  
  "github.com/kryshhzz/zed/database"
  "github.com/kryshhzz/zed/message"
  "github.com/kryshhzz/zed/user"
  )


func WriteToRoom(c *websocket.Conn,mt int,msg []byte,usr user.User){
  
  var msgProto message.Message 
  
  room_num := c.Params("id")
  
  for _,xc := range TotalCurrentConns {
    if xc.conn.Params("id") == room_num {
      
      var jsonMap map[string]interface{}
      json.Unmarshal(msg, &jsonMap)

      jsonMap["SenderID"] = usr.ID

      // sending message to everyone 
      if _, ok := jsonMap["message"]; ok {
        
        // checking reply id
        repId,err := strconv.Atoi(jsonMap["replyID"].(string))
        if err != nil {
          log.Println("ReplyID number not an integer",err)
          RemoveFromTotalCons(xc.conn)
        	break
        }
        
        //adding time
        layout1 := "Mon Jan 2 2006 15:04:05 MST "
        layout2 := "15:04"
        timeStr := time.Now()
        timeStr.String()
        log.Println(timeStr.Format(layout1))

        jsonMap["time"] = timeStr.Format(layout2)
        
        //mt = len(jsonMapStr)
        room_num_int,err := strconv.Atoi(room_num)
        if err != nil {
          log.Println("Room number not an integer",err)
      		RemoveFromTotalCons(xc.conn)
      		break
        }
        
        msgProto.Sender = usr
        msgProto.SenderID = usr.ID
        msgProto.Text = jsonMap["message"].(string)
        msgProto.RoomRefer = uint(room_num_int)
        
        var repToMsg message.Message
        
        // setting the message as reply in the db
        if repId != 0{
          
          r := database.DBConn.Preload("Sender").First(&repToMsg,repId)
          
            if r.Error != nil {
              log.Println("No message to reply !",r.Error.Error())
          		RemoveFromTotalCons(xc.conn)
          		break
            }
            
            msgProto.RepToID = uint(repId)
            msgProto.ReplyTo = &repToMsg
          }

        // sending the reply text to the user / ws
        if repId != 0 {
          jsonMap["ReplyTo"] = repToMsg.Sender.Username
          jsonMap["ReplyText"] = repToMsg.Text
        }else{
          jsonMap["ReplyTo"] = ""
          jsonMap["ReplyText"] = ""
        }
        
        r := database.DBConn.Create(&msgProto)
        if r.Error != nil {
          if strings.Contains(r.Error.Error(),"UNIQUE constraint failed: messages.id") != true{
            log.Println("Inserting erorr !",r.Error.Error())
        		RemoveFromTotalCons(xc.conn)
        		break
          }
        }else{
          log.Println("saved the msg")
        }
        
        // sending the msg's id and sender
        jsonMap["actID"] = msgProto.ID
        jsonMap["SenderUsername"] = usr.Username
        
        jsonMapStr,err := json.Marshal(jsonMap)
        log.Println(err)
        
        if err := xc.conn.WriteMessage(mt, []byte(jsonMapStr)); err != nil {
      		log.Println("write:", err)
      		RemoveFromTotalCons(xc.conn)
      		break
      	}
     
      // sending the typing status to everyone except the user
      }else{
        if xc.conn != c {
          if err := xc.conn.WriteMessage(mt,msg); err != nil {
        		log.Println("write:", err)
        		RemoveFromTotalCons(xc.conn)
        		break
      	  }
        }
      }
      
      
    }
  }
}