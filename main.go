package main

import (
    "net/http"
    "time"
    "os"
    "log"
    "fmt"

    "gopkg.in/telegram-bot-api.v4"
    _ "github.com/joho/godotenv/autoload"
    "github.com/aws/aws-lambda-go/lambda"

    "./lambdaparser"
    "./models"
)


var currentTime = time.Now().UTC()


func main() {
    isProduction := os.Getenv("IS_PRODUCTION")

    if isProduction == "true" {
        lambda.Start(handler)
    } else {
        coreExecutor()
    }
}


func handler(request lambdaparser.Request) (lambdaparser.Response, error) {
    requestBody, err := lambdaparser.ProcessRequest(request.Body)
    if err != nil { 
        return lambdaparser.Response{
            Message: "Invalid request received",
            Ok:      false,
        }, nil        
    }

    coreExecutor()

    return lambdaparser.Response{
        Message: fmt.Sprintf("Processed request ID %f", requestBody.UpdateID),
        Ok:      true,
    }, nil    
}


func coreExecutor() {
    models.RedisClient = models.InitializeDatabase()

    bot, updates := startBotServer()

    for update := range updates {
        log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

        // Here we initialize the db and then assign it to a global var of RedisClient
        // as defined in models.go
        models.StoreUser(update.Message.From.UserName, update.Message.From.FirstName, update.Message.From.ID)
        // models.UpdateMarketData(message.Text, currentTime.Format("200601021504"))

        replies := determineUpdateResponse(update)
        for _, reply := range replies {
            bot.Send(reply)
        }
    }
}


func determineUpdateResponse(update tgbotapi.Update) []tgbotapi.Chattable {
    switch update.Message.Text {
    case "/help", "/tolong", "/list":
        return customHelpHandler(update)
    default:
        return unkownHandler(update)
    }
}


func startBotServer() (*tgbotapi.BotAPI, tgbotapi.UpdatesChannel) {
    bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_TOKEN"))
    if err != nil {
        log.Fatal(err)
    }

    _, err = bot.SetWebhook(tgbotapi.NewWebhook(os.Getenv("BOT_WEBHOOK_URL")))
    if err != nil {
        log.Fatal(err)
    }

    updates := bot.ListenForWebhook("/")
    port    := os.Getenv("BOT_WEBHOOK_PORT")

    go http.ListenAndServe(":" + port, nil)

    log.Printf("Bot is up and running on port %s", port)

    return bot, updates
}


// Handlers
func unkownHandler(update tgbotapi.Update) []tgbotapi.Chattable {
    t := make([]tgbotapi.Chattable, 2)

    chatID := update.Message.Chat.ID
    t[0] = tgbotapi.NewMessage(chatID, fmt.Sprintf("Maaf %s, saya tidak mengerti permintaanmu barusan", update.Message.From.FirstName))
    t[1] = tgbotapi.NewMessage(chatID, "Kamu dapat melihat perintah yang tersedia dengan mengetik /help")

    return t
}


func customHelpHandler(update tgbotapi.Update) []tgbotapi.Chattable {
    t := make([]tgbotapi.Chattable, 3)

    chatID := update.Message.Chat.ID
    t[0] = tgbotapi.NewMessage(chatID, "Inilah daftar perintah yang tersedia untuk kamu:")
    t[1] = tgbotapi.NewMessage(chatID, "/koin - Untuk mengetahui semua koin yang dapat ditanyakan harganya")
    t[2] = tgbotapi.NewMessage(chatID, "/harga NAMA_SINGKAT_KOIN || contohnya, ketik: /harga btc")

    return t
}