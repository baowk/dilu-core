package base

import (
	"github.com/baowk/dilu-core/common/consts"
	"github.com/baowk/dilu-core/common/utils"
	"github.com/baowk/dilu-core/core/errs"
	"github.com/gin-gonic/gin"
)

type BaseApi struct {
}

func (e *BaseApi) GetReqId(c *gin.Context) string {
	return utils.GetReqId(c)
}

func (e *BaseApi) GetUserId(c *gin.Context) int {
	return c.GetInt(consts.USER_ID)
}

func (e *BaseApi) GetTenantId(c *gin.Context) int {
	return c.GetInt(consts.TENANT_ID)
}

func (e *BaseApi) Error(c *gin.Context, err error) {
	resMsg(c, FAILURE, err.Error())

}

func (e *BaseApi) Fail(c *gin.Context, code int, msg string, data ...any) {
	resMsg(c, code, msg, data...)
}

func (e *BaseApi) Err(c *gin.Context, err errs.IError) {
	errer(c, err)
}

func (e *BaseApi) Ok(c *gin.Context, data ...any) {
	ok(c, data...)
}

func (e *BaseApi) ResCustom(c *gin.Context, opts ...Option) {
	result(c, opts...)
}

func (e *BaseApi) Page(c *gin.Context, list any, total int64, page, size int) {
	pageResp(c, list, total, page, size)
}

//封装后代码路径指定到这里所以去掉
// func (e *BaseApi) LogError(c *gin.Context, err error) {
// 	core.Log.Error(fmt.Sprintf("REQID:%s", e.GetReqId(c)), zap.Error(err))
// }

// func (e *BaseApi) LogInfo(c *gin.Context, key string, val any) {
// 	ccore.Log.Info("REQID"+e.GetReqId(c), zap.Reflect("data", data))
// }
