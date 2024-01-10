package mongodb

import (
	"context"
	"github.com/baowk/dilu-core/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

var (
	client *mongo.Client
	ctx    = context.Background()
	cancel context.CancelFunc
	err    error
)

const (
	defaultTimeout = 50 * time.Second
	maxPoolSize    = 10
)

// MongoInit
// @Description: 初始化mongo
// @param mongoUrl
func MongoInit(conf config.Mongodb) {
	// 初始化链接
	timeout := time.Duration(conf.Timeout) * time.Second
	cleanFunc, err := New(conf.URL, timeout)
	if err != nil {
		cleanFunc()
		panic("mongo connect ping failed, err:" + err.Error())
	}
}

// New 创建mongo客户端实例
func New(url string, timeout time.Duration) (func(), error) {

	if t := timeout; t > 0 {
		ctx, cancel = context.WithTimeout(ctx, t)
		defer cancel()
	} else {
		ctx, cancel = context.WithTimeout(ctx, defaultTimeout)
	}

	// 链接
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(url).SetMaxPoolSize(maxPoolSize))
	if err != nil {
		return nil, err
	}

	// 断开
	cleanFunc := func() {
		err = client.Disconnect(context.Background())
		if err != nil {
			log.Printf("Mongo disconnect error: %s", err.Error())
		}
	}

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	return cleanFunc, nil
}

// CreateMongoCollection
// @Description: 创建mongo集合的服务
// @param dbName
// @param colName
// @return BaseCollection
// @return error
func CreateMongoCollection(dbName, colName string) BaseCollection {
	dataBase := client.Database(dbName)
	return &BaseCollectionImpl{
		DbName:     dbName,
		ColName:    colName,
		DataBase:   dataBase,
		Collection: dataBase.Collection(colName),
	}
}
