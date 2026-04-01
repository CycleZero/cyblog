package infra

import (
	"gorm.io/gorm"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewData,
	NewRedisClient,
	NewCustomRedisClient,
	NewNatsMQ,
	NewMinioClient,
	NewMongoDb,
	NewLLmClient,
	ProvideDB,
)

// ProvideDB 提供 *gorm.DB 实例
func ProvideDB(data *Data) *gorm.DB {
	return data.DB
}
