package user

import (
	"app/dao/model"
	"errors"
	"reflect"
	"runtime"

	"app/comm"

	"github.com/gin-gonic/gin"
	"github.com/zjutjh/mygo/foundation/reply"
	"github.com/zjutjh/mygo/kit"
	"github.com/zjutjh/mygo/ndb"
	"github.com/zjutjh/mygo/nlog"
	"github.com/zjutjh/mygo/swagger"
	"gorm.io/gorm"
)

// VerifyHandler API router注册点
func VerifyHandler() gin.HandlerFunc {
	api := VerifyApi{}
	swagger.CM[runtime.FuncForPC(reflect.ValueOf(hfVerify).Pointer()).Name()] = api
	return hfVerify
}

type VerifyApi struct {
	Info     struct{}          `name:"答案校验" desc:"校验答案并加分"`
	Request  VerifyApiRequest  // API请求参数 (Uri/Header/Query/Body)
	Response VerifyApiResponse // API响应数据 (Body中的Data部分)
}

type VerifyApiRequest struct {
	Body struct {
		ContainerId string `json:"container_id"`
		ProblemId   int64  `json:"problem_id"`
		Answer      string `json:"answer"`
	}
}

type VerifyApiResponse struct {
	Correct bool `json:"correct"`
}

// Run Api业务逻辑执行点
func (v *VerifyApi) Run(ctx *gin.Context) kit.Code {
	c := ctx.Request.Context()
	// 在此编写具体接口业务逻辑
	problemId := v.Request.Body.ProblemId
	answer := v.Request.Body.Answer

	for _, p := range comm.BizConf.Problems {
		if p.Id == problemId {
			// 答案错误
			if p.Answer != answer {
				v.Response.Correct = false
				return comm.CodeWrongAnswer
			}
			// 答案正确
			v.Response.Correct = true
			// 加分
			err := ndb.Pick().WithContext(c).Create(&model.Submission{
				ContainerID: v.Request.Body.ContainerId,
				QuestionID:  v.Request.Body.ProblemId,
			}).Error
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return comm.CodeSubmitDuplicate
			}
			if err != nil {
				return kit.CodeDatabaseError
			}
			return comm.CodeCorrectAnswer
		}
	}
	return comm.CodeProblemNotFound
}

// Init Api初始化 进行参数校验和绑定
func (v *VerifyApi) Init(ctx *gin.Context) (err error) {
	err = ctx.ShouldBindJSON(&v.Request.Body)
	if err != nil {
		return err
	}
	return err
}

// hfVerify API执行入口
func hfVerify(ctx *gin.Context) {
	api := &VerifyApi{}
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
