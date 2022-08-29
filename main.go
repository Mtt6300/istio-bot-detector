package main

import (
	"istio-botdetector/config"
	"istio-botdetector/detector"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

func main() {
	proxywasm.SetVMContext(&vmContext{})
}

type vmContext struct {
	types.DefaultVMContext
}

func (*vmContext) NewPluginContext(contextID uint32) types.PluginContext {
	return &pluginContext{}
}

type pluginContext struct {
	types.DefaultPluginContext
	contextID     uint32
	Configuration config.PluginConfiguration
	Detector      detector.Detector
	CacheBucket   detector.CacheBucket
}

func (ctx *pluginContext) NewHttpContext(contextID uint32) types.HttpContext {
	return &httpContext{Configuration: ctx.Configuration, Detector: ctx.Detector, CacheBucket: ctx.CacheBucket}
}

type httpContext struct {
	types.DefaultHttpContext
	contextID     uint32
	Configuration config.PluginConfiguration
	Detector      detector.Detector
	CacheBucket   detector.CacheBucket
}

func (ctx *pluginContext) OnPluginStart(pluginConfigurationSize int) types.OnPluginStartStatus {
	configDataByte, err := proxywasm.GetPluginConfiguration()
	if err != nil {
		proxywasm.LogCriticalf("error reading plugin configuration: %v", err)
		return types.OnPluginStartStatusFailed
	}

	if config, err := config.ParsPluginConfiguration(configDataByte); err != nil {
		proxywasm.LogCriticalf("error parsing plugin configuration: %v", err)
		return types.OnPluginStartStatusFailed
	} else {
		ctx.Configuration = config
	}

	if d, err := detector.InitializeDetector(ctx.Configuration); err != nil {
		proxywasm.LogCriticalf("error initializing detector: %v", err)
		return types.OnPluginStartStatusFailed
	} else {
		ctx.Detector = d
	}

	if cacheB, err := detector.InitializeCacheBucket(ctx.Configuration); err != nil {
		proxywasm.LogCriticalf("error initializing cache bucket: %v", err)
		return types.OnPluginStartStatusFailed
	} else {
		ctx.CacheBucket = cacheB
	}

	return types.OnPluginStartStatusOK
}

func (ctx *httpContext) OnHttpRequestHeaders(numHeaders int, endOfStream bool) types.Action {
	userAgent, _ := proxywasm.GetHttpRequestHeader("user-agent")

	if detector.IsBot(userAgent, ctx.Detector, ctx.Configuration, ctx.CacheBucket) {
		if err := proxywasm.SendHttpResponse(403, [][2]string{}, []byte("Bot is not allowed"), -1); err != nil {
			proxywasm.LogErrorf("failed to send the 403 response: %v", err)
		}
		return types.ActionPause
	}

	return types.ActionContinue
}
