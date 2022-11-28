package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// fmt.Printf("event.PathParamenter %v", request.PathParameters["ip"])
	fmt.Println(request)
	pathParam, found := request.PathParameters["proxy"]

	if found {
		queryIP, _ := url.QueryUnescape(pathParam)

		// Return a response with a 200 OK status
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Body:       "Param received:" + queryIP,
		}, nil
	}

	return clientError(http.StatusNotFound)
}

func clientError(status int) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       http.StatusText(status),
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
