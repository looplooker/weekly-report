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

type DeepServer struct{}

func NewDeep() *DeepServer {
	return &DeepServer{}
}

func (a *DeepServer) GetReport(paths, command string) string {
	commitInfo := commit.GetCommit(paths, command)
	// 获取git提交信息
	fmt.Println("获取到的Git信息：\n", commitInfo)

	apiKey := os.Getenv("DEEPSEEK_KEY") // 替换为你的 DeepSeek API 密钥
	client := NewDeepSeekClient(apiKey)

	// 初始化对话
	messages := []Message{
		{Role: "system", Content: os.Getenv("DEEPSEEK_PROMPT")},
		{Role: "user", Content: commitInfo},
	}

	// 发送对话请求
	response, err := client.Chat("deepseek-chat", messages)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// 获取助手的回复
	assistantMessage := response.Choices[0].Message
	return assistantMessage.Content
}

// DeepSeekClient 是 DeepSeek API 的客户端
type DeepSeekClient struct {
	apiKey string
	client *http.Client
}

// NewDeepSeekClient 创建一个新的 DeepSeekClient 实例
func NewDeepSeekClient(apiKey string) *DeepSeekClient {
	return &DeepSeekClient{
		apiKey: apiKey,
		client: &http.Client{},
	}
}

// Message 表示对话中的一条消息
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest 是发送给 DeepSeek API 的请求结构
type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

// ChatResponse 是 DeepSeek API 返回的响应结构
type ChatResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

// Chat 发送一个对话请求到 DeepSeek API
func (c *DeepSeekClient) Chat(model string, messages []Message) (*ChatResponse, error) {
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
	req, err := http.NewRequest("POST", "https://api.deepseek.com/v1/chat/completions", bytes.NewBuffer(requestBody))
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
