package user

import (
	"errors"
	"reflect"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zjutjh/mygo/foundation/reply"
	"github.com/zjutjh/mygo/kit"
	"github.com/zjutjh/mygo/ndb"
	"github.com/zjutjh/mygo/nlog"
	"github.com/zjutjh/mygo/swagger"
	"gorm.io/gorm"

	"app/comm"
	"app/dao/query"
)

// InfoHandler API router注册点
func InfoHandler() gin.HandlerFunc {
	api := InfoApi{}
	swagger.CM[runtime.FuncForPC(reflect.ValueOf(hfInfo).Pointer()).Name()] = api
	return hfInfo
}

type InfoApi struct {
	Info     struct{}        `name:"获取用户信息" desc:"获取用户的得分、提交时间等信息"`
	Request  InfoApiRequest  // API请求参数 (Uri/Header/Query/Body)
	Response InfoApiResponse // API响应数据 (Body中的Data部分)
}

type InfoApiRequest struct {
	Query struct {
		ContainerID string `form:"container_id" binding:"required"` // 容器ID
	}
}

type ProblemResp struct {
	ID    int64  `json:"id"`    // 题目ID
	Score int    `json:"score"` // 分数
	Name  string `json:"name"`  // 题目名称
}
type InfoApiResponse struct {
	Username          string         `json:"username"`           // 用户名
	TotalScore        int            `json:"total_score"`        // 总分
	LastSubmitTime    time.Time      `json:"last_submit"`        // 最后提交时间
	CompletedProblems []*ProblemResp `json:"completed_problems"` // 已完成的题目列表
	AllProblems       []*ProblemResp `json:"all_problems"`       // 所有题目列表
}

// Run Api业务逻辑执行点
func (i *InfoApi) Run(ctx *gin.Context) kit.Code {
	c := ctx.Request.Context()
	q := query.Use(ndb.Pick())

	// 获取用户信息
	user, err := q.User.WithContext(c).Where(q.User.ContainerID.Eq(i.Request.Query.ContainerID)).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		nlog.Pick().WithContext(c).Error("容器%s未注册", i.Request.Query.ContainerID)
		return comm.CodeContainerNotRegistered
	}
	if err != nil {
		nlog.Pick().WithContext(c).Error("获取用户信息失败:", err)
		return comm.CodeDatabaseError
	}

	// 获取用户的所有提交记录
	submissions, err := q.Submission.WithContext(c).Where(q.Submission.ContainerID.Eq(i.Request.Query.ContainerID)).Find()
	if err != nil {
		nlog.Pick().WithContext(c).Error("获取提交记录失败:", err)
		return comm.CodeDatabaseError
	}

	// 初始化响应数据
	i.Response = InfoApiResponse{
		Username:          user.Username,
		TotalScore:        0,
		LastSubmitTime:    user.UpdatedAt,
		CompletedProblems: make([]*ProblemResp, 0, len(submissions)),
		AllProblems:       make([]*ProblemResp, 0, len(comm.BizConf.Problems)),
	}

	// 计算总分和更新最后提交时间
	for _, submission := range submissions {
		// 更新最后提交时间
		if submission.UpdatedAt.After(i.Response.LastSubmitTime) {
			i.Response.LastSubmitTime = submission.UpdatedAt
		}

		// 计算分数并记录完成的题目
		if p, exists := comm.ProblemMap[submission.QuestionID]; exists {
			i.Response.TotalScore += p.Score
			i.Response.CompletedProblems = append(i.Response.CompletedProblems, &ProblemResp{
				ID:    p.Id,
				Score: p.Score,
				Name:  p.Name,
			})
		}
	}

	// 添加所有题目信息
	for _, p := range comm.BizConf.Problems {
		i.Response.AllProblems = append(i.Response.AllProblems, &ProblemResp{
			ID:    p.Id,
			Score: p.Score,
			Name:  p.Name,
		})
	}

	return comm.CodeOK
}

// Init Api初始化 进行参数校验和绑定
func (i *InfoApi) Init(ctx *gin.Context) (err error) {
	err = ctx.ShouldBindQuery(&i.Request.Query)
	return err
}

// hfInfo API执行入口
func hfInfo(ctx *gin.Context) {
	api := &InfoApi{}
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
