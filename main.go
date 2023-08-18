package main

import "account/cmd"

// 包的最后一个单词作为导出的变量
// 也可以用别名
// c "account/cmd"

func main() {
	cmd.RunServer()
}
