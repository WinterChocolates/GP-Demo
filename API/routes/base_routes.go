package routes

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"API/middlewares"
	"API/storage/cache"
	"API/storage/database"
	"API/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type BaseRouter struct {
	*gin.Engine
	controllers Controllers
}

func NewBaseRouter(engine *gin.Engine, ctrls Controllers) *BaseRouter {
	return &BaseRouter{
		Engine:      engine,
		controllers: ctrls,
	}
}

func (r *BaseRouter) ApplyGlobalMiddlewares() {
	r.Use(
		middlewares.RequestContext(),
		middlewares.EnhancedLogger(zap.L()),
		middlewares.ErrorLogger(zap.L()),
		middlewares.EnhancedSecurityHeaders(),
		middlewares.EnhancedCORSMiddleware(),
		middlewares.EnhancedRateLimiter(zap.L()),
		middlewares.PerformanceMonitor(zap.L()),
		gin.Recovery(),
	)
}

func (r *BaseRouter) SetupErrorHandlers() {
	r.NoRoute(func(c *gin.Context) {
		utils.RespondWithJSON(c, http.StatusNotFound, 404, "请求的资源不存在", nil)
	})

	r.NoMethod(func(c *gin.Context) {
		utils.RespondWithJSON(c, http.StatusMethodNotAllowed, 405, "不支持的请求方法", nil)
	})
}

func (r *BaseRouter) enhancedHealthCheckHandler(c *gin.Context) {
	type healthCheck func(context.Context) error

	checks := map[string]healthCheck{
		"database": database.CheckMySQLHealth,
		"redis":    cache.CheckRedisHealth,
	}

	results := make(map[string]string)
	overallStatus := http.StatusOK

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	for name, check := range checks {
		if err := check(ctx); err != nil {
			results[name] = fmt.Sprintf("异常: %v", err)
			overallStatus = http.StatusServiceUnavailable
		} else {
			results[name] = "正常"
		}
	}

	c.JSON(overallStatus, gin.H{
		"status":  utils.MapStatusToString(overallStatus),
		"checks":  results,
		"version": "1.0.0",
		"time":    time.Now().Format(time.RFC3339),
	})
}
