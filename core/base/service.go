package base

type BaseService struct {
}

// func (*BaseService) ParsePage(page, size int) (limit, offset int) {
// 	if page < 1 {
// 		page = 1
// 	}
// 	if size < 1 {
// 		size = 20
// 	}
// 	if size > 1000 {
// 		size = 1000
// 	}
// 	limit = size
// 	offset = (page - 1) * size
// 	return
// }

// func (s *BaseService) LogError(reqId string, err error) {
// 	core.Log.Error(fmt.Sprintf("REQID:%s", reqId), zap.Error(err))
// }

// func (e *BaseService) LogInfo(reqId string, key string, val any) {
// 	core.Log.Info(fmt.Sprintf("REQID:%s", reqId), zap.String(key, fmt.Sprintf("%v", val)))
// }
