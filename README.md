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

三种数据库选型
database/sql
> https://pkg.go.dev/database/sql
优点：
1. 官方包，得到Go社区的支持与维护
2. 轻量级，不会引入额外的依赖
3. 良好的跨数据兼容性
缺点：
1. 缺少高级功能，如查询构建器、关联、迁移等
2. 手动管理sql语句，易出错且难以维护

gorm
> https://gorm.io/docs/
优点：
1. 提供了丰富的功能，如查询构建器、关联、迁移等
2. 支持自动迁移数据库
3. 提供了更高级的查询接口，减少了手动编写sql语句的需求
4. 良好的文档和社区支持
缺点：
1. 依赖更多的外部库，可能导致项目臃肿
2. 抽象层可能导致性能损失
3. 可能需要更多的学习成本

sqlc
> https://docs.sqlc.dev/en/latest/tutorials/getting-started-postgresql.html
优点：
1. 将 sql 查询转换为类型安全的go代码，提高代码的可读性和安全性
2. 通过生成代码，可以减少手动编写sql语句的错误
3. 支持postgresql和mysql
4. 自动生成代码，易与维护
缺点：
1. 没有支持所有类型的数据库
2. 可能需要对sql语句进行调整以生成正确的go代码
3. 自动生成的代码可能难以理解和调试
4. 功能有限，缺少一些高级功能，如关联、迁移等

基于任务的学习模式
创建数据表
数据迁移（给user添加一个字段）
基本增删改查
错误处理
性能测试

使用 cobra 创建命令行程序
> https://github.com/spf13/cobra
实现类似 rails 的 bin/rails db:create 等命令
go get -u github.com/spf13/cobra@latest

go build;./account server

连接数据库
构建连接字符串
> https://www.connectionstrings.com/

不会写sql怎么办
> https://devdocs.io/postgresql~14/sql-createtable

go build;./account db create