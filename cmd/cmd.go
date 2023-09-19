// cmd 包用来承担项目开发过程中的所有命令行任务，例如启动服务器、同步数据库等
package cmd

import (
	"account/internal/database"
	"account/internal/email"
	"account/internal/jwt_helper"
	"account/internal/router"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Run() {
	rootCmd := &cobra.Command{
		Use: "account",
	}
	srvCmd := &cobra.Command{
		Use: "server",
		Run: func(cmd *cobra.Command, args []string) {
			RunServer()
		},
	}
	generateHmacSecretCmd := &cobra.Command{
		Use: "generateHmacSecret",
		Run: func(cmd *cobra.Command, args []string) {
			// 生成jwt密钥并保存到本地
			bytes, _ := jwt_helper.GenerateHmacSecret()
			keyPath := viper.GetString("jwt.hmac.keyPath")
			if err := os.WriteFile(keyPath, bytes, 0644); err != nil {
				log.Fatalln(err)
			}
			fmt.Println("HMAC key has been saved in ", keyPath)
		},
	}
	dbCmd := &cobra.Command{
		Use: "db",
	}
	emailCmd := &cobra.Command{
		Use: "email",
		Run: func(cmd *cobra.Command, args []string) {
			email.Send()
		},
	}
	mgrCreateCmd := &cobra.Command{
		Use: "migrate:create",
		Run: func(cmd *cobra.Command, args []string) {
			database.MigrateCreate(args[0])
		},
	}
	mgrtUpCmd := &cobra.Command{
		Use: "migrate:up",
		Run: func(cmd *cobra.Command, args []string) {
			database.MigrateUp()
		},
	}
	mgrtDownCmd := &cobra.Command{
		Use: "migrate:down",
		Run: func(cmd *cobra.Command, args []string) {
			database.MigrateDown()
		},
	}
	crudCmd := &cobra.Command{
		Use: "crud",
		Run: func(cmd *cobra.Command, args []string) {
			database.Crud()
		},
	}
	coverCmd := &cobra.Command{
		Use: "coverage",
		Run: func(cmd *cobra.Command, args []string) {
			// 预先在根目录下创建 coverage 文件夹
			os.MkdirAll("coverage", os.ModePerm)
			// 使用 exec 执行命令行
			if err := exec.Command("MailHog").Start(); err != nil {
				log.Println(err)
			}
			if err := exec.Command(
				"go", "test", "-coverprofile=coverage/cover.out", "./...",
			).Run(); err != nil {
				log.Fatalln(err)
			}
			if err := exec.Command(
				"go", "tool", "cover", "-html=coverage/cover.out", "-o", "coverage/index.html",
			).Run(); err != nil {
				log.Fatalln(err)
			}
			// 使用 gin 开启本地服务
			var port string
			if len(args) > 0 {
				port = args[0]
			} else {
				port = "8888"
			}
			fmt.Println("http://localhost:" + port + "/coverage/index.html")
			if err := http.ListenAndServe(":"+port, http.FileServer(http.Dir("."))); err != nil {
				log.Fatalln(err)
			}
		},
	}
	database.Connect()
	// 会在当前函数执行结束后执行
	defer database.Close()

	rootCmd.AddCommand(srvCmd, dbCmd, emailCmd, generateHmacSecretCmd, coverCmd)
	dbCmd.AddCommand(mgrCreateCmd, mgrtDownCmd, mgrtUpCmd, crudCmd)

	rootCmd.Execute()
}

func RunServer() {
	r := router.New()
	r.Run(":8080")
}
