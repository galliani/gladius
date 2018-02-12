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

    // "./bot"
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

        txtMsg := tgbotapi.NewMessage(update.Message.Chat.ID, "Yo")
        bot.Send(txtMsg)
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
