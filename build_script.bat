@echo off
setlocal
set GOOS=linux
set GOARCH=amd64
go build -o build/github-hook main.go
copy .\action.sh .\build\
endlocal