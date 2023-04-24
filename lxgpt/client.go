package lxgpt

import (
	"fmt"
	"net/http"

	"github.com/lanxinplus/lanxinplus-openapi-go-sdk/sdk"
)

type Client struct {
	lxIns      *lxClient
	metricsIns IMetrics
	serverPort string
}

type ClientConfig struct {
	// lx
	LxAPIUrl  string
	AppID     string
	AppSecret string
	OrgID     string

	// server
	ServerPort string
	Metrics    IMetrics
}

func New(config *ClientConfig) *Client {
	res := new(Client)

	res.metricsIns = config.Metrics
	if res.metricsIns == nil {
		res.metricsIns = new(noneMetrics)
	}

	client, _ := sdk.NewClientWithResponses(config.LxAPIUrl)
	res.lxIns = newLxClient(client, res.metricsIns, config.AppID, config.AppSecret, config.OrgID)

	res.serverPort = config.ServerPort

	return res
}

func (r *Client) Start() error {
	http.HandleFunc("/message", func(w http.ResponseWriter, req *http.Request) {
		r.lxIns.sendText(req.Context(), req)
	})

	fmt.Printf("start server: %s ...\n", r.serverPort)
	return http.ListenAndServe(fmt.Sprint("[::]:", r.serverPort), nil)
}
