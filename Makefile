-include .env

# ---------------- Database Start ---------------------------------------
jet-gen:
	jet -dsn=${BUYBETTER_DEV_SUPABASE_DSN} -path=./.jetgen

# ---------------- Database End -----------------------------------------

# ---------------- Golang Utils Start ---------------------------------------
test:
	go test ./... -v

lint:
	golangci-lint run ./...

# ---------------- Golang Utils End -----------------------------------------

# ---------------- Terraform Start ---------------------------------------
terraform-init:
	terraform -chdir=terraform/ init -backend-config=backend-config.tfvars
terraform-plan:
	terraform -chdir=terraform/ plan
terraform-apply:
	terraform -chdir=terraform/ apply -auto-approve
terraform-destroy:
	terraform -chdir=terraform/ destroy -auto-approve

# ---------------- Terraform End -----------------------------------------

# ---------------- AWS lambda Start -----------------------------------
FUNCTION_NAME := transform-category
FUNCTION_PATH := ./bin/function.zip
AWS_REGION := us-east-2

QUEUE_URL := ${SQS_QUEUE_URL}
MESSAGE_BODY := '{"version":"v1-lambda","shuffle":true,"train_ratio":60,"validate_ratio":20,"test_ratio":20}'

build-lambda:
	GOOS=linux GOARCH=amd64 go build -o ./bin/bootstrap ./cmd/lambda/main.go
	zip -j $(FUNCTION_PATH) ./bin/bootstrap

deploy-lambda:
	aws lambda update-function-code \
        --function-name $(FUNCTION_NAME) \
        --zip-file fileb://$(FUNCTION_PATH) \
        --region $(AWS_REGION)

publish-lambda:
	aws lambda publish-version \
        --function-name $(FUNCTION_NAME) \
        --region $(AWS_REGION)

sent-message:
	aws sqs send-message \
		--queue-url $(QUEUE_URL) \
		--message-body $(MESSAGE_BODY) \
		--region $(AWS_REGION)
# ---------------- AWS lambda End -----------------------------------