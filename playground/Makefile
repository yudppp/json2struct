build:
	GOOS=js GOARCH=wasm go build -o playground.wasm .
	mv ./playground.wasm ../docs
	cp "$(shell go env GOROOT)/lib/wasm/wasm_exec.js" ../docs/