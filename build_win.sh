rm -rf ./e32_33config.exe
clear

#export WORK=./work
export CGO_ENABLED=1
export CC=x86_64-w64-mingw32-gcc
#export CXX=x86_64-w64-mingw32-g++
export GOOS=windows
export GOARCH=amd64
go build -x -v -ldflags "-s -w" -o e32_33config.exe
# -x -v -ldflags "-s -w"
