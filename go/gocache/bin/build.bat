@echo off 

setlocal 
echo %~dp0

cd/d %~dp0
cd ..

go build -o bin/ cmd/multi-node.go

echo build successfully
endlocal

pause