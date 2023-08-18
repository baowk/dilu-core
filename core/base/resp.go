package base

import (
	"net/http"

	"github.com/baowk/dilu-core/common/consts"
	"github.com/baowk/dilu-core/core/errs"
	"github.com/gin-gonic/gin"
)

const (
	OK      = 200
	FAILURE = 500
)

type Resp struct {
	ReqId string `json:"reqId"`          //`json:"请求id"`
	Code  int    `json:"code"`           //返回码
	Msg   string `json:"msg,omitempty"`  //消息
	Data  any    `json:"data,omitempty"` //数据
}

type PageResp struct {
	List  any   `json:"list"`  //数据列表
	Total int64 `json:"total"` //总条数
	Size  int   `json:"size"`  //分页大小
	Page  int   `json:"page"`  //当前第几页
}

//type RespFunc func()

func Ok(c *gin.Context, data any) {
	c.AbortWithStatusJSON(http.StatusOK, Resp{
		ReqId: c.GetString(consts.REQ_ID),
		Code:  OK,
		Msg:   "ok",
		Data:  data,
	})
}

func Err(c *gin.Context, err errs.IError, msg string) {
	Fail(c, err.Code(), msg)
}

func Fail(c *gin.Context, code int, msg string, data ...any) {
	c.AbortWithStatusJSON(http.StatusOK, Resp{
		ReqId: c.GetString(consts.REQ_ID),
		Code:  code,
		Msg:   msg,
		Data:  data,
	})
}

func Page(c *gin.Context, list any, total int64, page int, pageSize int) {
	p := PageResp{
		Page:  page,
		Total: total,
		Size:  pageSize,
		List:  list,
	}
	Ok(c, p)
}
