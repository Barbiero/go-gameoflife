build:
	GOOS=js GOARCH=wasm go build -o main.wasm

serve: build
	npx serve -p 3000

run: serve
