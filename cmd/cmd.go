// cmd 包用来承担项目开发过程中的所有命令行任务，例如启动服务器、同步数据库等
package cmd

import (
	"account/internal/database"
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
	createCmd := &cobra.Command{
		Use: "create",
		Run: func(cmd *cobra.Command, args []string) {
			database.PgCreateTables()
		},
	}

	rootCmd.AddCommand(srvCmd)
	rootCmd.AddCommand(dbCmd)
	dbCmd.AddCommand(createCmd)
	database.PgConnect()
	// 会在当前函数执行结束后执行
	defer database.PgClose()
	rootCmd.Execute()
}

func RunServer() {
	r := router.New()
	r.Run(":8080")
}
