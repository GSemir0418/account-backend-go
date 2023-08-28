// cmd 包用来承担项目开发过程中的所有命令行任务，例如启动服务器、同步数据库等
package cmd

import (
	"account/internal/database"
	"account/internal/email"
	"account/internal/router"

	"github.com/spf13/cobra"
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

	database.Connect()
	// 会在当前函数执行结束后执行
	defer database.Close()

	rootCmd.AddCommand(srvCmd, dbCmd, emailCmd)
	dbCmd.AddCommand(mgrCreateCmd, mgrtDownCmd, mgrtUpCmd, crudCmd)

	rootCmd.Execute()
}

func RunServer() {
	r := router.New()
	r.Run(":8080")
}
