package cmd

import (
	"exam-cli/comm"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "获取用户信息",
	Run: func(cmd *cobra.Command, args []string) {
		var response comm.InfoResponse
		client := resty.New()
		_, err := client.R().
			SetQueryParam("container_id", conf.ContainerId).
			SetResult(&response).
			Post(conf.APIBaseURL + "/api/user/info")

		if err != nil {
			fmt.Println("Error sending request:", err)
			return
		}

		fmt.Printf("用户名: %s\n", response.Data.Username)
		fmt.Printf("当前分数: %d\n", response.Data.TotalScore)

		fmt.Println("已提交的题目:")
		for _, p := range response.Data.CompletedProblems {
			fmt.Printf(" - [%d] %s (%d分)\n", p.ID, p.Name, p.Score)
		}

		fmt.Println("所有题目:")
		for _, p := range response.Data.AllProblems {
			fmt.Printf(" - [%d] %s (%d分)\n", p.ID, p.Name, p.Score)
		}
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
