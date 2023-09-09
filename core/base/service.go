package base

import (
	"github.com/baowk/dilu-core/core"
	"github.com/gin-gonic/gin"
)

type BaseService struct {
	DbName string
	C      *gin.Context
}

func (s *BaseService) MakeContext(c *gin.Context) {
	s.C = c
}

func (s *BaseService) Create(data any) error {
	return core.Db(s.DbName).Create(data).Error
}

func (s *BaseService) Save(data any) error {
	return core.Db(s.DbName).Save(data).Error
}

func (s *BaseService) Get(id any, model any) error {
	return core.Db(s.DbName).First(model, id).Error
}

func (s *BaseService) DelWhere(where any) error {
	return core.Db(s.DbName).Delete(where).Error
}

func (s *BaseService) DelIds(model any, ids any) error {
	return core.Db(s.DbName).Delete(model, ids).Error
}

func (s *BaseService) Page(where any, data any, total *int64, limit, offset int) error {
	return core.Db(s.DbName).Where(where).Limit(limit).Offset(offset).
		Find(data).Limit(-1).Offset(-1).Count(total).Error
}

func (s *BaseService) UpdateWhere(where any, updates map[string]any) error {
	return core.Db(s.DbName).Where(where).Updates(updates).Error
}

func (s *BaseService) UpdateWhereModel(where any, updates any) error {
	return core.Db(s.DbName).Where(where).Updates(updates).Error
}

func (s *BaseService) GetByWhere(where any, model any) error {
	return core.Db(s.DbName).Where(where).Find(model).Error
}
