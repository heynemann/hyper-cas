package serve

import (
	"fmt"

	"github.com/valyala/fasthttp"
	"github.com/vtex/hyper-cas/utils"
	"go.uber.org/zap"
)

type LabelHandler struct {
	App *App
}

func NewLabelHandler(app *App) *LabelHandler {
	return &LabelHandler{App: app}
}

func (handler *LabelHandler) handlePut(ctx *fasthttp.RequestCtx) error {
	label := string(ctx.PostArgs().Peek("label"))
	hash := string(ctx.PostArgs().Peek("hash"))
	logger := utils.LoggerWith(zap.String("label", label), zap.String("hash", hash))
	if label == "" || hash == "" {
		err := fmt.Errorf("Both label and hash must be set (label: '%s', hash: '%s')", label, hash)
		logger.Error("Failed to save label.", zap.Error(err))
		return err
	}
	err := handler.App.Storage.StoreLabel(label, hash)
	if err != nil {
		logger.Error("Failed to store label.", zap.Error(err))
		return err
	}
	logger.Debug("Label stored successfully.")
	return nil
}

func (handler *LabelHandler) handleGet(ctx *fasthttp.RequestCtx) error {
	label := ctx.UserValue("label").(string)
	logger := utils.LoggerWith(zap.String("label", label))
	if !handler.App.Storage.HasLabel(label) {
		logger.Info("Label was not found in storage.")
		ctx.SetStatusCode(404)
		return nil
	}
	contents, err := handler.App.Storage.GetLabel(label)
	if err != nil {
		logger.Error("Could not retrieve label from storage.")
		return err
	}
	ctx.SetBodyString(contents)
	logger.Debug("Label retrieved successfully.")
	return nil
}

func (handler *LabelHandler) handleHead(ctx *fasthttp.RequestCtx) error {
	label := ctx.UserValue("label").(string)
	logger := utils.LoggerWith(zap.String("label", label))
	if has := handler.App.Storage.HasLabel(label); has {
		logger.Debug("Label found.")
		ctx.SetStatusCode(200)
	} else {
		logger.Debug("Label not found.")
		ctx.SetStatusCode(404)
	}
	return nil
}
