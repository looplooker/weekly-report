package server

import "github.com/looplooker/weekly-report/internal/provider"

type Ai interface {
	GetReport(paths, command string) string
}

func NewAi(p string) Ai {
	switch p {
	case "chatglm":
		return provider.NewGlm()
	case "deepseek":
		return provider.NewDeep()
	}

	return nil
}
