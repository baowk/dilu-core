package base

import (
	"net/http"

	"github.com/baowk/dilu-core/common/consts"
	"github.com/baowk/dilu-core/core/errs"
	"github.com/baowk/dilu-core/core/i18n"
	"github.com/gin-gonic/gin"
)

const (
	OK      = 200
	FAILURE = 500
)

type Resp struct {
	ReqId string `json:"reqId,omitempty"` //`json:"请求id"`
	Code  int    `json:"code"`            //返回码
	Msg   string `json:"msg,omitempty"`   //消息
	Data  any    `json:"data,omitempty"`  //数据
}

type PageResp struct {
	List        any   `json:"list"`        //数据列表
	Total       int64 `json:"total"`       //总条数
	PageSize    int   `json:"pageSize"`    //分页大小
	CurrentPage int   `json:"currentPage"` //当前第几页
}

type Option func(resp *Resp)

// func NewResp(opts ...Option) *Resp {
// 	r := new(Resp)
// 	for _, f := range opts {
// 		f(r)
// 	}
// 	return r
// }

func WithReqId(reqId string) Option {
	return func(resp *Resp) {
		resp.ReqId = reqId
	}
}

func WithCode(code int) Option {
	return func(resp *Resp) {
		resp.Code = code
	}
}

func WithMsg(msg string) Option {
	return func(resp *Resp) {
		resp.Msg = msg
	}
}

func WithData(data any) Option {
	return func(resp *Resp) {
		resp.Data = data
	}
}

func result(c *gin.Context, opts ...Option) {
	r := new(Resp)
	for _, f := range opts {
		f(r)
	}
	c.AbortWithStatusJSON(http.StatusOK, *r)
}

func pureJSON(c *gin.Context, data any) {
	c.PureJSON(http.StatusOK, Resp{
		ReqId: c.GetString(consts.REQ_ID),
		Code:  200,
		Msg:   "OK",
		Data:  data,
	})
}

func ok(c *gin.Context, data ...any) {
	resMsg(c, http.StatusOK, "OK", data...)
}

func errer(c *gin.Context, err errs.IError) {
	resMsg(c, err.Code(), "")
}

func resMsg(c *gin.Context, code int, msg string, data ...any) {
	if msg == "" {
		msg = i18n.Lang.GetMsg(code, c)
	}
	if len(data) == 0 {
		c.JSON(http.StatusOK, Resp{
			ReqId: c.GetString(consts.REQ_ID),
			Code:  code,
			Msg:   msg,
		})
	} else if len(data) == 1 {
		c.JSON(http.StatusOK, Resp{
			ReqId: c.GetString(consts.REQ_ID),
			Code:  code,
			Msg:   msg,
			Data:  data[0],
		})
	} else {
		c.JSON(http.StatusOK, Resp{
			ReqId: c.GetString(consts.REQ_ID),
			Code:  code,
			Msg:   msg,
			Data:  data,
		})
	}
}

func resMsgWithAbort(c *gin.Context, code int, msg string, data ...any) {
	if msg == "" {
		msg = i18n.Lang.GetMsg(code, c)
	}
	if len(data) == 0 {
		c.AbortWithStatusJSON(http.StatusOK, Resp{
			ReqId: c.GetString(consts.REQ_ID),
			Code:  code,
			Msg:   msg,
		})
	} else if len(data) == 1 {
		c.AbortWithStatusJSON(http.StatusOK, Resp{
			ReqId: c.GetString(consts.REQ_ID),
			Code:  code,
			Msg:   msg,
			Data:  data[0],
		})
	} else {
		c.AbortWithStatusJSON(http.StatusOK, Resp{
			ReqId: c.GetString(consts.REQ_ID),
			Code:  code,
			Msg:   msg,
			Data:  data,
		})
	}
}

func pageResp(c *gin.Context, list any, total int64, page int, pageSize int) {
	p := PageResp{
		CurrentPage: page,
		Total:       total,
		PageSize:    pageSize,
		List:        list,
	}
	resMsg(c, http.StatusOK, "OK", p)
}
