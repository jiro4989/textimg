APPNAME := $(shell basename `pwd`)
VERSION := $(shell grep Version internal/global/version.go | grep -Eo '`.*$$' | tr -d '`')
SRCS := $(shell find . -name "*.go" -type f )
XBUILD_TARGETS := \
	-os="windows linux darwin" \
	-arch="386 amd64" 
DIST_DIR := dist
README := README.*

.PHONY: help
help: ## ドキュメントのヘルプを表示する。
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: test
test: ## テストコードを実行する
	go test -cover ./...

.PHONY: clean
clean: ## バイナリ、配布物ディレクトリを削除する
	-rm -rf bin
	-rm -rf $(DIST_DIR)

.PHONY: docker-build
docker-build: ## Dockerイメージをビルドする
	docker-compose build --no-cache

.PHONY: docker-test
docker-test: ## Docker環境でgo testを実行する
	docker-compose run --rm base go test -tags ondocker -cover ./...

.PHONY: docker-push
docker-push: ## DockerHubにイメージをPushする
	docker push jiro4989/textimg
