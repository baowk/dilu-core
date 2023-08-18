package inter

import (
	"github.com/baowk/dilu-core/common/middleware"
	"github.com/baowk/dilu-core/config"
	"github.com/gin-gonic/gin"
)

func Init(r *gin.Engine, cfg *config.AppCfg) {
	middleware.InitMiddleware(r, cfg)
}
