.PHONY: run
run: build
	kubectl proxy -w static 
	
.PHONY: build
build:
	GOOS=js GOARCH=wasm go build -o static/main.wasm main.go