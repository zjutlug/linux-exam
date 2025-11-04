package cmd

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"
)

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Username          string `json:"username"`
		TotalScore        int    `json:"total_score"`
		LastSubmit        string `json:"last_submit"`
		CompletedProblems []struct {
			ID    int    `json:"id"`
			Score int    `json:"score"`
			Name  string `json:"name"`
		} `json:"completed_problems"`
		AllProblems []struct {
			ID    int    `json:"id"`
			Score int    `json:"score"`
			Name  string `json:"name"`
		} `json:"all_problems"`
	} `json:"data"`
}

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "获取用户信息",
	Run: func(cmd *cobra.Command, args []string) {
		var response Response
		config := readConfig()
		client := resty.New()
		_, err := client.R().
			SetQueryParam("container_id", config["container_id"]).
			SetResult(&response).
			Post(config["api_base_url"] + "/api/user/info")

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
