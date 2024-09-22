package main

import (
    "bytes"
    "context"
    "encoding/base64"
	"fmt"
    "os"
    "time"
    "github.com/aws/aws-lambda-go/events"
    "github.com/aws/aws-lambda-go/lambda"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3"
)

var s3Client *s3.S3
var bucketName = os.Getenv("kc-wedding-image-bucket-storage")

func init() {
    sess := session.Must(session.NewSession())
    s3Client = s3.New(sess)
}

// uploadHandler handles the file upload
func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    imageBase64 := request.Body

    // Decode the base64 image
    imageData, err := base64.StdEncoding.DecodeString(imageBase64)
    if err != nil {
        return events.APIGatewayProxyResponse{StatusCode: 400, Body: fmt.Sprintf("Error decoding image: %s", err)}, nil
    }


    // Limit file size to 100MB
    if len(imageData) > 100*1024*1024 {
        return events.APIGatewayProxyResponse{StatusCode: 400, Body: "File is too large"}, nil
    }

    fileType := request.Headers["Content-Type"]
    if fileType == "" {
        return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Content-Type header is missing"}, nil
    }

    // Restrict file types to images only
    allowedTypes := map[string]bool{
        "image/gif": true,
        "image/heif": true,
        "image/jpeg": true,
        "image/png":  true,
        "image/webp": true,
    }

    if !allowedTypes[fileType] {
        return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Invalid file type"}, nil
    }

    // TODO: Hash images to prevent repeats
    
        
    // TODO: Compress files


    // TODO: Reformat images to webp for size

    // Create a unique file name using the input file type
    fileExtension := GetExtension(fileType)
    fileName := fmt.Sprintf("image_%d%s", time.Now().Unix(), fileExtension)

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

// GetExtension maps MIME types to file extensions
func GetExtension(fileType string) string {
    switch fileType {
    case "image/gif":
        return ".gif"
    case "image/heif":
        return ".heif"
    case "image/jpeg":
        return ".jpg"
    case "image/png":
        return ".png"
    case "image/webp":
        return ".webp"
    default:
        return ".jpg" // Fallback if type is unknown
    }
}

func main() {
    lambda.Start(handler)
}

