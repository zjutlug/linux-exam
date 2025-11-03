package user

import (
	"app/dao/model"
	"errors"
	"reflect"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/zjutjh/mygo/foundation/reply"
	"github.com/zjutjh/mygo/kit"
	"github.com/zjutjh/mygo/ndb"
	"github.com/zjutjh/mygo/nlog"
	"github.com/zjutjh/mygo/swagger"
	"gorm.io/gorm"

	"app/comm"
)

// RegisterHandler API router注册点
func RegisterHandler() gin.HandlerFunc {
	api := RegisterApi{}
	swagger.CM[runtime.FuncForPC(reflect.ValueOf(hfRegister).Pointer()).Name()] = api
	return hfRegister
}

type RegisterApi struct {
	Info     struct{}            `name:"用户注册" desc:"API描述"`
	Request  RegisterApiRequest  // API请求参数 (Uri/Header/Query/Body)
	Response RegisterApiResponse // API响应数据 (Body中的Data部分)
}

type RegisterApiRequest struct {
	Query struct {
		Username    string `form:"username"`
		ContainerId string `form:"container_id"`
	}
}

type RegisterApiResponse struct{}

// Run Api业务逻辑执行点
func (r *RegisterApi) Run(ctx *gin.Context) kit.Code {
	c := ctx.Request.Context()
	err := ndb.Pick().WithContext(c).Create(&model.User{
		Username:    r.Request.Query.Username,
		ContainerID: r.Request.Query.ContainerId,
	}).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return comm.CodeUserExist
		}
		nlog.Pick().WithContext(ctx).Error(err)
		return comm.CodeUnknownError
	}
	return comm.CodeOK
}

// Init Api初始化 进行参数校验和绑定
func (r *RegisterApi) Init(ctx *gin.Context) (err error) {
	err = ctx.ShouldBindQuery(&r.Request.Query)
	if err != nil {
		return err
	}
	return err
}

// hfRegister API执行入口
func hfRegister(ctx *gin.Context) {
	api := &RegisterApi{}
	err := api.Init(ctx)
	if err != nil {
		nlog.Pick().WithContext(ctx).WithError(err).Warn("参数绑定校验错误")
		reply.Fail(ctx, comm.CodeParameterInvalid)
		return
	}
	code := api.Run(ctx)
	if !ctx.IsAborted() {
		if code == comm.CodeOK {
			reply.Success(ctx, api.Response)
		} else {
			reply.Fail(ctx, code)
		}
	}
}
