package bot

import (
    "log"
    "os"
    
    "github.com/yanzay/tbot"
    _ "github.com/joho/godotenv/autoload"

    "../models"
)

var vipPublicAPI = os.Getenv("MARKET_API_URL")


func Run() {
    // Here we initialize the db and then assign it to a global var of RedisClient
    // as defined in models.go
    models.RedisClient = models.InitializeDatabase()

    bot, err := tbot.NewServer(os.Getenv("TELEGRAM_TOKEN"))
    if err != nil {
        log.Fatal(err)
    }

    go bot.HandleFunc("/koin", listAllIdrCoins)
    go bot.HandleFunc("/harga {coin}", retrieveIdrTradeStat)

    // Help handlers
    go bot.HandleFunc("/help", customHelpHandler)
    go bot.HandleFunc("/tolong", customHelpHandler)
    go bot.HandleFunc("/list", customHelpHandler)

    // Set default handler if you want to process unmatched input
    bot.HandleDefault(unkownHandler)

    bot.ListenAndServe()
}
