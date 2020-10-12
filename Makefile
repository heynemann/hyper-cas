.PHONY: serve route sync-simple build

serve:
		@go run main.go serve --config hyper-cas.yaml

route:
		# @go run main.go route --config hyper-cas.yaml
		@docker run -v "`pwd`/route-nginx.conf:/etc/nginx/nginx.conf" -v "/tmp/hyper-cas:/app" -v "/tmp/hyper-cas/sites:/etc/nginx/conf.d" -p 8000:80 -it nginx:latest

route-daemon:
		@docker run -v "`pwd`/route-nginx.conf:/etc/nginx/nginx.conf" -v "/tmp/hyper-cas:/app" -v "/tmp/hyper-cas/sites:/etc/nginx/conf.d" -p 8000:80 -d nginx:latest

sync-simple:
		@go run main.go sync --config ./hyper-cas.yaml --label master ./fixtures/simple

build:
		@mkdir -p build
		@GOOS=darwin go build -o ./build/hyper-cas-mac main.go
		@GOOS=linux go build -o ./build/hyper-cas main.go
		@GOOS=windows go build -o ./build/hyper-cas-win main.go

test:
	@go test ./...
