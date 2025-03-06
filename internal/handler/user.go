package handler

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jinzhu/copier"

	"github.com/go-dev-frame/sponge/pkg/errcode"
	"github.com/go-dev-frame/sponge/pkg/gin/middleware"
	"github.com/go-dev-frame/sponge/pkg/logger"
	"github.com/go-dev-frame/sponge/pkg/mgo/query"

	schoolV1 "school/api/school/v1"
	"school/internal/cache"
	"school/internal/dao"
	"school/internal/database"
	"school/internal/ecode"
	"school/internal/model"
)

var _ schoolV1.UserLogicer = (*userHandler)(nil)
var _ time.Time

type userHandler struct {
	userDao dao.UserDao
}

// NewUserHandler create a handler
func NewUserHandler() schoolV1.UserLogicer {
	collectionName := new(model.User).TableName()
	return &userHandler{
		userDao: dao.NewUserDao(
			database.GetDB().Collection(collectionName),
			cache.NewUserCache(database.GetCacheType()),
		),
	}
}

// Create a record
func (h *userHandler) Create(ctx context.Context, req *schoolV1.CreateUserRequest) (*schoolV1.CreateUserReply, error) {
	err := req.Validate()
	if err != nil {
		logger.Warn("req.Validate error", logger.Err(err), logger.Any("req", req), middleware.CtxRequestIDField(ctx))
		return nil, ecode.InvalidParams.Err()
	}

	user := &model.User{}
	err = copier.Copy(user, req)
	if err != nil {
		return nil, ecode.ErrCreateUser.Err()
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	err = h.userDao.Create(ctx, user)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("user", user), middleware.CtxRequestIDField(ctx))
		return nil, ecode.InternalServerError.Err()
	}

	return &schoolV1.CreateUserReply{Id: user.ID.Hex()}, nil
}

// DeleteByID delete a record by id
func (h *userHandler) DeleteByID(ctx context.Context, req *schoolV1.DeleteUserByIDRequest) (*schoolV1.DeleteUserByIDReply, error) {
	err := req.Validate()
	if err != nil {
		logger.Warn("req.Validate error", logger.Err(err), logger.Any("req", req), middleware.CtxRequestIDField(ctx))
		return nil, ecode.InvalidParams.Err()
	}

	err = h.userDao.DeleteByID(ctx, req.Id)
	if err != nil {
		logger.Warn("DeleteByID error", logger.Err(err), middleware.CtxRequestIDField(ctx))
		return nil, ecode.InternalServerError.Err()
	}

	return &schoolV1.DeleteUserByIDReply{}, nil
}

// UpdateByID update a record by id
func (h *userHandler) UpdateByID(ctx context.Context, req *schoolV1.UpdateUserByIDRequest) (*schoolV1.UpdateUserByIDReply, error) {
	err := req.Validate()
	if err != nil {
		logger.Warn("req.Validate error", logger.Err(err), logger.Any("req", req), middleware.CtxRequestIDField(ctx))
		return nil, ecode.InvalidParams.Err()
	}

	user := &model.User{}
	err = copier.Copy(user, req)
	if err != nil {
		return nil, ecode.ErrUpdateByIDUser.Err()
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here
	user.ID = database.ToObjectID(req.Id)

	err = h.userDao.UpdateByID(ctx, user)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("user", user), middleware.CtxRequestIDField(ctx))
		return nil, ecode.InternalServerError.Err()
	}

	return &schoolV1.UpdateUserByIDReply{}, nil
}

// GetByID get a record by id
func (h *userHandler) GetByID(ctx context.Context, req *schoolV1.GetUserByIDRequest) (*schoolV1.GetUserByIDReply, error) {
	err := req.Validate()
	if err != nil {
		logger.Warn("req.Validate error", logger.Err(err), logger.Any("req", req), middleware.CtxRequestIDField(ctx))
		return nil, ecode.InvalidParams.Err()
	}

	record, err := h.userDao.GetByID(ctx, req.Id)
	if err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			logger.Warn("GetByID error", logger.Err(err), logger.Any("id", req.Id), middleware.CtxRequestIDField(ctx))
			return nil, ecode.NotFound.Err()
		}
		logger.Error("GetByID error", logger.Err(err), logger.Any("id", req.Id), middleware.CtxRequestIDField(ctx))
		return nil, ecode.InternalServerError.Err()
	}

	data, err := convertUser(record)
	if err != nil {
		logger.Warn("convertUser error", logger.Err(err), logger.Any("user", record), middleware.CtxRequestIDField(ctx))
		return nil, ecode.ErrGetByIDUser.Err()
	}

	return &schoolV1.GetUserByIDReply{
		User: data,
	}, nil
}

