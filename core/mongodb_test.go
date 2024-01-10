package core

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/baowk/dilu-core/core/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
	"time"
)

var (
	userCollection *mongo.Collection
	// ……多个
)

type UserDto struct {
	Name string `json:"name" bson:"name"`
	Age  int64  `json:"age" bson:"age"`
}

func (a *UserDto) BsonByte(item interface{}) {
	bytes, err := bson.Marshal(item)
	if err != nil {
		panic(err)
	}
	bson.Unmarshal(bytes, a)
}

// @Description: 初始化mongo
// @param mongoUrl
func TestMongo(t *testing.T) {

	// 初始化链接
	timeout := time.Duration(600) * time.Second
	url := "mongodb://127.0.0.1:27017"
	cleanFunc, err := mongodb.New(url, timeout)
	if err != nil {
		cleanFunc()
		fmt.Println("======================")
		fmt.Println(err)
		fmt.Println("======================")
	}

	ctx := context.TODO()

	//  初始化user的集合对象
	userCollection := mongodb.CreateMongoCollection("test", "user")

	// 初始化索引
	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{"name", 1}, {"age", -1}},
			Options: options.Index().SetUnique(true),
		},
	}
	if err := userCollection.CreateIndexes(ctx, indexes); err != nil {
		fmt.Printf("init indexes error: %v", err.Error())
	}

	filter := bson.M{
		"name": bson.M{
			"$regex":   "ike",
			"$options": "i",
		},
	}
	list, err := userCollection.SelectList(ctx, filter, nil)
	if err != nil {
		panic(err)
	}

	userList := make([]*UserDto, 0)
	for _, item := range list {
		app := new(UserDto)
		app.BsonByte(item)
		userList = append(userList, app)
	}

	fmt.Printf("userList: length:%v\n", len(userList))
	bytes, _ := json.Marshal(userList)
	fmt.Printf("userList: data:%v\n", string(bytes))
}
