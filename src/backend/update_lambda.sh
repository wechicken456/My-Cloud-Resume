#!/bin/bash
GOOS=linux GOARCH=amd64 go build -tags lambda.norpc -o bootstrap main.go
zip myFunction.zip bootstrap
aws lambda update-function-code \
    --function-name Cloud-Resume-API-Handler-Lambda \
    --zip-file fileb://myFunction.zip \
    --region us-east-1
