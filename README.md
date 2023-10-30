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

gin 中间件
将接口鉴权流程提取出来，即 GetMe 的逻辑 提取作为中间件
注意返回值和参数
修改 router.go
修改 测试的 setup
运行测试

开发 items api
1. 创建路由和控制器
ctrl shift p go stubs
输入: ctrl *ItemController account/internal/controller.Controller
2. 注册路由
3. 数据库迁移
1 创建迁移文件
go build; ./account db migrate:create create_items_table
2 编写迁移命令（创建枚举类型及数据表）
因为每次迁移最好只做一件事，使用事务，确保迁移过程的多次关联操作不会出错，或者引起数据库异常（类型与表格必须同时存在）
pgsql 支持定义枚举类型(kind)，以及数组类型（tag_ids）,mysql中只能用关联表来做
3 执行迁移命令
go build; ./account db migrate:up
如果出错了，且不能使用回滚来解决，那么手动修改数据库的 schema_migrations 表格，将数据设置为上次的版本，dirty为false
update schema_migrations set version=3,dirty=false;
手动修改数据库或迁移文件错误后，重新执行迁移即可
4. sqlc 生成
1 写model
schema.sql 写建表语句
2 写方法
queries/items.sql 写创建方法
3 sqlc generate
5. 写 Create API
1 声明 body 结构体，将请求上下文绑定到 body 中
2 获取当前登录用户
3 items 表新增
4 返回 item
6. 单元测试
1 初始化 q r c
2 注册路由
3 初始化 w 构造请求 req
4 创建用户 添加请求头权限字段
5 r 发起请求
6 断言

文档格式
注意要紧贴控制器方法 中间不要有空行
且注释中间也不能有空行 要使用空的//相连
// CreateItem
//
//	@Summary	创建账目
//	@Accept		json
//	@Produce	json
//	@Param		amount		body		int		true	"金额（单位：分）"	example(100)
//	@Param		kind		body		queries.Kind	true	"类型"		example(expenses)
//	@Param		tag_ids		body		[]string		true	"标签ID列表"	example([1,2,3])
//	@Success	200			{object}	api.CreateItemResponse
//	@Failure	401			{string}	string	无效的JWT
//	@Failure	422			{string}	string	参数错误
//	@Router		/api/v1/items [post]

支持 JWT 测试：
在main中
//	@securityDefinitions.apiKey	Bearer
//	@in							header
//	@name						Authorization
在需要权限的文档加上：能自动带上 Authorizition 请求头
//	@Security	Bearer

之后我们在测试前获取到validationCode 然后访问session接口，将返回的 jwt 保存在 文档上方 Authorize 的位置（别忘 Bearer ） 需要权限的接口就会自动携带这个全局 jwt 来请求了


bug 修复！花了一下午！！！
debug 配置
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Go Debug",
      "type": "go",
      "request": "launch",
      "mode": "auto",
      "program": "${workspaceFolder}",
      "args": [
        "server"
      ]
    }
  ]
}
注意中间件的声明位置，要在所有路由之前Use

item 分页
- 取出参数使用c.Request.URL.Query() 而不是绑定json请求体了
query := c.Request.URL.Query()
记得要先判断query["page"]切片是否有值
pageStr := query["page"][0]
因为是从url取参数 所以不需要定义struct了
- 使用 nav-inc/datetime 以支持 ISO8601 格式的时间解析
datetime.Parse(happenedBeforeStr, time.Local)
- 使用 sqlc.arg(happened_after) 为生成的函数参数命名，否则自动命名很难用
SELECT * from items
WHERE happened_at >= sqlc.arg(happened_after) AND happened_at < sqlc.arg(happened_before)
ORDER BY happened_at DESC ;
- 构造请求时，ISO8601 中的 + 会被转译，因此要使用url.QueryEscape包一下
req, _ := http.NewRequest(
		"GET",
		"/api/v1/items/balance?happened_after="+url.QueryEscape("2020-01-01T00:00:00+0800")+
			"&happened_before="+url.QueryEscape("2020-01-02T00:00:00+0800"),
		nil,
	)

## 开发 tags api
0. 先过一遍测试用例
go test ./...
1. 创建 controller 文件与 test 文件,router 注册路由
ctrl shift p go stubs
输入: ctrl *TagController account/internal/controller.Controller
2. 创建数据库迁移文件 
go build .; ./account db migrate:create create_tags_table
写上升降级的代码，同步数据库
go build .; ./account db migrate:up
3. 复制建表语句到 config/schema.sql以供sqlc生成model代码
sqlc generate
4. 声明 api 出入参类型
5. 写创建 tag 的 sql 语句
6. 写 controller
7. 测试
8. 生成文档
报错：Error parsing type definition 'queries.Tag': cannot find type definition: sql.NullTime
因为 deleted_at 字段在数据库中是一个可以为null的时间戳
而 go 语言的 time.Time 是一个结构体，永远不会为null
为了表示null，database/sql 库提供了 NullTime类型
type NullTime struct {
    Time  time.Time
    Valid bool // Valid is true if Time is not NULL
}
在使用 swag init 命令时要加上 --parseDependency，能够解析依赖的类型
但这与我们的要求不一致
搜索 sqlc overwrite null type
sqlc.yaml加上如下配置
overrides:
        - db_type: "pg_catalog.timestamp"
          go_type:
            import: "time"
            type: "Time"
            pointer: true
          nullable: true
