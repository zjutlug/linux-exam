package dashboard

import (
	"app/dao/query"
	"reflect"
	"runtime"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zjutjh/mygo/foundation/reply"
	"github.com/zjutjh/mygo/kit"
	"github.com/zjutjh/mygo/ndb"
	"github.com/zjutjh/mygo/nlog"
	"github.com/zjutjh/mygo/swagger"

	"app/comm"
)

// ScoreHandler API router注册点
func ScoreHandler() gin.HandlerFunc {
	api := ScoreApi{}
	swagger.CM[runtime.FuncForPC(reflect.ValueOf(hfScore).Pointer()).Name()] = api
	return hfScore
}

// ScoreApi 定义看板分数API的结构
type ScoreApi struct {
	Info     struct{}         `name:"获取看板分数" desc:"获取看板分数"`
	Request  ScoreApiRequest  // API请求参数 (Uri/Header/Query/Body)
	Response ScoreApiResponse // API响应数据 (Body中的Data部分)
}

// ScoreApiRequest 请求参数结构体
type ScoreApiRequest struct {
}

// ScoreResp 单个用户得分响应结构
type ScoreResp struct {
	Username   string    `json:"username"`    // 用户名
	Score      int       `json:"score"`       // 总分
	LastSubmit time.Time `json:"last_submit"` // 最后提交时间
}

// ScoreApiResponse 响应数据类型
type ScoreApiResponse []ScoreResp

// Run Api业务逻辑执行点
func (s *ScoreApi) Run(ctx *gin.Context) kit.Code {
	c := ctx.Request.Context()

	// 获取所有用户
	q := query.Use(ndb.Pick())
	users, err := q.User.WithContext(c).Find()
	if err != nil {
		nlog.Pick().WithContext(ctx).Error("获取用户数据失败:", err)
		return comm.CodeDatabaseError
	}

	// 获取所有提交记录
	submissions, err := q.Submission.WithContext(c).Find()
	if err != nil {
		nlog.Pick().WithContext(ctx).Error("获取提交记录失败:", err)
		return comm.CodeDatabaseError
	}

	// 创建用户得分统计映射
	userScoreMap := make(map[string]*ScoreResp, len(users))

	// 初始化所有用户的得分记录
	for _, user := range users {
		userScoreMap[user.ContainerID] = &ScoreResp{
			Username:   user.Username,
			Score:      0,
			LastSubmit: time.Time{}, // 初始化为零值时间
		}
	}

	// 创建题目分数查找映射以提高查询效率
	problemScores := make(map[int64]int)
	for _, problem := range comm.BizConf.Problems {
		problemScores[problem.Id] = problem.Score
	}

	// 计算每个用户的得分并记录最后提交时间
	for _, submission := range submissions {
		// 单个问题的分数
		if score, ok := userScoreMap[submission.ContainerID]; ok {
			if problemScore, exists := problemScores[submission.QuestionID]; exists {
				score.Score += problemScore
				// 更新最后提交时间（如果当前提交时间更晚）
				if submission.UpdatedAt.After(score.LastSubmit) {
					score.LastSubmit = submission.UpdatedAt
				}
			}
		}
	}

	// 转换为切片并排序
	resp := make([]ScoreResp, 0, len(userScoreMap))
	for _, score := range userScoreMap {
		resp = append(resp, *score)
	}

	sort.Slice(resp, func(i, j int) bool {
		// 首先按分数降序排序
		if resp[i].Score != resp[j].Score {
			return resp[i].Score > resp[j].Score
		}
		// 分数相同时按最后提交时间升序排序（先提交的排前面）
		return resp[i].LastSubmit.Before(resp[j].LastSubmit)
	})

	s.Response = resp
	return comm.CodeOK
}

// Init Api初始化 进行参数校验和绑定
func (s *ScoreApi) Init(ctx *gin.Context) (err error) {
	return err
}

// hfScore API执行入口
func hfScore(ctx *gin.Context) {
	api := &ScoreApi{}
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
