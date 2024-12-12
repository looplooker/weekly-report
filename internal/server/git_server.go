package server

import (
	"fmt"
	"os/exec"
	"strings"
)

func getCommit(paths string, command string) string {
	pathList := strings.Fields(paths)
	// 使用 map 来存储唯一的提交信息
	uniqueResults := make(map[string]bool)
	var results []string

	for _, path := range pathList {
		cmdArgs := parseGitCommand(command)

		cmd := exec.Command("git", cmdArgs...)
		cmd.Dir = path

		fmt.Printf("执行目录: %s\n", path)
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
			// 按行分割输出
			lines := strings.Split(strings.TrimSpace(string(output)), "\n")

			// 对每一行进行去重处理
			var uniqueLines []string
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line != "" && !uniqueResults[line] {
					uniqueResults[line] = true
					uniqueLines = append(uniqueLines, line)
				}
			}

			// 只有当有唯一的行时才添加到结果中
			if len(uniqueLines) > 0 {
				result := strings.Join(uniqueLines, "\n")
				//results = append(results, fmt.Sprintf("=== %s ===\n%s", path, result))
				results = append(results, fmt.Sprintf("\n%s", result))
			}
		}
	}

	return strings.Join(results, "\n\n")
}

// parseGitCommand 解析git命令字符串为参数数组
func parseGitCommand(command string) []string {
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
