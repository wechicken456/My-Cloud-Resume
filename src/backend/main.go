package main

import (
	"context"
	"log"
	appConfig "main/internal/config"
	"main/internal/handlers"
	"main/internal/service"
	"main/internal/storage"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	ses "github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

func main() {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("Couldn't load AWS config: %s", err)
	}

	appCfg := appConfig.Load()

	// Initialize AWS clients
	dynamoClient := dynamodb.NewFromConfig(cfg)
	sesClient := ses.NewFromConfig(cfg)
	snsClient := sns.NewFromConfig(cfg)

	// Initialize storage
	store := storage.New(dynamoClient, appCfg.DynamoDBTable, appCfg.SessionTable)

	// Initialize services
	counterService := service.NewVisitorService(store)
	likesService := service.NewLikesService(store)
	contactService := service.NewContactService(appCfg)
	notificationService := service.NewNotificationService(sesClient, snsClient, appCfg)

	// Initialize handler
	apiHandler := handlers.NewAPIHandler(counterService, likesService, contactService, notificationService)

	lambda.Start(apiHandler.HandleRequest)
}
