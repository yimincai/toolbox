BINARY_NAME=toolbox

dep:
	go mod tidy && go mod vendor && go fmt

build:
	GOARCH=amd64 GOOS=linux go build -o dist/${BINARY_NAME}_linux_amd64 main.go
	GOARCH=arm64 GOOS=darwin go build -o dist/${BINARY_NAME}_darwin_arm64 main.go

run:
	go run main.go

# dev: 
# 	CompileDaemon --command="./dist/daemon-build" -build="go build -o dist/daemon-build" -color=true -graceful-kill=true

clean:
	go clean
	rm dist/${BINARY_NAME}_linux_amd64

# image:
# 	docker buildx build --no-cache --platform linux/amd64 -t yimincai/${BINARY_NAME} --load .