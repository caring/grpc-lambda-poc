package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/caring/grpc-lambda-poc/pb"
	"google.golang.org/grpc"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

const serviceAddress = "localhost:3000"

// Gateway is our lambda handler invoked by the `lambda.Start` function call
func Gateway(ctx context.Context) (Response, error) {
	conn, err := grpc.Dial(serviceAddress, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
		return Response{StatusCode: 500}, err
	}
	defer conn.Close()

	c := pb.NewCallscoringClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	serviceResponse, err := c.GetValidAttributes(ctx, &pb.GetAttributesRequest{})
	if err != nil {
		log.Fatalf("could not fetch valid attributes: %v", err)
		return Response{StatusCode: 500}, err
	}

	var buf bytes.Buffer

	body, err := json.Marshal(serviceResponse)
	if err != nil {
		log.Fatalf("could not marshal service response: %v", err)
		return Response{StatusCode: 500}, err
	}

	json.HTMLEscape(&buf, body)

	response := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type":           "application/json",
			"X-MyCompany-Func-Reply": "gateway-handler",
		},
	}

	return response, nil
}

func main() {
	lambda.Start(Gateway)
}