// List of records by query parameters
func (h *userHandler) List(ctx context.Context, req *schoolV1.ListUserRequest) (*schoolV1.ListUserReply, error) {
	err := req.Validate()
	if err != nil {
		logger.Warn("req.Validate error", logger.Err(err), logger.Any("req", req), middleware.CtxRequestIDField(ctx))
		return nil, ecode.InvalidParams.Err()
	}

	params := &query.Params{}
	err = copier.Copy(params, req.Params)
	if err != nil {
		return nil, ecode.ErrListUser.Err()
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	records, total, err := h.userDao.GetByColumns(ctx, params)
	if err != nil {
		if strings.Contains(err.Error(), "query params error:") {
			logger.Warn("GetByColumns error", logger.Err(err), logger.Any("params", params), middleware.CtxRequestIDField(ctx))
			return nil, ecode.InvalidParams.Err()
		}
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("params", params), middleware.CtxRequestIDField(ctx))
		return nil, ecode.InternalServerError.Err()
	}

	users := []*schoolV1.User{}
	for _, record := range records {
		data, err := convertUser(record)
		if err != nil {
			logger.Warn("convertUser error", logger.Err(err), logger.Any("id", record.ID), middleware.CtxRequestIDField(ctx))
			continue
		}
		users = append(users, data)
	}

	return &schoolV1.ListUserReply{
		Total: total,
		Users: users,
	}, nil
}

// DeleteByIDs delete records by batch id
func (h *userHandler) DeleteByIDs(ctx context.Context, req *schoolV1.DeleteUserByIDsRequest) (*schoolV1.DeleteUserByIDsReply, error) {
	err := req.Validate()
	if err != nil {
		logger.Warn("req.Validate error", logger.Err(err), logger.Any("req", req), middleware.CtxRequestIDField(ctx))
		return nil, ecode.InvalidParams.Err()
	}

	err = h.userDao.DeleteByIDs(ctx, req.Ids)
	if err != nil {
		logger.Warn("DeleteByIDs error", logger.Err(err), middleware.CtxRequestIDField(ctx))
		return nil, ecode.InternalServerError.Err()
	}

	return &schoolV1.DeleteUserByIDsReply{}, nil
}

// GetByCondition get a record by condition
func (h *userHandler) GetByCondition(ctx context.Context, req *schoolV1.GetUserByConditionRequest) (*schoolV1.GetUserByConditionReply, error) {
	err := req.Validate()
	if err != nil {
		logger.Warn("req.Validate error", logger.Err(err), logger.Any("req", req), middleware.CtxRequestIDField(ctx))
		return nil, ecode.InvalidParams.Err()
	}

	conditions := &query.Conditions{}
	for _, v := range req.Conditions.GetColumns() {
		column := query.Column{}
		_ = copier.Copy(&column, v)
		conditions.Columns = append(conditions.Columns, column)
	}
	err = conditions.CheckValid()
	if err != nil {
		logger.Warn("Parameters error", logger.Err(err), logger.Any("conditions", conditions), middleware.CtxRequestIDField(ctx))
		return nil, ecode.InvalidParams.Err()
	}

	record, err := h.userDao.GetByCondition(ctx, conditions)
	if err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			logger.Warn("GetByID error", logger.Err(err), logger.Any("req", req), middleware.CtxRequestIDField(ctx))
			return nil, ecode.NotFound.Err()
		}
		logger.Error("GetByID error", logger.Err(err), logger.Any("req", req), middleware.CtxRequestIDField(ctx))
		return nil, ecode.InternalServerError.Err()
	}

	data, err := convertUser(record)
	if err != nil {
		logger.Warn("convertUser error", logger.Err(err), logger.Any("user", record), middleware.CtxRequestIDField(ctx))
		return nil, ecode.ErrGetByIDUser.Err()
	}

	return &schoolV1.GetUserByConditionReply{
		User: data,
	}, nil
}

