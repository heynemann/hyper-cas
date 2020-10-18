package serve

import (
	"fmt"

	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
	"github.com/vtex/hyper-cas/utils"
	"go.uber.org/zap"
)

type FileHandler struct {
	App *App
}

func NewFileHandler(app *App) *FileHandler {
	viper.SetDefault("gzipSourceFiles", false)
	return &FileHandler{App: app}
}

func (handler *FileHandler) handlePut(ctx *fasthttp.RequestCtx) error {
	value := ctx.Request.Body()
	hash := utils.HashBytes(value)
	strHash := fmt.Sprintf("%x", hash)
	if viper.GetBool("gzipSourceFiles") {
		var err error
		value, err = utils.Zip(value)
		if err != nil {
			utils.LogError("Failed to zip source file.", zap.Error(err))
			return err
		}
	}
	err := handler.App.Storage.Store(strHash, value)
	if err != nil {
		utils.LogError("Failed to store file.", zap.ByteString("hash", hash), zap.Error(err))
		return err
	}
	ctx.SetBodyString(strHash)
	utils.LogInfo("Successfully stored file.", zap.ByteString("hash", hash))
	return nil
}

func (handler *FileHandler) handleGet(ctx *fasthttp.RequestCtx) error {
	hash := ctx.UserValue("hash").(string)
	logger := utils.LoggerWith(zap.String("hash", hash))
	contents, err := handler.App.Storage.Get(hash)
	if contents == nil {
		logger.Debug("File not found for specified hash.")
		ctx.SetStatusCode(404)
		return nil
	}
	if err != nil {
		logger.Error("Failed to retrieve file.", zap.Error(err))
		return err
	}
	ctx.SetBody(contents)
	logger.Debug("File retrieved successfully.")
	return nil
}

func (handler *FileHandler) handleHead(ctx *fasthttp.RequestCtx) error {
	hash := ctx.UserValue("hash").(string)
	logger := utils.LoggerWith(zap.String("hash", hash))
	if has := handler.App.Storage.Has(hash); has {
		logger.Debug("File exists.")
		ctx.SetStatusCode(200)
	} else {
		logger.Debug("File not found.")
		ctx.SetStatusCode(404)
	}
	return nil
}
