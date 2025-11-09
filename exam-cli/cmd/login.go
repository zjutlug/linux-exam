package cmd

import (
	"bufio"
	"exam-cli/comm"
	"exam-cli/conf"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"
	"syscall"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:    "login",
	Short:  "登录终端",
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		// 检查是否已经注册
		var response comm.InfoResponse
		client := resty.New()
		_, err := client.R().
			SetQueryParam("container_id", conf.Pick().ContainerId).
			SetResult(&response).
			Post(conf.Pick().APIBaseURL + "/api/user/info")
		if err == nil && response.Code == 0 {
			fmt.Println("欢迎回来, " + response.Data.Username)
			enterShell()
		}
		reader := bufio.NewReader(os.Stdin)
		re := regexp.MustCompile(`^[A-Za-z0-9_-]+$`)
		var name string
		for {
			fmt.Print("请输入用户名（仅字母数字和 _ -，回车确认）: ")
			line, err := reader.ReadString('\n')
			if err != nil {
				_, _ = fmt.Fprintln(os.Stderr, "读取失败：", err)
				os.Exit(1)
			}
			name = strings.TrimSpace(line)
			if name == "" {
				fmt.Println("用户名不能为空，请重试。")
				continue
			}
			if !re.MatchString(name) {
				fmt.Println("用户名只能包含字母数字、下划线或短横线。")
				continue
			}
			if len(name) > 20 {
				fmt.Println("用户名过长")
				continue
			}
			// 向服务器发送注册请求
			client := resty.New()
			params := url.Values{}
			params.Add("username", name)
			params.Add("container_id", conf.Pick().ContainerId)
			resp, err := client.R().
				SetBody(map[string]interface{}{
					"username":     name,
					"container_id": conf.Pick().ContainerId,
				}).
				Post(conf.Pick().APIBaseURL + "/api/user/register")
			msg := string(resp.Body())
			if err != nil || msg != "ok" {
				fmt.Printf("注册失败:%s, 请重试\n", msg)
				continue
			}
			break
		}

		fmt.Printf("欢迎, %s! 正在进入 shell...\n", name)
		enterShell()
	},
}

// enterShell 进入bash shell环境
func enterShell() {
	if err := syscall.Exec("/bin/bash", nil, os.Environ()); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "执行 shell 失败：", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
