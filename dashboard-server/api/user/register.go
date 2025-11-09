package user

import (
	"app/dao/model"
	"errors"
	"reflect"
	"runtime"

	"github.com/gin-gonic/gin"
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
	Info     struct{}            `name:"用户注册" desc:"用户注册, 提交用户名和容器ID进行注册"`
	Request  RegisterApiRequest  // API请求参数 (Uri/Header/Query/Body)
	Response RegisterApiResponse // API响应数据 (Body中的Data部分)
}

type RegisterApiRequest struct {
	Body struct {
		Username    string `json:"username" binding:"required"`
		ContainerId string `json:"container_id" binding:"required"`
	}
}

type RegisterApiResponse struct{}

// Run Api业务逻辑执行点
func (r *RegisterApi) Run(ctx *gin.Context) kit.Code {
	c := ctx.Request.Context()
	err := ndb.Pick().WithContext(c).Create(&model.User{
		Username:    r.Request.Body.Username,
		ContainerID: r.Request.Body.ContainerId,
	}).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return comm.CodeUserExist
		}
		nlog.Pick().WithContext(ctx).Error(err)
		return comm.CodeUnknownError
	}
	nlog.Pick().WithContext(c).Infof("容器[%s]注册成功, 用户名为: %s", r.Request.Body.ContainerId, r.Request.Body.Username)
	return comm.CodeOK
}

// Init Api初始化 进行参数校验和绑定
func (r *RegisterApi) Init(ctx *gin.Context) (err error) {
	err = ctx.ShouldBindJSON(&r.Request.Body)
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
		ctx.String(400, "登录失败")
		return
	}
	code := api.Run(ctx)
	if !ctx.IsAborted() {
		if code == comm.CodeOK {
			ctx.String(200, "ok")
		} else if code == comm.CodeUserExist {
			ctx.String(400, "用户已存在")
		} else {
			ctx.String(400, "未知异常")
		}
	}
}
