set GOOS=js
set GOARCH=wasm
go build -o website/main.wasm src/main.go
pause