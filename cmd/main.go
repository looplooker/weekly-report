package main

import (
	"flag"
	"fmt"
	"github.com/looplooker/weekly-report/internal/server"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// 定义命令行参数
	paths := flag.String("paths", "C:/dev/xxx", "项目根目录列表，用空格分隔")
	command := flag.String("command", "log --since='1 week ago' --pretty=format:'%h - %s (%cr) <%an>'", "要执行的git命令")
	flag.Parse()

	// 检查必要参数
	if *paths == "" {
		log.Fatal("请指定项目路径")
	}
	if *command == "" {
		log.Fatal("请指定执行命令")
	}

	// 加载 .env 文件
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	sendkey := os.Getenv("SEND_KEY")
	title := "新鲜的周报出炉啦~"
	ai := server.NewAi()
	desp := ai.GetReport(*paths, *command)

	resp, err := server.ScSend(sendkey, title, desp, nil)
	if err != nil {
		log.Fatal("Error:", err)
	} else {
		fmt.Println("Response:", resp)
	}
}
