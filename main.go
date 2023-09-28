package main

import (
	"account/cmd"
	viper_config "account/config"
)

//	@title						记账
//	@description				记账应用接口文档
//
//	@contact.name				GSemir
//	@contact.url				http://gsemir0418.github.com/
//	@contact.email				gsemir0418@gmail.com
//
//	@host						localhost:8080
//	@BasePath					/
//
//	@securityDefinitions.apiKey	Bearer
//	@in							header
//	@name						Authorization
//
//	@externalDocs.description	OpenAPI
//	@externalDocs.url			https://swagger.io/resources/open-api/
func main() {
	// 读取本地密钥
	viper_config.LoadViperConfig()
	// 初始化命令行程序
	cmd.Run()
}
