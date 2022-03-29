textimg: parser/grammer.peg.go *.go */*.go
	go fmt ./...
	go build

parser/grammer.peg.go: parser/grammer.peg
	peg parser/grammer.peg

.PHONY: help
help: ## ドキュメントのヘルプを表示する。
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: test
test: textimg ## テストコードを実行する
	go test -cover ./...

.PHONY: docker-build
docker-build: ## Dockerイメージをビルドする
	docker-compose build

.PHONY: docker-test
docker-test: ## Docker環境でgo testを実行する
	docker-compose run --rm base go test -tags docker -cover ./...

.PHONY: docker-push
docker-push: ## DockerHubにイメージをPushする
	docker push jiro4989/textimg

.PHONY: setup-tools
setup-tools: ## 開発時に使うツールをインストールする
	go install github.com/pointlander/peg@latest
	wget https://raw.githubusercontent.com/ekalinin/github-markdown-toc/master/gh-md-toc
	sudo install -m 0755 ./gh-md-toc /usr/local/bin/
	-rm -f gh-md-toc
