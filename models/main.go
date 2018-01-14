package models
import(
    "fmt"
    "os"
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/postgres"    
)

var (
    // DBCon is the connection handle
    // for the database
    DbCon *gorm.DB
)

type User struct {
  gorm.Model
  FirstName string
  LastName string
  UserName string
  TelegramUserUid int
}

func InitializeDatabase() (db *gorm.DB) {
    db, err := gorm.Open("postgres", os.Getenv("POSTGRES_ADDR"))
    if err != nil {
        panic(err)
    }    
    
    db.AutoMigrate(&User{})

    fmt.Println("Successfully connected!")
    return db
}