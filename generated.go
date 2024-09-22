package main

import (
    "context"
    "encoding/base64"
    "fmt"
    "github.com/aws/aws-lambda-go/events"
    "github.com/aws/aws-lambda-go/lambda"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3"
    "os"
    "time"
)

var s3Client *s3.S3
var bucketName = os.Getenv("BUCKET_NAME")

func init() {
    sess := session.Must(session.NewSession())
    s3Client = s3.New(sess)
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    imageBase64 := request.Body

    // Decode the base64 image
    imageData, err := base64.StdEncoding.DecodeString(imageBase64)
    if err != nil {
        return events.APIGatewayProxyResponse{StatusCode: 400, Body: fmt.Sprintf("Error decoding image: %s", err)}, nil
    }

    // Create a unique file name
    fileName := fmt.Sprintf("image_%d.jpg", time.Now().Unix())

    // Upload to S3
    _, err = s3Client.PutObject(&s3.PutObjectInput{
        Bucket: aws.String(bucketName),
        Key:    aws.String(fileName),
        Body:   aws.ReadSeekCloser(bytes.NewReader(imageData)),
        ContentType: aws.String("image/jpeg"),
    })

    if err != nil {
        return events.APIGatewayProxyResponse{StatusCode: 500, Body: fmt.Sprintf("Failed to upload image: %s", err)}, nil
    }

    return events.APIGatewayProxyResponse{
        StatusCode: 200,
        Body:       fmt.Sprintf("Image uploaded successfully: %s", fileName),
    }, nil
}

func main() {
    lambda.Start(handler)
}