使得 deleted_at 字段在go中的类型为 *time.Time 指针，指针是可以为 空值的

controller 主要负责控制下类型 处理报错 取接口数据 全扔给sqlc

目前测试未通过时，控制台输出的信息过多，不方便debug
造成该现象的原因：gin的logger中间件 和 debug模式
把gin框架的logger中间件删掉
进入gin.Default的源码，会发现Default会使用两个中间件
engine.Use(Logger(), Recovery())
按照源码的流程写一遍，不使用logger即可
将debug模式改为release模式（log中有相关提示）

小问题
在items controller 中，我们需要在查询参数中取happened_after属性
之前是在 query 对象中取的，并对数组越界异常进行了处理
query := c.Request.URL.Query()
	if len(query["happened_after"]) > 0 {
		happenedAfterStr := query["happened_after"][0]
		pt, err := datetime.Parse(happenedAfterStr, time.Local)
		if err == nil {
			happenedAfter = pt
		}
	}
其实用query.Get()就可以解决

删除 kind 类型
校验使用go来做，就不需要postgres的kind类型来限制了
写一个 migration 来删除数据库的kind类型
使用事务进行表配置的更改即类型的删除
修改 schema.sql 重新生成 sqlc

go更新数据库常见问题
如果在更新数据时，某字段传了空字符串""，那么要不要覆盖原有数据
声明接口类型时，指定了某字段（例如）类型为 string，且不能为空
所以即使不传这个字段，go 会将该字段赋值为 ""
这时会导致调用更新接口时，不仅要传更新的字段，也要将剩余字段传回给 go，否则 go 会使用 "" 覆盖掉初始值
- 方案一 在sql语句中使用CASE WHEN THEN忽略空字符串的情况
UPDATE tags SET sign = CASE WHEN @sign = '' THEN sign ELSE @sign END WHERE id = @id RETURNING *;
方案一可以解决必填字符串类型字段的问题
```
- name: UpdateTag :one
UPDATE tags  
SET 
  user_id = @user_id,
  name = CASE WHEN @name::varchar = '' THEN name ELSE @name END,
  sign = CASE WHEN @sign::varchar = '' THEN sign ELSE @sign END,
  kind = CASE WHEN @kind::varchar = '' THEN kind ELSE @kind END
WHERE id = @id
RETURNING id, user_id, name, sign, kind, deleted_at, created_at, updated_at;
```

声明接口类型时，指定了某字段类型为 string，同时可以为空
此时传 nil 表示不更新该字段，传 "" 表示将字段更新为空字符串
- 方案二 使用 NullString 类型
当声明数据库的 schema 时，某字段类型为 VARCHAR(100) 没有 NOT NULL 关键字
在 sqlc 生成的 go 代码中，该字段的类型会被指定为 sql.NullString 而不是 string
type NullString struct {
  String string
  Valid bool // Valid is true if String is not NULL
}
NullString 类型是一个结构体，导致该字段会在响应体中变成一个对象返回给前端，而且在请求体中也要将该字段作为对象传回来
完全不符合我们的开发使用习惯

- 方案二改进 重写 NullString 类型
重写 NullString 类型主要是重写结构体的序列化和反序列化方法：MarshalJSON UnmarshalJSON；数据库读写方法：Scan 和 Value
然后修改 sqlc 配置，指定 varchar 类型非空时的类型为我们重写的 MyNullString 即可
其实也可以让 MyNullString 继承 sql.NullString

- 方案三 使用 *string
可以在 sqlc 的配置中将 varchar 类型非空的字段指定为字符串的指针类型（pointer: true）
如果不传值（或传 null），那么 go 会默认变成 nil（null）
如果传空字符串，那么 go 会默认变成空字符串
麻烦的点在于，如果要在go中使用这个字段的值，需要提前进行判空处理（因为指针为空时，取值会报错）

- 方案四 使用 null 库
go get gopkg.in/guregu/null.v4
```
- db_type: "pg_catalog.varchar"
          nullable: true
          go_type:
            import: "gopkg.in/guregu/null.v4"
            type: "String"
            pointer: false
```
修改api出入参类型为 null.String
注意使用时要判断结构体的数据而不是直接作为字符串判断
assert.Equal(t, "xxx", j.Resource.X.String)