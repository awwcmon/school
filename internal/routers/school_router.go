// Code generated by https://github.com/go-dev-frame/sponge

package routers

import (
	"github.com/gin-gonic/gin"

	"github.com/go-dev-frame/sponge/pkg/logger"
	//"github.com/go-dev-frame/sponge/pkg/middleware"

	schoolV1 "school/api/school/v1"
	"school/internal/handler"
)

func init() {
	allMiddlewareFns = append(allMiddlewareFns, func(c *middlewareConfig) {
		schoolMiddlewares(c)
	})

	allRouteFns = append(allRouteFns,
		func(r *gin.Engine, groupPathMiddlewares map[string][]gin.HandlerFunc, singlePathMiddlewares map[string][]gin.HandlerFunc) {
			schoolRouter(r, groupPathMiddlewares, singlePathMiddlewares, handler.NewSchoolHandler())
		})
}

func schoolRouter(
	r *gin.Engine,
	groupPathMiddlewares map[string][]gin.HandlerFunc,
	singlePathMiddlewares map[string][]gin.HandlerFunc,
	iService schoolV1.SchoolLogicer) {
	schoolV1.RegisterSchoolRouter(
		r,
		groupPathMiddlewares,
		singlePathMiddlewares,
		iService,
		schoolV1.WithSchoolLogger(logger.Get()),
		schoolV1.WithSchoolHTTPResponse(),
		schoolV1.WithSchoolErrorToHTTPCode(
		// Set some error codes to standard http return codes,
		// by default there is already ecode.InternalServerError and ecode.ServiceUnavailable
		// example:
		// 	ecode.Forbidden, ecode.LimitExceed,
		),
	)
}

// you can set the middleware of a route group, or set the middleware of a single route,
// or you can mix them, pay attention to the duplication of middleware when mixing them,
// it is recommended to set the middleware of a single route in preference
func schoolMiddlewares(c *middlewareConfig) {
	// set up group route middleware, group path is left prefix rules,
	// if the left prefix is hit, the middleware will take effect, e.g. group route is /api/v1, route /api/v1/school/:id  will take effect
	// c.setGroupPath("/api/v1/school", middleware.Auth())

	// set up single route middleware, just uncomment the code and fill in the middlewares, nothing else needs to be changed
	//c.setSinglePath("POST", "/api/v1/login", middleware.Auth())    // Login 登录注释
	//c.setSinglePath("GET", "/api/v1/hello", middleware.Auth())    // Hello ......
}
