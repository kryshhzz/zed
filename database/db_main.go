package database

import (
  "gorm.io/gorm"
  _ "gorm.io/driver/sqlite"
  )
  
var DBConn *gorm.DB