package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

// BaseCollection
// @Description: 定义操作的接口
type BaseCollection interface {

	//  SelectPage
	//  @Description: 分页查询
	//  @param ctx
	//  @param filter
	//  @param sort
	//  @param skip
	//  @param limit
	//  @return int64
	//  @return []interface{}
	//  @return error
	SelectPage(ctx context.Context, filter interface{}, sort interface{}, skip, limit int64) (int64, []interface{}, error)

	//  SelectList
	//  @Description: 查询列表
	//  @param ctx
	//  @param filter
	//  @param sort
	//  @return []interface{}
	//  @return error
	SelectList(ctx context.Context, filter interface{}, sort interface{}) ([]interface{}, error)

	//  SelectOne
	//  @Description: 查询单条
	//  @param ctx
	//  @param filter
	//  @return interface{}
	//  @return error
	SelectOne(ctx context.Context, filter interface{}) (interface{}, error)

	//  SelectCount
	//  @Description: 查询统计
	//  @param ctx
	//  @param filter
	//  @return int64
	//  @return error
	SelectCount(ctx context.Context, filter interface{}) (int64, error)

	//  UpdateOne
	//  @Description: 更新单条
	//  @param ctx
	//  @param filter
	//  @param update
	//  @return int64
	//  @return error
	UpdateOne(ctx context.Context, filter, update interface{}) (int64, error)

	//  UpdateMany
	//  @Description: 更新多条
	//  @param ctx
	//  @param filter
	//  @param update
	//  @return int64
	//  @return error
	UpdateMany(ctx context.Context, filter, update interface{}) (int64, error)

	//  Delete
	//  @Description: 根据条件删除
	//  @param ctx
	//  @param filter
	//  @return int64
	//  @return error
	Delete(ctx context.Context, filter interface{}) (int64, error)

	//  InsetOne
	//  @Description: 插入单条
	//  @param ctx
	//  @param model
	//  @return interface{}
	//  @return error
	InsetOne(ctx context.Context, model interface{}) (interface{}, error)

	//  InsertMany
	//  @Description: 插入多条
	//  @param ctx
	//  @param models
	//  @return []interface{}
	//  @return error
	InsertMany(ctx context.Context, models []interface{}) ([]interface{}, error)

	//  Aggregate
	//  @Description: 聚合查询
	//  @param ctx
	//  @param pipeline
	//  @param result
	//  @return error
	Aggregate(ctx context.Context, pipeline interface{}, result interface{}) error

	//  CreateIndexes
	//  @Description: 创建索引，用于初始化时调用
	//  @param ctx
	//  @param indexes
	//  @return error
	CreateIndexes(ctx context.Context, indexes []mongo.IndexModel) error

	//  GetCollection
	//  @Description: 获取当前的*mongo.Collection对象
	//  @return *mongo.Collection
	GetCollection() *mongo.Collection
}
