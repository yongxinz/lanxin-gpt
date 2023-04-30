package lxgpt

import (
	"fmt"
	"net/http"

	"github.com/yongxinz/lanxinplus-openapi-go-sdk/sdk"
)

type Client struct {
	lxIns      *lxClient
	chatGPTIns *chatGPTClient
	metricsIns IMetrics
	serverPort string
}

type ClientConfig struct {
	// lx
	LxAPIUrl   string
	AppID      string
	AppSecret  string
	HookToken  string
	HookSecret string

	// ChatGPT
	ChatGPTAPIKey      string
	ChatGPTProxy       string
	ChatGPTProxyEnable bool

	// server
	ServerPort string
	Metrics    IMetrics
}

func New(config *ClientConfig) *Client {
	cli := new(Client)

	cli.metricsIns = config.Metrics
	if cli.metricsIns == nil {
		cli.metricsIns = new(noneMetrics)
	}

	client, _ := sdk.NewClientWithResponses(config.LxAPIUrl)
	cli.lxIns = newLxClient(
		client, cli.metricsIns, config.AppID, config.AppSecret, config.HookToken, config.HookSecret)

	cli.chatGPTIns = newChatGPTClient(
		config.ChatGPTAPIKey, config.ChatGPTProxy, config.ChatGPTProxyEnable, cli.metricsIns)

	cli.serverPort = config.ServerPort

	return cli
}

func (r *Client) Start() error {
	http.HandleFunc("/message", func(w http.ResponseWriter, req *http.Request) {
		r.lxIns.SendText(req.Context(), req)
	})

	http.HandleFunc("/webhook/message", func(w http.ResponseWriter, req *http.Request) {
		r.lxIns.WebHookHandler(req.Context(), req)
	})

	http.HandleFunc("/chat", func(w http.ResponseWriter, req *http.Request) {
		r.LxMessageReceiverHandler(req.Context(), req)
	})

	fmt.Printf("start server: %s ...\n", r.serverPort)
	return http.ListenAndServe(fmt.Sprint("[::]:", r.serverPort), nil)
}
