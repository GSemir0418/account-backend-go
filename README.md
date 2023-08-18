go mod init account

https://gin-gonic.com/docs/

https://gin-gonic.com/docs/quickstart/

go get -u github.com/gin-gonic/gin

go build;./account

go mod tidy 自动安装依赖包,删除多余依赖

go test ./... 递归执行全部测试文件
go test ./test/...

docker 体积越来越大怎么办
首先开启所有用到的容器
执行 
docker system prune
docker volume prune