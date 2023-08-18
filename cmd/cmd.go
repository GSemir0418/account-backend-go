// cmd 包用来承担项目开发过程中的所有命令行任务，例如启动服务器、同步数据库等
package cmd

import "account/internal/router"

// go中首字母大写会被默认导出
// go中变量名一般用一个字母
func RunServer() {
	r := router.New()
	// Listen and Server in 0.0.0.0:8080
	// If 127.0.0.1:8080, server will not be eble to serve request from outside web
	r.Run(":8080")
}