// ListByIDs list of records by batch id
func (h *userHandler) ListByIDs(ctx context.Context, req *schoolV1.ListUserByIDsRequest) (*schoolV1.ListUserByIDsReply, error) {
	err := req.Validate()
	if err != nil {
		logger.Warn("req.Validate error", logger.Err(err), logger.Any("req", req), middleware.CtxRequestIDField(ctx))
		return nil, ecode.InvalidParams.Err()
	}

	userMap, err := h.userDao.GetByIDs(ctx, req.Ids)
	if err != nil {
		logger.Error("GetByIDs error", logger.Err(err), logger.Any("ids", req.Ids), middleware.CtxRequestIDField(ctx))
		return nil, ecode.InternalServerError.Err()
	}

	users := []*schoolV1.User{}
	for _, id := range req.Ids {
		if v, ok := userMap[id]; ok {
			record, err := convertUser(v)
			if err != nil {
				logger.Warn("convertUser error", logger.Err(err), logger.Any("user", v), middleware.CtxRequestIDField(ctx))
				return nil, ecode.InternalServerError.Err()
			}
			users = append(users, record)
		}
	}

	return &schoolV1.ListUserByIDsReply{
		Users: users,
	}, nil
}

// ListByLastID get records by last id
func (h *userHandler) ListByLastID(ctx context.Context, req *schoolV1.ListUserByLastIDRequest) (*schoolV1.ListUserByLastIDReply, error) {
	err := req.Validate()
	if err != nil {
		logger.Warn("req.Validate error", logger.Err(err), logger.Any("req", req), middleware.CtxRequestIDField(ctx))
		return nil, ecode.InvalidParams.Err()
	}
	if req.LastID == "" {
		req.LastID = database.MaxObjectID
	}

	records, err := h.userDao.GetByLastID(ctx, req.LastID, int(req.Limit), req.Sort)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("req", req), middleware.CtxRequestIDField(ctx))
		return nil, ecode.InternalServerError.Err()
	}

	users := []*schoolV1.User{}
	for _, record := range records {
		data, err := convertUser(record)
		if err != nil {
			logger.Warn("convertUser error", logger.Err(err), logger.Any("id", record.ID), middleware.CtxRequestIDField(ctx))
			continue
		}
		users = append(users, data)
	}

	return &schoolV1.ListUserByLastIDReply{
		Users: users,
	}, nil
}

// Login login
func (h *userHandler) Login(ctx context.Context, req *schoolV1.LoginRequest) (*schoolV1.LoginResult, error) {
	c, ctx := middleware.AdaptCtx(ctx)
	c.ShouldBindJSON(req)
	secretKey := []byte("")

	// JWT 载荷
	claims := jwt.MapClaims{
		"username": "vben",
		"password": "123456",
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // 24小时后过期
	}

	// 生成 Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := token.SignedString(secretKey)
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"id":          0,
			"password":    "123456",
			"realName":    "Vben",
			"roles":       []string{"super"},
			"username":    "vben",
			"accessToken": signedToken,
		},
		"error":   nil,
		"message": "ok",
	})
	return nil, errcode.SkipResponse
}

// Info ......
func (h *userHandler) Info(ctx context.Context, req *schoolV1.LoginRequest) (*schoolV1.LoginResult, error) {
	c, ctx := middleware.AdaptCtx(ctx)
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"id":       0,
			"realName": "Vben",
			"roles":    []string{"super"},
			"username": "vben",
		},
		"error":   nil,
		"message": "ok",
	})
	return nil, errcode.SkipResponse
}

// Codes ......
func (h *userHandler) Codes(ctx context.Context, req *schoolV1.LoginRequest) (*schoolV1.LoginResult, error) {
	c, ctx := middleware.AdaptCtx(ctx)
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": []string{
			"AC_100100",
			"AC_100110",
			"AC_100120",
			"AC_100010",
		},
		"error":   nil,
		"message": "ok",
	})
	return nil, errcode.SkipResponse
}

// Logout ......
func (h *userHandler) Logout(ctx context.Context, req *schoolV1.LoginRequest) (*schoolV1.LoginResult, error) {
	c, ctx := middleware.AdaptCtx(ctx)
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"data":    "",
		"error":   nil,
		"message": "ok",
	})
	return nil, errcode.SkipResponse
}

func convertUser(record *model.User) (*schoolV1.User, error) {
	value := &schoolV1.User{}
	err := copier.Copy(value, record)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here, e.g. CreatedAt, UpdatedAt
	value.Id = record.ID.Hex()

	return value, nil
}
