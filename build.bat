set GOOS=js
set GOARCH=wasm
go build -o main.wasm main.go
pause