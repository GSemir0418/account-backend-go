package database

import (
	"database/sql"
	"fmt"

	"log"

	// 引入pg的驱动，前面加_表示不直接访问这个包
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "gsemir"
	password = "gsemir"
	dbname   = "go_account_dev"
)

func PgConnect() {
	// 构建连接字符串
	connectStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", connectStr)
	if err != nil {
		log.Fatal(err)
	}
	// 及时在不同文件下，包名相同的变量可以访问
	DB = db
	// 两个步骤共享同一个 err
	err = DB.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Successfully connected to db")
}

func PgClose() {
	DB.Close()
	log.Println("Closed db")
}

func PgCreateTables() {
	// 创建 Users 表
	_, err := DB.Exec(`CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		email VARCHAR(100) NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Successfully created user table")
}

func Migrate() {
	// 示例：给 users 表添加手机字段
	_, err := DB.Exec(`ALTER TABLE users ADD COLUMN phone VARCHAR(50)`)
	if err != nil {
		// Fatal 会终止程序，影响后续命令
		// log.Fatal(err)
		log.Println(err)
	}
	log.Println("Successfully add phone column to users table")

	// 示例：新增 items 表
	// 复用 err
	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS items (
		id SERIAL PRIMARY KEY,
		amount VARCHAR(100) NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		log.Println(err)
	}
	log.Println("Successfully create items table")

	// 示例：修改 created_at 字段类型为 WITHOUT TIME ZONE
	_, err = DB.Exec(`ALTER TABLE items ALTER COLUMN created_at TYPE TIMESTAMP`)
	if err != nil {
		// Fatal 会终止程序，影响后续命令
		// log.Fatal(err)
		log.Println(err)
	}
	log.Println("Successfully change the type of created_at column")

	// 继续往下写同步的命令
	// 给 User 的 email 字段添加唯一索引
	_, err = DB.Exec(`CREATE UNIQUE INDEX users_email_index ON users (email)`)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("added unique index on users email")
	}
}

func Crud() {
	// 需要返回值用 Query，不需要就用 Exec
	// 创建一个 user
	_, err := DB.Exec(`INSERT INTO users (email) VALUES ('1@qq.com')`)
	// 错误处理：细分错误的类型
	if err != nil {
		// log.Fatalln(err)
		switch x := err.(type) {
		// 数据库错误
		case *pq.Error:
			pqError := err.(*pq.Error)
			log.Println(pqError.Code.Name())
			log.Println(pqError.Message)
		default:
			log.Println(x)
		}
	} else {
		log.Println("Successfully create a user")
	}

	// 修改一个 user
	_, err = DB.Exec(`Update users SET phone = 138123456789 where email = '1@qq.com'`)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Successfully update a user")
	}

	// 分页查询
	// 预准备语句 psql使用 $1 $2 作为参数占位符 mysql 使用?占位
	stmt, err := DB.Prepare("SELECT phone FROM users where email = $1 offset $2 limit $3")
	if err != nil {
		log.Fatalln(err)
	}
	result, err := stmt.Query("1@qq.com", 0, 3)
	if err != nil {
		log.Println(err)
	} else {
		// 返回值是一个迭代器而不是数组，节省内存
		for result.Next() {
			var phone string
			result.Scan(&phone)
			log.Println("phone", phone)
		}
		log.Println("Successfully read users")
	}

}
