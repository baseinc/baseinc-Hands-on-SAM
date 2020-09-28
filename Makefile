.PHONY: build test local-run local-up dev-deploy

build:
	sam build

test:
	@cd ./hello-world/ && \
	go test -v ./...

# ローカルでlambdaを実行する
local-run: build
	sam local invoke \
		--parameter-overrides 'ENV=local AppName=Hello'

# apiとして起動する場合
local-up: build
	sam local start-api \
		--parameter-overrides 'ENV=local AppName=Hello'

# デプロイする
dev-deploy: build
	sam deploy \
	--s3-bucket sam-app-stack \
	--parameter-overrides 'ENV=dev AppName=Hello' \
	--force-upload \
	--no-fail-on-empty-changeset