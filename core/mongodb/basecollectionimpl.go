package mongodb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// BaseCollectionImpl
// @Description: collection 的实现
type BaseCollectionImpl struct {
	DbName     string
	ColName    string
	DataBase   *mongo.Database
	Collection *mongo.Collection
}

// 分页查询
func (b *BaseCollectionImpl) SelectPage(ctx context.Context, filter interface{}, sort interface{}, skip, limit int64) (int64, []interface{}, error) {
	var err error

	resultCount, err := b.Collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, nil, err
	}

	opts := options.Find().SetSort(sort).SetSkip(skip).SetLimit(limit)
	finder, err := b.Collection.Find(ctx, filter, opts)
	if err != nil {
		return resultCount, nil, err
	}

	result := make([]interface{}, 0)
	if err := finder.All(ctx, &result); err != nil {
		return resultCount, nil, err
	}
	return resultCount, result, nil
}

// 查询列表
func (b *BaseCollectionImpl) SelectList(ctx context.Context, filter interface{}, sort interface{}) ([]interface{}, error) {
	var err error

	opts := options.Find().SetSort(sort)
	finder, err := b.Collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	result := make([]interface{}, 0)
	if err := finder.All(ctx, &result); err != nil {
		return nil, err
	}
	return result, err
}

// 查询单条
func (b *BaseCollectionImpl) SelectOne(ctx context.Context, filter interface{}) (interface{}, error) {
	result := new(interface{})
	err := b.Collection.FindOne(ctx, filter, options.FindOne()).Decode(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// 查询统计
func (b *BaseCollectionImpl) SelectCount(ctx context.Context, filter interface{}) (int64, error) {
	return b.Collection.CountDocuments(ctx, filter)
}

// 更新单条
func (b *BaseCollectionImpl) UpdateOne(ctx context.Context, filter, update interface{}) (int64, error) {
	result, err := b.Collection.UpdateOne(ctx, filter, update, options.Update())
	if err != nil {
		return 0, err
	}
	if result.MatchedCount == 0 {
		return 0, fmt.Errorf("Update result: %s ", "document not found")
	}
	return result.MatchedCount, nil
}

// 更新多条
func (b *BaseCollectionImpl) UpdateMany(ctx context.Context, filter, update interface{}) (int64, error) {
	result, err := b.Collection.UpdateMany(ctx, filter, update, options.Update())
	if err != nil {
		return 0, err
	}
	if result.MatchedCount == 0 {
		return 0, fmt.Errorf("Update result: %s ", "document not found")
	}
	return result.MatchedCount, nil
}

// 根据条件删除
func (b *BaseCollectionImpl) Delete(ctx context.Context, filter interface{}) (int64, error) {
	result, err := b.Collection.DeleteMany(ctx, filter, options.Delete())
	if err != nil {
		return 0, err
	}
	if result.DeletedCount == 0 {
		return 0, fmt.Errorf("DeleteOne result: %s ", "document not found")
	}
	return result.DeletedCount, nil
}

// 插入单条
func (b *BaseCollectionImpl) InsetOne(ctx context.Context, model interface{}) (interface{}, error) {
	result, err := b.Collection.InsertOne(ctx, model, options.InsertOne())
	if err != nil {
		return nil, err
	}
	return result.InsertedID, err
}

// 插入多条
func (b *BaseCollectionImpl) InsertMany(ctx context.Context, models []interface{}) ([]interface{}, error) {
	result, err := b.Collection.InsertMany(ctx, models, options.InsertMany())
	if err != nil {
		return nil, err
	}
	return result.InsertedIDs, err
}

// 聚合查询
func (b *BaseCollectionImpl) Aggregate(ctx context.Context, pipeline interface{}, result interface{}) error {
	finder, err := b.Collection.Aggregate(ctx, pipeline, options.Aggregate())
	if err != nil {
		return err
	}
	if err := finder.All(ctx, &result); err != nil {
		return err
	}
	return nil
}

// 创建索引，用于初始化时调用
func (b *BaseCollectionImpl) CreateIndexes(ctx context.Context, indexes []mongo.IndexModel) error {
	_, err := b.Collection.Indexes().CreateMany(ctx, indexes, options.CreateIndexes())
	return err
}

// 获取当前的*mongo.Collection对象
func (b *BaseCollectionImpl) GetCollection() *mongo.Collection {
	return b.Collection
}
