package lxgpt

import (
	"context"
)

func (r *Client) replyChatGPTMessage(msg string) error {
	return r.lxIns.WebHook(context.Background(), msg)
}
