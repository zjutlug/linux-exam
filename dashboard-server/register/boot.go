package register

import (
	"fmt"

	"github.com/zjutjh/mygo/config"
	"github.com/zjutjh/mygo/feishu"
	"github.com/zjutjh/mygo/foundation/kernel"
	"github.com/zjutjh/mygo/kit"
	"github.com/zjutjh/mygo/ndb"
	"github.com/zjutjh/mygo/nlog"

	"app/comm"
	"app/register/generate"
)

func Boot() kernel.BootList {
	return kernel.BootList{
		// 基础引导器
		feishu.Boot(),   // 飞书Bot (消息提醒) 吗 的 不 让 删
		nlog.Boot(),     // 业务日志
		generate.Boot(), // 导入生成代码

		// Client引导器
		ndb.Boot(), // DB

		// 业务引导器
		BizConfBoot(),
		AppBoot(),
	}
}

// BizConfBoot 初始化应用业务配置引导器
func BizConfBoot() func() error {
	return func() error {
		// 解析默认配置
		comm.BizConf = &comm.BizConfig{}
		err := config.Pick().UnmarshalKey("biz", comm.BizConf)
		if err != nil {
			return fmt.Errorf("%w: 解析应用业务配置错误: %w", kit.ErrDataUnmarshal, err)
		}
		// 提取problemMap
		comm.ProblemMap = make(map[int64]comm.ProblemConfig, len(comm.BizConf.Problems))
		for _, p := range comm.BizConf.Problems {
			comm.ProblemMap[p.Id] = p
		}
		return nil
	}
}

// AppBoot 应用定制引导器
func AppBoot() func() error {
	return func() error {
		// 可以在这里编写业务初始引导逻辑
		return nil
	}
}
