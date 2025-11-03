package comm

// BizConf 业务配置
var BizConf *BizConfig

type ProblemConfig struct {
	Id     int64
	Answer string
	Score  int
	Name   string
}

type BizConfig struct {
	Problems []ProblemConfig
}

var ProblemMap map[int64]ProblemConfig
