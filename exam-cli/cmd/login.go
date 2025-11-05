package cmd

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Interactive login (container entrypoint)",
	Run: func(cmd *cobra.Command, args []string) {
		// 必须是交互式终端
		if !term.IsTerminal(int(os.Stdin.Fd())) {
			_, _ = fmt.Fprintln(os.Stderr, "需要交互式终端。")
			os.Exit(1)
		}

		// 检查本地注册状态文件
		homeDir, err := os.UserHomeDir()
		if err != nil {
			homeDir = "/root" // 容器中默认使用 root 用户
		}
		registeredFile := filepath.Join(homeDir, ".registered")

		// 如果存在注册文件，直接读取用户名并登录
		if content, err := os.ReadFile(registeredFile); err == nil {
			username := strings.TrimSpace(string(content))
			if username != "" {
				fmt.Printf("欢迎回来, %s!\n", username)
				enterShell(username)
				return
			}
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

			// 向服务器发送注册请求
			client := resty.New()
			params := url.Values{}
			params.Add("username", name)
			params.Add("container_id", conf.ContainerId)

			resp, err := client.R().
				SetQueryParamsFromValues(params).
				Post(conf.APIBaseURL + "/api/user/register")

			if err != nil {
				fmt.Printf("注册失败: %v\n", err)
				continue
			}

			// 检查响应状态码
			if resp.StatusCode() != 200 {
				fmt.Printf(resp.String())
				continue
			}

			// 注册成功，创建本地状态文件
			if err := os.WriteFile(registeredFile, []byte(name), 0644); err != nil {
				fmt.Printf("警告：无法创建本地状态文件：%v\n", err)
			}

			break
		}

		fmt.Printf("欢迎, %s! 正在进入 shell...\n", name)
		enterShell(name)
	},
}

// enterShell 进入bash shell环境
func enterShell(username string) {
	env := append(os.Environ(), "USER="+username)
	bashPath := "/bin/bash" // 容器中的 bash 路径
	args := []string{"bash", "--login"}
	if err := syscall.Exec(bashPath, args, env); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "执行 shell 失败：", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
