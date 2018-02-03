package main

import (
    "fmt"
    "github.com/aws/aws-lambda-go/lambda"

    "./bot"
    "./parser"
)


func handler(request parser.Request) (parser.Response, error) {
    requestBody := parser.ProcessRequest(request.Body)

    bot.Run()

    return parser.Response{
        Message: fmt.Sprintf("Processed request ID %f", requestBody.UpdateID),
        Ok:      true,
    }, nil    
}


func main() {
    lambda.Start(handler)
}
