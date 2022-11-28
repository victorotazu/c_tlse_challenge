package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type extapi struct {
	IpApiKey string
}

var httpClient = &http.Client{}

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	pathParam, found := request.PathParameters["proxy"]

	if found {
		queryIP, _ := url.QueryUnescape(pathParam)
		ip := net.ParseIP(queryIP)

		if ip == nil {
			clientError(http.StatusBadRequest)
		}

		fmt.Printf("Parsed IP: %v\n", ip.String())
		// Return a response with a 200 OK status
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Body:       geolocationdata(ip.String()),
		}, nil
	}

	return clientError(http.StatusNotFound)
}

func geolocationdata(ip string) string {
	var result extapi
	json.Unmarshal([]byte(os.Getenv("extapikey")), &result)

	url := "https://ipapi.co/" + ip + "/json/?key=" + result.IpApiKey

	resp, err := httpClient.Get(url)
	if err != nil {
		// It would be interesting to implement a fallback
		// in case the external api is down. MaxMind db perhaps
		serverError(http.StatusInternalServerError)
	}

	defer resp.Body.Close()

	geodata, err := io.ReadAll(resp.Body)
	if err != nil {
		serverError(http.StatusInternalServerError)
	}

	return string(geodata)
}

func clientError(status int) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       http.StatusText(status),
	}, nil
}

func serverError(status int) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       http.StatusText(status),
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
