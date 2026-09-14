DOCKER_IMAGE := jiro4989/textimg

textimg: parser/grammar.peg.go *.go */*.go
	go fmt ./...
	go build

parser/grammar.peg.go: parser/grammar.peg
	peg parser/grammar.peg

.PHONY: help
help: ## ドキュメントのヘルプを表示する。
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: test
test: textimg ## テストコードを実行する
	go test -cover ./...

.PHONY: docker-build
docker-build: ## Docker イメージをビルドする
	docker compose build

.PHONY: docker-test
docker-test: ## Docker 環境で go test を実行する
	docker compose run --rm base go test -tags docker -cover ./...

.PHONY: docker-push
docker-push: ## Docker イメージをビルドして DockerHub にイメージを Push する
	@if [ "$(TAG)" = "" ]; then echo "[ERR] TAG 変数は必須です"; exit 1; fi
	docker build --no-cache -t $(DOCKER_IMAGE):$(TAG) .
	docker push $(DOCKER_IMAGE):$(TAG)

.PHONY: setup-tools
setup-tools: ## 開発時に使うツールをインストールする
	go install github.com/pointlander/peg@latest
