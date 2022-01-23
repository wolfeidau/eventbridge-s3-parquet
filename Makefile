APPNAME := eventbridge-s3-parquet
STAGE ?= dev
BRANCH ?= master

default: clean build lint test

deploy: default deploy-storage deploy-s3-events

include build/go.mk

.PHONY: archive
archive:
	@echo "--- build an archive"
	@cd dist && zip -X -9 ./handler.zip *-lambda

.PHONY: deploy-storage
deploy-storage:
	@echo "--- deploy stack $(APPNAME)-$(STAGE)-$(BRANCH)-storage"

	@sam deploy \
		--no-fail-on-empty-changeset \
		--template-file sam/app/storage.yaml \
		--capabilities CAPABILITY_IAM \
		--tags "environment=$(STAGE)" "branch=$(BRANCH)" "service=$(APPNAME)" \
		--stack-name $(APPNAME)-$(STAGE)-$(BRANCH)-storage \
		--parameter-overrides AppName=$(APPNAME) Stage=$(STAGE) Branch=$(BRANCH)

.PHONY: deploy-s3-events
deploy-s3-events:
	$(eval SAM_BUCKET := $(shell aws ssm get-parameter --name '/config/$(STAGE)/$(BRANCH)/deploy_bucket' --query 'Parameter.Value' --output text))
	$(eval DATA_BUCKET_NAME := $(shell aws ssm get-parameter --name '/config/$(STAGE)/$(BRANCH)/$(APPNAME)/data_bucket_name' --query 'Parameter.Value' --output text))
	
	@sam deploy \
		--no-fail-on-empty-changeset \
		--template-file sam/app/s3-events.yaml \
		--capabilities CAPABILITY_IAM \
		--s3-bucket $(SAM_BUCKET) \
		--s3-prefix sam/$(GIT_HASH) \
		--tags "environment=$(STAGE)" "branch=$(BRANCH)" "service=$(APPNAME)" \
		--stack-name $(APPNAME)-$(STAGE)-$(BRANCH)-s3-events \
		--parameter-overrides AppName=$(APPNAME) Stage=$(STAGE) Branch=$(BRANCH) \
			DataBucketName=$(DATA_BUCKET_NAME) 

.PHONY: logs-s3-events
logs-s3-events:
	@sam logs \
		--stack-name $(APPNAME)-$(STAGE)-$(BRANCH)-s3-events -t