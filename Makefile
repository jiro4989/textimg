.PHONY: help
help: ## ドキュメントのヘルプを表示する。
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: test
test: ## テストコードを実行する
	go test -cover ./...

.PHONY: docker-build
docker-build: ## Dockerイメージをビルドする
	docker-compose build

.PHONY: docker-test
docker-test: ## Docker環境でgo testを実行する
	docker-compose run --rm base go test -tags ondocker -cover ./...

.PHONY: docker-push
docker-push: ## DockerHubにイメージをPushする
	docker push jiro4989/textimg
