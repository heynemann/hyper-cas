.PHONY: serve route sync-simple build

serve:
		@go run main.go serve --config hyper-cas.yaml

route:
		@go run main.go route --config hyper-cas.yaml

sync-simple:
		@go run main.go sync --config ./hyper-cas.yaml --label master ./fixtures/simple

build:
		@mkdir -p build
		@GOOS=darwin go build -o ./build/hyper-cas-mac main.go
		@GOOS=linux go build -o ./build/hyper-cas main.go
		@GOOS=windows go build -o ./build/hyper-cas-win main.go

test:
	@go test ./...
