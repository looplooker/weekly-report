package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type getTokenReq struct {
	ApiKey    string `json:"api_key"`
	ApiSecret string `json:"api_secret"`
}

type tokenRes struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}
type getTokenRes struct {
	Message string   `json:"message"`
	Result  tokenRes `json:"result"`
	Status  int      `json:"status"`
}

type sessionReq struct {
	AssistantId string `json:"assistant_id"`
	Prompt      string `json:"prompt"`
}

type content struct {
	Text string `json:"text"`
	Type string `json:"type"`
}
type part struct {
	Role      string    `json:"role"`
	Content   []content `json:"content"`
	Status    string    `json:"status"`
	CreatedAt string    `json:"created_at"`
}
type sessionResult struct {
	HistoryId      string `json:"history_id"`
	ConversationId string `json:"conversation_id"`
	Output         []part `json:"output"`
	Status         string `json:"status"`
}
type sessionRes struct {
	Message string        `json:"message"`
	Result  sessionResult `json:"result"`
	Status  int           `json:"status"`
}

type AiServer struct{}

func NewAi() *AiServer {
	return &AiServer{}
}

func (a *AiServer) GetReport(paths, command string) string {
	token := getToken()
	commitInfo := getCommit(paths, command)
	// 获取git提交信息
	fmt.Println("获取到的Git信息：\n", commitInfo)

	return startSession(token, commitInfo)
}

func getToken() string {
	url := os.Getenv("AI_API") + os.Getenv("AI_GET_TOKEN")
	params := &getTokenReq{
		ApiKey:    os.Getenv("AI_KEY"),
		ApiSecret: os.Getenv("AI_SECRET"),
	}

	payloadBytes, err := json.Marshal(params)
	if err != nil {
		log.Fatal(err)
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json;charset=utf-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var res getTokenRes
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		log.Fatal(err)
	}
	if res.Status != 0 {
		log.Fatal(res.Message)
	}

	return res.Result.AccessToken
}

func startSession(token, prompt string) string {
	url := os.Getenv("AI_API") + os.Getenv("AI_SESSION")
	params := &sessionReq{
		AssistantId: os.Getenv("ASSISTANT_ID"),
		Prompt:      prompt,
	}

	payloadBytes, err := json.Marshal(params)
	if err != nil {
		log.Fatal(err)
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var res sessionRes
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		log.Fatal(err)
	}
	if res.Status != 0 {
		log.Fatal(res.Message)
	}

	return res.Result.Output[0].Content[0].Text
}
