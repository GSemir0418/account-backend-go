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
工作流！
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

清空数据表
DELETE FROM users;

go test -benchmem -bench "Crud" ./...

2997	    367429 ns/op	    2073 B/op	      47 allocs/op

执行了 2997 次，每次操作使用了 0.3ms，消耗 207KB 内存，开辟内存的次数 47 次
1 ns 是 10^-9 s

// 删除 table
DROP TABLE users, items;

// 查看表结构
SELECT table_name, column_name, data_type FROM information_schema.columns WHERE table_name='users';
// 查看索引
SELECT tablename, indexname, indexdef FROM pg_indexes WHERE schemaname='public' ORDER BY tablename, indexname;

1818	    607539 ns/op	   37587 B/op	     555 allocs/op

gorm 典型的使用空间弥补时间

提供数据库建表语句(CREATE/ALTER TABLE)和查询语句，自动编译为 Go 的 struct 和 func
sqlc 是编译器，可以将 sql 编译为 Go，类型安全
sqlc 不支持创建数据表以及数据迁移

选型网站 star-history.com

选择 golang-migrate 作为迁移工具
> https://github.com/golang-migrate/migrate#cli-usage
> as cli
> https://github.com/golang-migrate/migrate/tree/master/cmd/migrate
> as lib
> https://github.com/golang-migrate/migrate#use-in-your-go-project
> get started
> https://github.com/golang-migrate/migrate/blob/master/GETTING_STARTED.md

安装
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
或 brew install golang-migrate
创建迁移文件
migrate create -ext sql -dir config/migrations -seq create_users_table

运行迁移文件
migrate -database "postgres://gsemir:gsemir@localhost:5432/go_account_dev?sslmode=disable" \
-source "file://$(pwd)/config/migrations" up

封装命令行
记得把最后一个引入的 github 改为 file
"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"


创建迁移文件命令
go build; ./account db migrate:create add_email_to_users 
运行迁移文件命令
go build; ./account db migrate:up
回滚迁移命令
go build; ./account db migrate:down

记得手动同步 sqlc 的 schema.sql
每新增一个sql语句就要重新执行 sql generate

go install 安装的是命令行工具
默认安装路径是 $HOME/go/bin
go get 安装的是项目依赖库

使用 swaggo 生成API文档
> https://github.com/swaggo/swag/blob/master/README_zh-CN.md

通用注释写在 main 中
router 加上 swagger 的 router
controller 前加注释
swag init
go build; ./account server
访问 http://localhost:8080/swagger/index.html


验证码 API

1. 迁移数据库 建表
2. 发送验证码的配置
使用 gomail 库
> https://github.com/go-gomail/gomail
gomail 示例
> https://pkg.go.dev/gopkg.in/gomail.v2#example-package

3. 密钥管理 将邮箱授权码保存到环境变量

.zshrc.local
export EMAIL_SMTP_PWD='xxxxxx'
source ~/.zshrc.local
此时再执行发送命令即可

go 通过 os.Getenv("EMAIL_SMTP_PWD") 获取到密码即可

但当协同开发时设置环境变量会很繁琐，且数据类型仅支持字符串，需要手动转换
有没有一个文件，里面有项目的密码配置，但不提交到 github 以保证安全性呢

4. 使用 viper 密钥管理
  
> Ruby 牛逼 一个 master.key 就可以解决了
> viper
> https://github.com/spf13/viper

main里面设置配置文件路径并读取
viper.config.json.example 可以提交
viper.config.json 不能提交
为了方便测试环境读取viper配置
将 viper 的配置文件统一存放到 $HOME/.account 项目目录下
使用绝对路径读取配置文件即可
封装Viper的读取逻辑 在 router 中引入将

创建第一个路由
拿到json请求体数据
1. 声明结构体实例
2. 将结构体实例绑定到上下文
3. 此时结构体实例body就是请求体了

测试构造请求体数据
strings.NewReader(`{"email": "test@qq.com"}`)

测试环境不会运行 main.go
且运行了也无法正确读取 viper 的配置（路径）
所以将 viper 的配置文件统一存放到 $HOME/.account 项目目录下
使用绝对路径读取配置文件即可
封装Viper的读取逻辑 在 router 中引入

用 sqlc 保存验证码到数据库
其实使用redis做这个需求更合适

对 sqlc 的配置进行调整 queries: "config/query.sql" => queries: "config/queries"
 拆分 query.sql，按字段分类编写 sql 查询语句

router 连接数据库

为了方便测试 封装 CountValidationCodes 方法
在测试环境下注意要先 router.New 连接数据库，然后再查询

生成真随机验证码
使用 crypto/rand 库
注意数字切片转换为字符编码的逻辑

使用 mail hog 简化邮件测试
> https://github.com/mailhog/MailHog
>
go install github.com/mailhog/MailHog@latest
无需在打开目标邮箱去看邮件是否正常发送了
需要将读取到的邮件配置的email.smtp.host和email.smtp.port覆盖(放到 router.New 之后)
MailHog 命令启动本地邮件服务器 8025 前端 1025 是后端
此外该服务器还提供了api用于读取全部收到的邮件，用于在测试代码中调用 这样连打开网页检查都不用了


登录 api
创建路由
获取请求体（定义struct 绑定json请求体）
定义查询validation_code语句：查询email和code且没有被使用过