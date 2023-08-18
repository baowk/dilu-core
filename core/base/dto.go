package base

type ReqBase struct {
	ReqId string `json:"reqId" form:"reqId"` // 请求id 链路跟踪
}

type ReqId struct {
	Id int `json:"id" form:"id"` // 主键ID
}

type ReqStrId struct {
	Id string `json:"id" form:"id"` // 主键ID
}

type ReqIds struct {
	Ids []int `json:"ids" form:"ids"` //多id
}

type ReqPage struct {
	Page int `json:"page" form:"page"` // 页码
	Size int `json:"size" form:"size"` // 每页大小
	//	Keyword string `json:"keyword" form:"keyword"` //关键字
}

func (e *ReqPage) GetPage() int {
	if e.Page < 1 {
		return 1
	}
	return e.Page
}

func (e *ReqPage) GetSize() int {
	if e.Size < 1 {
		return 10
	}
	return e.Size
}

func (e *ReqPage) GetOffset() int {
	return (e.GetPage() - 1) * e.GetSize()
}
