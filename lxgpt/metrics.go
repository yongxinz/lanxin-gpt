package lxgpt

type IMetrics interface {
	// lxgpt api
	EmitChatGPTApiFailed()
	EmitChatGPTApiSuccess()

	// lx api
	EmitLxApiFailed()
	EmitLxApiSuccess()

	// app
	EmitAppSuccess()
	EmitAppFailed()
}

type noneMetrics struct{}

func (r *noneMetrics) EmitChatGPTApiFailed() {}

func (r *noneMetrics) EmitChatGPTApiSuccess() {}

func (r *noneMetrics) EmitLxApiFailed() {}

func (r *noneMetrics) EmitLxApiSuccess() {}

func (r *noneMetrics) EmitAppSuccess() {}

func (r *noneMetrics) EmitAppFailed() {}
