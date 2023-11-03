BINARY_NAME=toolbox

dep:
	@go mod tidy && go fmt

build:
	@GOARCH=arm64 GOOS=darwin go build -o bin/${BINARY_NAME}_darwin_arm64 main.go

build-race:
	@GOARCH=arm64 GOOS=darwin go build -race -o bin/${BINARY_NAME}_darwin_arm64 main.go


run: build
	@echo $(MAKECMDGOALS)
	./bin/${BINARY_NAME}_darwin_arm64

dev: 
	@CompileDaemon --command="./bin/daemon-build" -build="go build -o bin/daemon-build" -color=true -graceful-kill=true

clean:
	go clean
	rm bin/${BINARY_NAME}_linux_arm64