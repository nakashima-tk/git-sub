package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// 直近の操作ブランチを取得する関数
func getRecentBranches(limit int) ([]string, error) {
	cmd := exec.Command("git", "reflog", "--format=%gs")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(out.String(), "\n")
	var branches []string
	seen := make(map[string]bool)

	re := regexp.MustCompile(`checkout: moving from .* to (\S+)`)

	for _, line := range lines {
		matches := re.FindStringSubmatch(line)
		if len(matches) == 2 {
			branch := matches[1]
			if !seen[branch] {
				branches = append(branches, branch)
				seen[branch] = true
			}
		}
		if len(branches) >= limit {
			break
		}
	}

	return branches, nil
}

// ユーザーにブランチ選択を促す
func selectBranch(branches []string) (string, error) {
	prompt := promptui.Select{
		Label: "切り替えるブランチを選択してください",
		Items: branches,
	}

	_, result, err := prompt.Run()
	if err != nil {
		return "", err
	}
	return result, nil
}

// Gitのブランチを切り替える
func checkoutBranch(branch string) error {
	cmd := exec.Command("git", "checkout", branch)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func main() {
	var rootCmd = &cobra.Command{
		Use:   "git-branch-selector",
		Short: "直近に操作したブランチを表示し、checkout できる",
		Run: func(cmd *cobra.Command, args []string) {
			branches, err := getRecentBranches(5)
			if err != nil {
				fmt.Println("エラー:", err)
				os.Exit(1)
			}

			if len(branches) == 0 {
				fmt.Println("最近操作したブランチが見つかりません。")
				os.Exit(1)
			}

			branch, err := selectBranch(branches)
			if err != nil {
				fmt.Println("選択エラー:", err)
				os.Exit(1)
			}

			fmt.Printf("ブランチ %s に切り替えます...\n", branch)
			if err := checkoutBranch(branch); err != nil {
				fmt.Println("ブランチ切り替えに失敗しました:", err)
				os.Exit(1)
			}
		},
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
