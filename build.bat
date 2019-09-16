set GOOS=js
set GOARCH=wasm
go build -o website/main.wasm main.go
pause