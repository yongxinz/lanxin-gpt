package lxgpt

import (
	"context"
	"log"
	"net/http"
	"strings"
)

func (r *Client) ReceiveChatGPTMessage(ctx context.Context, msg string) (err error) {
	defer func() {
		if err != nil {
			r.metricsIns.EmitAppFailed()
		} else {
			r.metricsIns.EmitAppSuccess()
		}
	}()

	var result string
	result, err = r.chatGPTIns.ChatGPTRequest(msg)
	if err != nil {
		log.Println("ChatGPT 请求失败, 请稍后重试. ", err)
		return r.replyChatGPTMessage("ChatGPT 请求失败, 请稍后重试.")
	} else if strings.TrimSpace(result) == "" {
		return r.replyChatGPTMessage("ChatGPT 请求失败, 请稍后重试.")
	}

	return r.replyChatGPTMessage(result)
}

func (r *Client) LxMessageReceiverHandler(ctx context.Context, t *http.Request) (string, error) {
	if err := t.ParseForm(); err != nil {
		return "", err
	}

	msg := t.PostFormValue("msg")
	go r.ReceiveChatGPTMessage(ctx, msg)

	return "", nil
}
