package main

import (
    "fmt"
    "time"

    "github.com/aws/aws-lambda-go/lambda"

    "./bot"
    "./lambdaparser"
    "./models"
)


var currentTime = time.Now().UTC()


func handler(request lambdaparser.Request) (lambdaparser.Response, error) {
    requestBody, err := lambdaparser.ProcessRequest(request.Body)
    if err != nil { 
        return lambdaparser.Response{
            Message: "Invalid request received",
            Ok:      false,
        }, nil        
    }

    message := requestBody.Message

    // Here we initialize the db and then assign it to a global var of RedisClient
    // as defined in models.go
    models.RedisClient = models.InitializeDatabase()
    models.StoreUser(message.From.Username, message.Chat.FirstName, message.From.ID)
    models.UpdateMarketData(message.Text, currentTime.Format("200601021504"))

    bot.Run()

    return lambdaparser.Response{
        Message: fmt.Sprintf("Processed request ID %f", requestBody.UpdateID),
        Ok:      true,
    }, nil    
}


func main() {
    lambda.Start(handler)
}