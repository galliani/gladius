package main

import (
    "fmt"
    "github.com/aws/aws-lambda-go/lambda"

    "./bot"
    "./lambdaparser"
    "./models"
)


func handler(request lambdaparser.Request) (lambdaparser.Response, error) {
    requestBody, err := lambdaparser.ProcessRequest(request.Body)
    if err != nil { panic(err) }

    // Here we initialize the db and then assign it to a global var of RedisClient
    // as defined in models.go
    models.RedisClient = models.InitializeDatabase()

    message := requestBody.Message
    models.StoreUser(message.From.Username, message.Chat.FirstName, message.From.ID)
    models.UpdateMarketData(message.Text)

    bot.Run()

    return lambdaparser.Response{
        Message: fmt.Sprintf("Processed request ID %f", requestBody.UpdateID),
        Ok:      true,
    }, nil    
}


func main() {
    lambda.Start(handler)
}