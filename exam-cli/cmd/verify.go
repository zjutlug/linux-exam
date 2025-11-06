package cmd

import (
	"fmt"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"
)

type response struct {
	Code int    `json:"code"`
	Msg  string `json:"message"`
}

var verifyCmd = &cobra.Command{
	Use:   "verify <container_id> <problem>",
	Short: "答案校验",
	Args:  cobra.ExactArgs(2), // 限定必须有两个参数
	Run: func(cmd *cobra.Command, args []string) {
		problemID, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("题目不存在")
			return
		}
		answer := args[1]

		url := conf.APIBaseURL + "/api/user/verify"

		client := resty.New()
		var resp response
		_, err = client.R().
			SetBody(map[string]interface{}{
				"container_id": conf.ContainerId,
				"problem_id":   problemID,
				"answer":       answer,
			}).
			SetResult(&resp).
			Post(url)

		if err != nil {
			fmt.Println("请求出错:", err)
			return
		}
		fmt.Println(resp.Msg)
	},
}

func init() {
	rootCmd.AddCommand(verifyCmd)
}
