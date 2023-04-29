package lxgpt

import (
	"fmt"
	"net/http"

	"github.com/yongxinz/lanxinplus-openapi-go-sdk/sdk"
)

type Client struct {
	lxIns      *lxClient
	metricsIns IMetrics
	serverPort string
}

type ClientConfig struct {
	// lx
	LxAPIUrl   string
	AppID      string
	AppSecret  string
	OrgID      string
	HookToken  string
	HookSecret string

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
	cli.lxIns = newLxClient(client, cli.metricsIns, config.AppID, config.AppSecret, config.OrgID, config.HookToken, config.HookSecret)

	cli.serverPort = config.ServerPort

	return cli
}

func (r *Client) Start() error {
	http.HandleFunc("/message", func(w http.ResponseWriter, req *http.Request) {
		r.lxIns.sendText(req.Context(), req)
	})

	http.HandleFunc("/webhook/message", func(w http.ResponseWriter, req *http.Request) {
		r.lxIns.WebHook(req.Context(), req)
	})

	fmt.Printf("start server: %s ...\n", r.serverPort)
	return http.ListenAndServe(fmt.Sprint("[::]:", r.serverPort), nil)
}
