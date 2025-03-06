// Code generated by https://github.com/go-dev-frame/sponge

package handler

import (
	"context"
	schoolV1 "school/api/school/v1"
	//"github.com/go-dev-frame/sponge/pkg/gin/middleware"
)

var _ schoolV1.SchoolLogicer = (*schoolHandler)(nil)

type schoolHandler struct {
	// example:
	// 	userDao dao.UserDao
}

// NewUserHandler create a handler
func NewSchoolHandler() schoolV1.SchoolLogicer {
	return &schoolHandler{
		// example:
		// 	userDao: dao.NewUserDao(
		// 		database.GetDB(),
		// 		cache.NewUserCache(database.GetCacheType()),
		// 	),
	}
}

// Hello ......
func (h *schoolHandler) Hello(ctx context.Context, req *schoolV1.HelloRequest) (*schoolV1.HelloReply, error) {
	panic("implement me")

	// fill in the business logic code here
	// example:
	//
	//	    err := req.Validate()
	//	    if err != nil {
	//		    logger.Warn("req.Validate error", logger.Err(err), logger.Any("req", req), middleware.CtxRequestIDField(ctx))
	//		    return nil, ecode.InvalidParams.Err()
	//	    }
	//
	//	    reply, err := h.userDao.Hello(ctx, &model.User{
	//     	Name: req.Name,
	//     })
	//	    if err != nil {
	//			logger.Warn("Hello error", logger.Err(err), middleware.CtxRequestIDField(ctx))
	//			return nil, ecode.InternalServerError.Err()
	//		}
	//
	//     return &schoolV1.HelloReply{
	//     	Msg: reply.Msg,
	//     }, nil
}
