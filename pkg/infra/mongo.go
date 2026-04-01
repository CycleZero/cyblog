package infra

// 该项目暂时弃用MongoDB

//func NewMongoDb(vc *viper.Viper) *mongo.Client {
//	host := vc.GetString("data.mongo.host")
//	port := vc.GetString("data.mongo.port")
//	user := vc.GetString("data.mongo.username")
//	password := vc.GetString("data.mongo.password")
//	dbName := vc.GetString("data.mongo.db_name")
//
//	log.GetLogger().Sugar().Infof("MongoDB 配置：host=%s, port=%s, user=%s, dbName=%s", host, port, user, dbName)
//
//	// 使用 ClientOptions 方式配置连接
//	clientOpts := options.Client().
//		ApplyURI("mongodb://" + host + ":" + port).
//		SetAuth(options.Credential{
//			Username: user,
//			Password: password,
//		}).
//		SetDirect(true). // 直接连接，不使用自动发现
//		SetServerSelectionTimeout(30 * time.Second)
//
//	log.GetLogger().Sugar().Infof("MongoDB 连接地址：%s:%s", host, port)
//
//	// 设置超时时间
//	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
//	defer cancel()
//
//	// 创建客户端并连接
//	client, err := mongo.Connect(clientOpts)
//	if err != nil {
//		log.GetLogger().Sugar().Errorf("MongoDB Connect 失败：%v", err)
//		panic(err)
//	}
//
//	// 验证连接（Ping 测试）
//	err = client.Ping(ctx, nil)
//	if err != nil {
//		log.GetLogger().Sugar().Errorf("MongoDB Ping 失败：%v", err)
//		panic(err)
//	}
//
//	log.GetLogger().Sugar().Infof("成功连接到 MongoDB")
//
//	return client
//}
