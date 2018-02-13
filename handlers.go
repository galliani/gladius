// package main

// import (
//     "log"
//     "os"
//     "strings"
//     "fmt"

//     // _ "github.com/joho/godotenv/autoload"

//     // "../models"    
// )


// func Run() {
//     bot, err := tbot.NewServer(os.Getenv("TELEGRAM_TOKEN"))
//     if err != nil {
//         log.Fatal(err)
//     }

//     go bot.HandleFunc("/koin", listAllIdrCoins)
//     go bot.HandleFunc("/harga {coin}", retrieveIdrTradeStat)

//     // Help handlers
//     go bot.HandleFunc("/help", customHelpHandler)
//     go bot.HandleFunc("/tolong", customHelpHandler)
//     go bot.HandleFunc("/list", customHelpHandler)

//     // Set default handler if you want to process unmatched input
//     bot.HandleDefault(unkownHandler)

//     bot.ListenAndServe()
// }



// func listAllIdrCoins(message *tbot.Message) {
//     coins := []string{"Bitcoin [btc]", "Bitcoin Cash [bch]", "Bitcoin Gold [btg]", "Litecoin [ltc]", "Ethereum [eth]", "Ethereum Classic [etc]", "Ripple [xrp]", "Lumens [xlm]", "Waves [waves]", "NXT [nxt]", "ZCoin [xzc]"}

//     message.Reply("Berikut adalah daftar koin yang nilai tukarnya dengan rupiah dapat kamu tanya ke saya")
//     for _, coin := range coins {
//         message.Reply(coin)
//     }
//     message.Reply("==========")
// }


// func retrieveIdrTradeStat(message *tbot.Message) {
//     coinTicker := strings.ToLower(message.Vars["coin"])

//     message.Replyf("Berikut info mengenai aktivitas perdagangan IDR-%s", strings.ToUpper(coinTicker))

//     stat := models.RetrieveMarketStats(coinTicker)
//     stat.DisplayAsMoney()

//     message.Replyf("Harga Terakhir: %s", stat.Last)
//     message.Replyf("Harga Beli #1: %s", stat.Buy)
//     message.Replyf("Harga Jual #1: %s", stat.Sell)
//     message.Replyf("Harga Tertinggi (24 jam): %s", stat.High)
//     message.Replyf("Harga Terendah (24 jam): %s", stat.Low)
//     message.Reply("==========")
// }
