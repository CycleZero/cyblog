package infra

import (
	"context"
	"cyblog/pkg/log"
	"encoding/json"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var globalData *Data

type Data struct {
	DB          *gorm.DB
	logger      *zap.Logger
	RedisClient *RedisClient
	MongoDb     *mongo.Client
	MinioClient *minio.Client
	NatsMQ      *NatsMQ
}

type RedisClient struct {
	*redis.Client
}

// target 为指针类型
func (r *RedisClient) GetObject(ctx context.Context, key string, target any) error {

	res, err := r.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(res), target)
}

func (r *RedisClient) PutObject(ctx context.Context, key string, target any, expiration time.Duration) error {

	str, err := json.Marshal(target)
	if err != nil {
		return err
	}
	return r.SetEx(ctx, key, string(str), expiration).Err()
}

func NewData(
	vc *viper.Viper,
	rdb *RedisClient,
	mongoc *mongo.Client,
	minioClient *minio.Client,
	natsMQ *NatsMQ,
) *Data {

	logger, _ := zap.NewDevelopment()
	var masterDB *gorm.DB
	host := vc.GetString("data.db.host")
	port := vc.GetString("data.db.port")
	user := vc.GetString("data.db.user")
	password := vc.GetString("data.db.password")
	dbname := vc.GetString("data.db.db_name")
	dsn := GetDsn(host, port, user, password, dbname)
	var err error
	masterDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		logger.Info("连接数据库失败")
		panic(err)
	}
	initTaskQueue()
	return &Data{
		DB:          masterDB,
		logger:      logger,
		MongoDb:     mongoc,
		RedisClient: rdb,
		MinioClient: minioClient,
		NatsMQ:      natsMQ,
	}
}
func GetDsn(host, port, user, password, dbname string) string {
	//gorm:gorm@tcp(localhost:9910)/gorm?charset=utf8&parseTime=True&loc=Local
	dsn := user + ":" + password + "@tcp(" + host + ":" + port + ")/" + dbname + "?charset=utf8&parseTime=True&loc=Local"
	log.GetLogger().Info("生成DSN: " + dsn)
	return dsn
}

func NewRedisClient(vc *viper.Viper) *redis.Client {

	host := vc.GetString("data.redis.host")
	port := vc.GetString("data.redis.port")
	password := vc.GetString("data.redis.password")
	log.GetLogger().Info("Redis连接信息: " + host + ":" + port)
	rdb := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: password,
		DB:       8, // use default DB
	})
	log.GetLogger().Info("连接Redis成功")
	return rdb
}

func NewCustomRedisClient(rdb *redis.Client) *RedisClient {
	return &RedisClient{rdb}
}

func NewMinioClient(vc *viper.Viper) *minio.Client {
	endpoint := vc.GetString("data.minio.endpoint")
	accessKeyID := vc.GetString("data.minio.accessKey")
	secretAccessKey := vc.GetString("data.minio.secretKey")
	bucketName := vc.GetString("data.minio.bucket")

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: false,
	})
	if err != nil {
		log.GetLogger().Error("Error creating minio client", zap.Error(err))
		panic(err)
	}

	ctx := context.Background()
	err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			log.GetLogger().Sugar().Infof("We already own %s\n", bucketName)
		} else {
			log.GetLogger().Sugar().Fatal(err)
		}
	} else {
		log.GetLogger().Sugar().Infof("Successfully created %s\n", bucketName)
	}
	return minioClient
}

func GetData() *Data {
	if globalData == nil {
		panic("请先初始化Data")
	}
	return globalData
}

func GetDB() *gorm.DB {
	return GetData().DB
}

func (d *Data) Close() error {
	err := d.MongoDb.Disconnect(context.Background())
	if err != nil {
		return err
	}
	return nil
}

// WithTransaction 从Context中获取父事务（父DB），子事务需自行调用db.Begin().无论是否成功，都会绑定context
func (d *Data) WithTransaction(c context.Context) *gorm.DB {
	if c == nil {
		return d.DB.WithContext(c)
	}
	db := GetTransaction(c)
	if db != nil {
		return db
	}
	return d.DB.WithContext(c)

}

func SetTransaction(c context.Context, tran *gorm.DB) context.Context {
	t := tran.WithContext(c)
	return context.WithValue(c, "transaction", t)
}

func GetTransaction(c context.Context) *gorm.DB {
	db, ok := c.Value("transaction").(*gorm.DB)
	if ok {
		return db
	}
	return nil
}
