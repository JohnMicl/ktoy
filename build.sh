go env -w GOPROXY=https://goproxy.cn,direct
export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:./release/x86_64/linux
go build -o ./release/x86_64/linux/ktoy
./release/x86_64/linux/ktoy
