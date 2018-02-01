package main

import (
    "log"
    "fmt"
    "encoding/json"
    "github.com/aws/aws-lambda-go/lambda"
)

// Taken from "github.com/aws/aws-lambda-go/lambda/events/apigw.go"
// Commented out for documentation purposes, like this
// is the proper struct to wrap the request sent to the lambda
type Request struct {
    // Resource              string                        `json:"resource"` // The resource path defined in API Gateway
    // Path                  string                        `json:"path"`     // The url path for the caller
    // HTTPMethod            string                        `json:"httpMethod"`
    // Headers               map[string]string             `json:"headers"`
    // QueryStringParameters map[string]string             `json:"queryStringParameters"`
    // PathParameters        map[string]string             `json:"pathParameters"`
    // StageVariables        map[string]string             `json:"stageVariables"`
    // RequestContext        APIGatewayProxyRequestContext `json:"requestContext"`
    // IsBase64Encoded       bool                          `json:"isBase64Encoded,omitempty"`
    Body                     string                        `json:"body"`
}

// {
//     "update_id": "xxxx",
//     "message": {
//         "message_id": 4347,
//         "from": {
//             "id": 23423114,
//             "is_bot": false,
//             "first_name": "User 1",
//             "username": "telegramuser",
//             "language_code": "en-US"
//         },
//         "chat": {
//             "id": 25838532,
//             "first_name": "User 1",
//             "username": "telegramuser",
//             "type": "private"
//         },
//         "date": 1516875119,
//         "text": "/idr eth",
//         "entities": [
//             {
//                 "offset": 0,
//                 "length": 4,
//                 "type": "bot_command"
//             }
//         ]
//     }
// }

type RequestBody struct {
    UpdateID      int         `json:"update_id"`
    Message       *Message    `json:"message"`
}

type Message struct {
    ID        int       `json:"message_id"`
    From      *From     `json:"from"`
    Chat      *Chat     `json:"from"`
    Date      int       `json:"date"`
    Text      string    `json:"text"`
}

type From struct {
    ID              int       `json:"id"`
    IsBot           bool      `json:"is_bot"`
    Username        string    `json:"username"`
    LanguageCode    string    `json:"language_code"`
}

type Chat struct {
    ID          int         `json:"id"`
    FirstName   string      `json:"first_name"`
    Username    string      `json:"username"`
    Type        string      `json:"type"`
}

type Response struct {
  Message string `json:"message"`
  Ok      bool   `json:"ok"`
}


func Handler(request Request) (Response, error) {
    log.Printf("Processing Lambda request %s\n", request.Body)

    requestBody := RequestBody{}
    err := json.Unmarshal([]byte(request.Body), &requestBody)
    if err != nil {
        panic(err)
    }


    // log.Printf("Body request received: %s", request.Body.Message)

    Run()

    return Response{
        Message: fmt.Sprintf("Processed request ID %f", requestBody.UpdateID),
        Ok:      true,
    }, nil    
}


func main() {
    lambda.Start(Handler)
}
