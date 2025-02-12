package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/looplooker/weekly-report/internal/commit"
	"io"
	"log"
	"net/http"
	"os"
)

type HubServer struct{}

func NewHub() *HubServer {
	return &HubServer{}
}

func (a *HubServer) GetReport(paths, command string) string {
	commitInfo := commit.GetCommit(paths, command)
	// 获取git提交信息
	fmt.Println("获取到的Git信息：\n", commitInfo)

	apiKey := os.Getenv("AIHUB_KEY") // 替换为你的 aihub API 密钥
	client := NewHubClient(apiKey)

	// 初始化对话
	messages := []Message{
		{Role: "system", Content: os.Getenv("AIHUB_PROMPT")},
		{Role: "user", Content: commitInfo},
	}

	// 发送对话请求
	response, err := client.Chat(os.Getenv("AIHUB_MODEL"), messages)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// 获取助手的回复
	assistantMessage := response.Choices[0].Message
	return assistantMessage.Content
}

// HubClient 是 Hub API 的客户端
type HubClient struct {
	apiKey string
	client *http.Client
}

// NewHubClient 创建一个新的 HubClient 实例
func NewHubClient(apiKey string) *HubClient {
	return &HubClient{
		apiKey: apiKey,
		client: &http.Client{},
	}
}

// Chat 发送一个对话请求到 Hub API
func (c *HubClient) Chat(model string, messages []Message) (*ChatResponse, error) {
	// 构造请求体
	request := ChatRequest{
		Model:    model,
		Messages: messages,
		Stream:   false,
	}

	// 将请求体编码为 JSON
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to encode request body: %v", err)
	}

	// 创建 HTTP 请求
	req, err := http.NewRequest("POST", "https://api.openai-hub.com/v1/chat/completions", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// 设置请求头
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d, response: %s", resp.StatusCode, string(body))
	}

	// 解析响应体
	var response ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &response, nil
}
