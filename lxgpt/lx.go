package lxgpt

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/yongxinz/lanxinplus-openapi-go-sdk/sdk"
)

type lxClient struct {
	cli        *sdk.ClientWithResponses
	metricsIns IMetrics

	// send message
	appID     string
	appSecret string
	orgID     string

	// webhook bot
	hookToken  string
	hookSecret string
}

func newLxClient(cli *sdk.ClientWithResponses, metricsIns IMetrics, appID, appSecret, orgID, hookToken, hookSecret string) *lxClient {
	return &lxClient{
		cli:        cli,
		metricsIns: metricsIns,

		appID:     appID,
		appSecret: appSecret,
		orgID:     orgID,

		hookToken:  hookToken,
		hookSecret: hookSecret,
	}
}

func (l *lxClient) sendText(ctx context.Context, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	mails := r.PostForm["mails"]
	msg := r.PostFormValue("msg")

	appToken := l.GetV1AppToken()
	userIds := l.GetUserIdList(appToken, mails)

	params := sdk.V1MessagesCreateParams{}
	params.SetAppToken(appToken)

	type Reminder struct {
		All     bool
		UserIds []*string
	}

	type Text struct {
		Content  string
		Reminder Reminder
	}

	type MsgData struct {
		Text Text
	}

	msgData := MsgData{
		Text: Text{
			Content: msg,
			Reminder: Reminder{
				All:     false,
				UserIds: userIds,
			},
		},
	}

	body := sdk.V1MessagesCreateRequestBody{}
	body.SetAccountId("").
		SetAttach("").
		SetDepartmentIdList([]*string{}).
		SetEntryId("").
		SetMsgData(msgData).
		SetMsgType("text").
		SetUserIdList(userIds)

	reqEditors := []sdk.RequestEditorFn{}

	resp, err := l.cli.V1MessagesCreateWithBodyWithResponse(ctx, &params, body, reqEditors...)
	if err != nil || resp.GetErrCode() != 0 {
		l.metricsIns.EmitLxApiFailed()
		log.Println("LxAPI 调用失败 请稍后重试. ", err)
	} else {
		l.metricsIns.EmitLxApiSuccess()
	}
	return err
}

func (l *lxClient) WebHook(ctx context.Context, t *http.Request) error {
	if err := t.ParseForm(); err != nil {
		return err
	}

	msg := t.PostFormValue("msg")

	sign := l.GenSign()
	params := sdk.V1BotHookMessagesCreateParams{}
	params.SetHookToken(l.hookToken)

	type Text struct {
		Content string
	}

	type MsgData struct {
		Text Text
	}

	msgData := MsgData{
		Text: Text{
			Content: msg,
		},
	}
	timestamp := fmt.Sprintf("%v", time.Now().Unix())
	body := sdk.V1BotHookMessagesCreateRequestBody{}
	body.SetTimestamp(timestamp).SetSign(sign).SetMsgType("text").SetMsgData(msgData)

	reqEditors := []sdk.RequestEditorFn{}

	resp, err := l.cli.V1BotHookMessagesCreateWithBodyWithResponse(ctx, &params, body, reqEditors...)
	if err != nil || resp.GetErrCode() != 0 {
		l.metricsIns.EmitLxApiFailed()
		log.Println("LxAPI 调用失败 请稍后重试. ", err)
	} else {
		l.metricsIns.EmitLxApiSuccess()
	}
	return err
}

func (l *lxClient) GetV1AppToken() string {
	params := sdk.V1AppTokenCreateParams{}
	params.SetGrantType("client_credential").
		SetAppid(l.appID).
		SetSecret(l.appSecret)

	resp, err := l.cli.V1AppTokenCreateWithResponse(context.TODO(), &params)
	if err != nil {
		panic(err)
	}
	if resp == nil ||
		resp.GetErrCode() != 0 ||
		resp.GetData() == nil ||
		resp.GetData().GetAppToken() == "" {
		panic("invalid app_token")
	}
	return resp.GetData().GetAppToken()
}

func (l *lxClient) GetUserIdList(token string, mails []string) (data []*string) {
	for _, mail := range mails {
		params := sdk.V2StaffsIdMappingFetchParams{}
		params.SetAppToken(token).SetIdType("mail").SetIdValue(mail).SetOrgId(l.orgID)

		resp, err := l.cli.V2StaffsIdMappingFetchWithResponse(context.TODO(), &params)
		if err != nil {
			panic(err)
		}
		data = append(data, resp.Data.StaffId)
	}

	return
}

func (l *lxClient) GenSign() string {
	timestamp := time.Now().Unix()
	secret := l.hookSecret

	stringToSign := fmt.Sprintf("%v", timestamp) + "@" + secret
	h := hmac.New(sha256.New, []byte(stringToSign))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	return signature
}
