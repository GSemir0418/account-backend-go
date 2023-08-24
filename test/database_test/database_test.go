package database_test

import (
	"account/internal/database"
	"testing"
)

// 以 BenchmarkXXX 命名
func BenchmarkCrud(b *testing.B) {
	database.Connect()
	// database.CreateTables()
	database.Migrate()
	defer database.Close()
	// b.N 表示测试次数
	for i := 0; i < b.N; i++ {
		database.Crud()
	}
}
