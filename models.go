package main
// import (
//     "os"
//     // "github.com/jinzhu/gorm"
//     // _ "github.com/jinzhu/gorm/dialects/postgres"    
// )

// var (
//     // DBCon is the connection handle
//     // for the database
//     DbCon *gorm.DB
// )

// The Stat struct should be able to decode a JSON object like this:
// {
//     "ticker": {
//         "high": "16500000",
//         "low": "13611000",
//         "vol_eth": "4453.35873651",
//         "vol_idr": "66475714301",
//         "last": "14301000",
//         "buy": "14301000",
//         "sell": "14375000",
//         "server_time": 1515107696
//     }
// }
type Stat struct {
    Ticker struct {
        High            string `json:"high"`
        Low             string `json:"low"`
        VolEth          string `json:"vol_eth"`
        VolIdr          string `json:"vol_idr"`
        Last            string `json:"last"`
        Buy             string `json:"buy"`
        Sell            string `json:"sell"`
        ServerTime      int    `json:"server_time"`
    }
}

// type User struct {
//   gorm.Model
//   FirstName string
//   LastName string
//   UserName string
//   TelegramUserUid int
// }

// func InitializeDatabase() (db *gorm.DB) {
//     db, err := gorm.Open("postgres", os.Getenv("POSTGRES_ADDR"))
//     if err != nil {
//         panic(err)
//     }    
    
//     db.AutoMigrate(&User{})

//     return db
// }