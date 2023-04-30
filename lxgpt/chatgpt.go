package lxgpt

import (
	"context"
	"log"
	"net/http"
	"net/url"

	"github.com/sashabaranov/go-openai"
)

type chatGPTClient struct {
	apiKey      string
	proxy       string
	proxyEnable bool
	metricsIns  IMetrics
}

func newChatGPTClient(apiKey, proxy string, proxyEnable bool, metricsIns IMetrics) *chatGPTClient {
	return &chatGPTClient{
		apiKey:      apiKey,
		proxy:       proxy,
		proxyEnable: proxyEnable,
		metricsIns:  metricsIns,
	}
}

func (r *chatGPTClient) ChatGPTRequest(msg string) (result string, err error) {
	defer func() {
		if err != nil {
			r.metricsIns.EmitChatGPTApiFailed()
		} else {
			r.metricsIns.EmitChatGPTApiSuccess()
		}
	}()

	log.Print("Receive message: ", msg)

	var client *openai.Client
	if r.proxyEnable {
		config := openai.DefaultConfig(r.apiKey)
		proxyUrl, err := url.Parse(r.proxy)
		if err != nil {
			panic(err)
		}
		transport := &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		}
		config.HTTPClient = &http.Client{
			Transport: transport,
		}

		client = openai.NewClientWithConfig(config)
	} else {
		client = openai.NewClient(r.apiKey)
	}

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: msg,
				},
			},
		},
	)
	if err != nil {
		return "ChatCompletion error", err
	}

	result = resp.Choices[0].Message.Content
	log.Println("msg: ", msg, "result: ", result)

	return
}
