@echo off 

setlocal 
echo %~dp0

del multi-node.exe

cd/d %~dp0
cd ..

go build -o bin/ cmd/multi-node.go

echo build successfully
endlocal

::延迟 5秒，启动测试
timeout /T 5

start cmd /k "cd/d %~dp0 && multi-node.exe -port=8001"
start cmd /k "cd/d %~dp0 && multi-node.exe -port=8002"
start cmd /k "cd/d %~dp0 && multi-node.exe -port=8003"

start cmd /k "cd/d %~dp0 && multi-node.exe -api"


::延迟 5秒，启动测试
timeout /T 5

echo start to tests

echo curl 127.0.0.1:9001/api?key=k1
curl 127.0.0.1:9001/api?key=k1
echo.

echo curl 127.0.0.1:9001/api?key=k2
curl 127.0.0.1:9001/api?key=k2
echo.

echo curl 127.0.0.1:9001/api?key=k3
curl 127.0.0.1:9001/api?key=k3
echo.

echo done
pause