package main

import "account/cmd"

// @title           记账
// @description     记账应用接口文档

// @contact.name   GSemir
// @contact.url    http://gsemir0418.github.com/
// @contact.email  gsemir0418@gmail.com

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth(JWT)

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/

func main() {
	// 初始化命令行程序
	cmd.Run()
}
