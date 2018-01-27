package main

import (
    // "log"
    // "github.com/aws/aws-lambda-go/events"
    // "github.com/aws/aws-lambda-go/lambda"
)


// func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
//     log.Printf("Processing Lambda request %s\n", request.RequestContext.RequestID)

//     if len(request.Body) >= 1 {
//         log.Printf("Body request received: %s", request.Body)
//     }

//     Run()

//     return events.APIGatewayProxyResponse{
//         Body:       "Hello " + request.Body,
//         StatusCode: 200,
//     }, nil
// }


func main() {
    Run()
    // lambda.Start(Handler)
}
