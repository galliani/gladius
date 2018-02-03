package main

import (
    "fmt"
    "github.com/aws/aws-lambda-go/lambda"
)


func handler(request Request) (Response, error) {
    requestBody := processLambdaRequest(request.Body)

    run()

    return Response{
        Message: fmt.Sprintf("Processed request ID %f", requestBody.UpdateID),
        Ok:      true,
    }, nil    
}


func main() {
    lambda.Start(handler)
}
