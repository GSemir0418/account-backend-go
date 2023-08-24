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
	createMgrCmd := &cobra.Command{
		Use: "create:migration",
		Run: func(cmd *cobra.Command, args []string) {
			database.CreateMgr(args[0])
		},
	}
	mgrtCmd := &cobra.Command{
		Use: "migrate",
		Run: func(cmd *cobra.Command, args []string) {
			database.Migrate()
		},
	}
	crudCmd := &cobra.Command{
		Use: "crud",
		Run: func(cmd *cobra.Command, args []string) {
			database.Crud()
		},
	}

	database.Connect()
	rootCmd.AddCommand(srvCmd, dbCmd)
	dbCmd.AddCommand(createMgrCmd, mgrtCmd, crudCmd)
	// 会在当前函数执行结束后执行
	defer database.Close()
	rootCmd.Execute()
}

func RunServer() {
	r := router.New()
	r.Run(":8080")
}
