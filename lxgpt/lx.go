package lxgpt

import (
	"context"
	"log"
	"net/http"

	"github.com/lanxinplus/lanxinplus-openapi-go-sdk/sdk"
)

type lxClient struct {
	cli        *sdk.ClientWithResponses
	metricsIns IMetrics

	appID     string
	appSecret string
	orgID     string
}

func newLxClient(cli *sdk.ClientWithResponses, metricsIns IMetrics, appID, appSecret, orgID string) *lxClient {
	return &lxClient{
		cli:        cli,
		metricsIns: metricsIns,

		appID:     appID,
		appSecret: appSecret,
		orgID:     orgID,
	}
}

func (r *lxClient) sendText(ctx context.Context, t *http.Request) error {
	if err := t.ParseForm(); err != nil {
		return err
	}

	mails := t.PostForm["mails"]
	msg := t.PostFormValue("msg")

	appToken := r.GetV1AppToken()
	userIds := r.GetUserIdList(appToken, mails)

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

	params := sdk.V1MessagesCreateParams{}
	params.SetAppToken(appToken)

	body := sdk.V1MessagesCreateRequestBody{}
	body.SetAccountId("").
		SetAttach("").
		SetDepartmentIdList([]*string{}).
		SetEntryId("").
		SetMsgData(msgData).
		SetMsgType("text").
		SetUserIdList(userIds)

	reqEditors := []sdk.RequestEditorFn{}

	_, err := r.cli.V1MessagesCreateWithBodyWithResponse(ctx, &params, body, reqEditors...)
	if err != nil {
		r.metricsIns.EmitLxApiFailed()
		log.Println("LxAPI 调用失败 请稍后重试. ", err)
	} else {
		r.metricsIns.EmitLxApiSuccess()
	}
	return err
}

func (r *lxClient) GetV1AppToken() string {
	params := sdk.V1AppTokenCreateParams{}
	params.SetGrantType("client_credential").
		SetAppid(r.appID).
		SetSecret(r.appSecret)

	resp, err := r.cli.V1AppTokenCreateWithResponse(context.TODO(), &params)
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

func (r *lxClient) GetUserIdList(token string, mails []string) (data []*string) {
	for _, mail := range mails {
		params := sdk.V2StaffsIdMappingFetchParams{}
		params.SetAppToken(token).SetIdType("mail").SetIdValue(mail).SetOrgId(r.orgID)

		resp, err := r.cli.V2StaffsIdMappingFetchWithResponse(context.TODO(), &params)
		if err != nil {
			panic(err)
		}
		data = append(data, resp.Data.StaffId)
	}

	return
}
