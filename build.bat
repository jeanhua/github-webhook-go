@echo off
setlocal
set GOOS=windows
set GOARCH=amd64
go build -o build/github-hook_windows.exe
endlocal

@echo off
setlocal
set GOOS=linux
set GOARCH=amd64
go build -o build/github-hook_amd64_linux
endlocal

@echo off
setlocal
set GOOS=darwin
set GOARCH=amd64
go build -o build/github-hook_amd64_macos
endlocal