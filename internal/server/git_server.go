package server

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

func getCommit(paths string, command string) string {
	pathList := strings.Fields(paths)
	uniqueResults := make(map[string]bool)
	var results []string

	for _, path := range pathList {
		var projectPaths []string

		// 检查路径是否以 '/' 结尾
		if strings.HasSuffix(path, "/") {
			// 读取目录下的所有子目录
			entries, err := os.ReadDir(path)
			if err != nil {
				fmt.Printf("读取目录错误 %s: %v\n", path, err)
				continue
			}

			// 遍历子目录
			for _, entry := range entries {
				if entry.IsDir() {
					projectPaths = append(projectPaths, filepath.Join(path, entry.Name()))
				}
			}
		} else {
			// 单个项目路径
			projectPaths = append(projectPaths, path)
		}

		// 处理所有项目路径
		for _, projectPath := range projectPaths {
			// 检查是否为 Git 仓库
			if !isGitRepo(projectPath) {
				fmt.Printf("跳过非 Git 仓库: %s\n", projectPath)
				continue
			}

			cmdArgs := parseGitCommand(command)
			cmd := exec.Command("git", cmdArgs...)
			cmd.Dir = projectPath

			fmt.Printf("执行目录: %s\n", projectPath)
			fmt.Printf("执行命令: git %s\n", command)

			output, err := cmd.Output()
			if err != nil {
				if exitErr, ok := err.(*exec.ExitError); ok {
					fmt.Printf("命令执行错误: %v\n", string(exitErr.Stderr))
				} else {
					fmt.Printf("其他错误: %v\n", err)
				}
				continue
			}

			if output != nil {
				lines := strings.Split(strings.TrimSpace(string(output)), "\n")
				var uniqueLines []string
				for _, line := range lines {
					line = strings.TrimSpace(line)
					if line != "" && !uniqueResults[line] {
						uniqueResults[line] = true
						uniqueLines = append(uniqueLines, line)
					}
				}

				if len(uniqueLines) > 0 {
					result := strings.Join(uniqueLines, "\n")
					results = append(results, fmt.Sprintf("\n%s", result))
				}
			}
		}
	}

	return strings.Join(results, "\n\n")
}

// isGitRepo 检查指定路径是否为 Git 仓库
func isGitRepo(path string) bool {
	gitPath := filepath.Join(path, ".git")
	_, err := os.Stat(gitPath)
	return err == nil
}

// parseGitCommand 解析git命令字符串为参数数组
func parseGitCommand(command string) []string {
	// 确保输入字符串是UTF-8编码
	if !utf8.ValidString(command) {
		// 尝试将字符串从系统默认编码转换为UTF-8
		bytes := []byte(command)
		command = string(bytes)
	}

	var args []string
	var current strings.Builder
	inQuote := false
	quoteChar := rune(0)

	// 替换所有的单引号为双引号，简化处理
	command = strings.ReplaceAll(command, "'", "\"")

	for _, char := range command {
		switch {
		case char == '"':
			if !inQuote {
				inQuote = true
				quoteChar = char
			} else if char == quoteChar {
				inQuote = false
				quoteChar = 0
				// 当引号结束时，保存当前参数
				if current.Len() > 0 {
					args = append(args, current.String())
					current.Reset()
				}
			} else {
				current.WriteRune(char)
			}
		case char == ' ' && !inQuote:
			if current.Len() > 0 {
				args = append(args, current.String())
				current.Reset()
			}
		default:
			current.WriteRune(char)
		}
	}

	// 处理最后一个参数
	if current.Len() > 0 {
		args = append(args, current.String())
	}

	// 移除空参数
	var cleanArgs []string
	for _, arg := range args {
		if strings.TrimSpace(arg) != "" {
			cleanArgs = append(cleanArgs, strings.TrimSpace(arg))
		}
	}
	fmt.Printf("解析后的参数: %#v\n", cleanArgs)

	return cleanArgs
}
