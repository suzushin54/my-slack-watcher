.PHONY: install build deploy help

## Install necessary packages
install:
	npm install @types/node@20 @aws-cdk/aws-lambda @aws-cdk/aws-apigateway @aws-cdk/aws-iam

## Build the project
build:
	npm run build

## Deploy the project
deploy:
	cdk deploy

## Display this help message
help:
	@echo "Usage: make [target]"
	@echo
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
