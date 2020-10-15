.PHONY: serve route sync-simple build

serve:
		@go run main.go serve --config hyper-cas.yaml

docker-serve: docker
	@docker run -it -p 2485:2485 -v "/tmp/hyper-cas/sites:/app/sites" -v "/tmp/hyper-cas/storage:/app/storage" hyper-cas:latest

route:
		@docker run -v "`pwd`/route-nginx.conf:/etc/nginx/nginx.conf" -v "/tmp/hyper-cas:/app" -v "/tmp/hyper-cas/sites:/etc/nginx/conf.d" -p 80:80 -it nginx:latest

route-daemon:
		@docker run -v "`pwd`/route-nginx.conf:/etc/nginx/nginx.conf" -v "/tmp/hyper-cas:/app" -v "/tmp/hyper-cas/sites:/etc/nginx/conf.d" -p 80:80 -d nginx:latest

sync-simple:
		@go run main.go sync --config ./hyper-cas.yaml --label master ./fixtures/simple

build:
		@mkdir -p build
		@GOOS=darwin CGO_ENABLED=0 go build -o ./build/hyper-cas-mac main.go
		@GOOS=linux CGO_ENABLED=0 go build -o ./build/hyper-cas main.go
		@GOOS=windows CGO_ENABLED=0 go build -o ./build/hyper-cas-win main.go

test:
	@go test ./...

docker:
	@docker build -t hyper-cas .

push-image: docker
	@docker tag hyper-cas:latest vtexcom/hyper-cas:latest
	@docker push vtexcom/hyper-cas:latest
