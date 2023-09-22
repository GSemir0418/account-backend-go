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
防止重复连接数据库（并发）
不断完善测试，注意执行流程及结构体和json之间的转换
注意函数参数类型

生成jwt
> https://github.com/golang-jwt/jwt
> 示例 https://pkg.go.dev/github.com/golang-jwt/jwt/v5#example-New-Hmac

go get -u github.com/golang-jwt/jwt/v5

创建 jwt_helper package
定义两个方法，分别用于生成加密用的密钥，以及根据用户id生成jwt字符串
在测试代码中使用fmt.Println打印 resBody 中的 JWT
在构造测试背景时报错使用log.Fatalln
在处理测试错误时使用t.Error
执行测试的代码中加 -v 可以显示详情
go test -timeout 30s -run '^TestSession$' account/test/controller_test -v

继续生成 user_id
有则读取，没有则创建
避免重复创建测试用户失败，将删除表的操作以及初始化测试的代码抽离出来
省略 jwt 解密过程，登录后返回userId

将HMAC密钥生成后保存到本地环境变量，避免重复生成（因为解密要用）
os.WriteFile保存到本地，viper环境变量中存文件路径即可/Users/gsemir/.account/jwt/hmac.key
通过命令行工具生成 jwt 密钥并保存


展示测试覆盖率
go test -coverprofile=coverage/coverage.out ./...
生成测试覆盖率html
go tool cover -html=coverage/coverage.out -o coverage/coverage.html
由于测试文件和测试代码不在同一个目录下，导致测试覆盖率无法正确展示
将测试文件复制过去后，统一包名为 controller，会报错，由于我们在测试中引入了 router 而router又引用了controller，导致循环引用
router 引用 controller 是必然的
那么就要想办法切断测试中对 router 的引用
复制初始化gin服务器、初始化路由、加载配置和连接数据库的操作到 setupTest函数中即可。

创建 controller 接口
为了统一管理 ctrler 开发与引入
各模块的 ctrler 就是一个结构体 包含 ctrler 接口规定的 Create Destroy Update Find RegisterRoutes
go 中没有继承与实现的概念，只要这个结构体有接口要求的方法就可以作为这个接口的实现

自动编写各模块 controller 结构体方法：
> https://github.com/josharian/impl
> go install github.com/josharian/impl@latest
生成代码到控制台供复制
impl 'ctrl *SessionController' account/internal/controller.Controller
也可以用vscode 插件 ctrl shift p go stubs
输入: ctrl *SessionController account/internal/controller.Controller

总结
r.POST("/api/v1/session", controller.CreateSession)
=>
每个 controller 自己注册路由，分配路由
v1 := rg.Group("v1")
v1.POST("session", ctrl.Create)
而 router 只负责循环调用每个 ctrler 的 register 方法


优化测试工作流
目前我们的TDD流程为
写测试 => 写代码 => 开启 MailHog => 执行测试 => 生成覆盖率 => 生成覆盖率 html => 开启http服务打开html => 根据覆盖率调整测试代码
可以将这些抽离为一个控制台命令
os/exec 包提供了运行脚本的功能
exec.Command(...).Start()和Run()的区别
start会执行这条命令，但不会等待其结束，直接执行下一行代码
这样能够将测试结果可视化，方便补充测试用例

me controller
get请求，获取用户信息
读取请求头的jwt信息，解码出id，去数据库中查询id对应的用户，返回此用户
go stub 插件，输入 ctrl *MeController account/internal/controller.Controller
注册路由
写测试
抽离测试初始化方法为setup_test_case.go
1. 读取环境变量配置（email，secret等）
2. 连接数据库实例 q
3. 初始化 gin 服务器实例 r
4. 初始化上下文参数 c
5. 删除 user 表，方便测试
6. 返回清理函数，开发者自行选择执行
共有三个测试用例公共变量 r q c

断言逻辑、
jwt 的 *jwt.Token 是一个结构体类型的指针，表示一个 jwt 对象的值

断点调试
go install -v github.com/go-delve/delve/cmd/dlv@latest
点 debug test 也行 vscode 自动安装 使用默认配置即可

重新生成 API 文档
https://github.com/swaggo/swag/blob/master/README_zh-CN.md#%E5%BF%AB%E9%80%9F%E5%BC%80%E5%A7%8B
controller 中添加注释
创建文档
swag init && go build . && ./account sever
为了方便获取到返回值的类型（api文档注释最多读取两层包的数据），封装一个 api 的包，里面注明各接口返回值类型（取sqlc生成的结构体类型即可）

更改 json 输出格式
目前我们接口的输出是由 sqlc 自动生成的，例如 created_at 字段，目前的输出是 createdAt，我们要统一json字段风格为snake
可以在 struct 后添加``json:"created_at"``注释，但 sqlc 不允许我们修改自动生成的文件，只能改 sqlc的配置
自动添加 json 注释 emit_json_tags: true
指定 json 风格 json_tags_case_style: snake
snake pascal camel
