package main
import (
    "log"
    "os"
    
    "github.com/yanzay/tbot"
    _ "github.com/joho/godotenv/autoload"

    // Custom lib packages
    "./models"
    "./handlers"
)


func Run() {
    // Here we initialize the db and then assign it to a global var of RedisClient
    // as defined in models.go
    models.RedisClient = models.InitializeDatabase()

    bot, err := tbot.NewServer(os.Getenv("TELEGRAM_TOKEN"))
    if err != nil {
        log.Fatal(err)
    }

    go bot.HandleFunc("/koin", handlers.ListAllIdrCoins)
    go bot.HandleFunc("/harga {coin}", handlers.RetrieveIdrTradeStat)

    // Help handlers
    go bot.HandleFunc("/help", handlers.CustomHelpHandler)
    go bot.HandleFunc("/tolong", handlers.CustomHelpHandler)
    go bot.HandleFunc("/list", handlers.CustomHelpHandler)

    // Set default handler if you want to process unmatched input
    bot.HandleDefault(handlers.UnkownHandler)

    bot.ListenAndServe()
}
