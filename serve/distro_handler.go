package serve

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/valyala/fasthttp"
	"github.com/vtex/hyper-cas/content"
	"github.com/vtex/hyper-cas/utils"
	"go.uber.org/zap"
)

type DistroHandler struct {
	App *App
}

func NewDistroHandler(app *App) *DistroHandler {
	return &DistroHandler{App: app}
}

func (handler *DistroHandler) handlePut(ctx *fasthttp.RequestCtx) error {
	value := ctx.Request.Body()
	scanner := bufio.NewScanner(bytes.NewReader(value))

	contents := []string{}
	items := []content.NodeItem{}
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), ":")
		if len(parts) != 2 {
			err := fmt.Errorf("The body should be composed of lines with {filepath}:{content hash} only.")
			utils.LogError("Failed to parse distribution body.", zap.Error(err))
			return err
		}
		items = append(items, content.NodeItem{
			parts[0],
			[]byte(parts[1]),
		})
		contents = append(contents, fmt.Sprintf("%s:%s", parts[0], parts[1]))
	}

	tree, err := content.NewTreeWithHashes(items)
	if err != nil {
		utils.LogError("Failed to calculate tree for distribution.", zap.Strings("items", contents))
		return err
	}
	root := tree.Root()
	hash := fmt.Sprintf("%x", root.Hash)
	utils.LogDebug("Distribution contents parsed successfully and tree calculated.", zap.String("hash", hash))

	if handler.App.Storage.HasDistro(hash) {
		utils.LogInfo("Distribution already exists on storage. Skipping distribution storage...", zap.String("hash", hash))
		ctx.SetStatusCode(200)
		ctx.SetBodyString(hash)
		return nil
	}
	err = handler.App.Storage.StoreDistro(hash, contents)
	if err != nil {
		utils.LogError("Failed to store distribution.", zap.String("hash", hash), zap.Error(err))
		return err
	}
	ctx.SetBodyString(hash)
	utils.LogDebug("Distribution stored successfully.", zap.String("hash", hash))

	return nil
}

func (handler *DistroHandler) handleGet(ctx *fasthttp.RequestCtx) error {
	distro := ctx.UserValue("distro").(string)
	if !handler.App.Storage.HasDistro(distro) {
		utils.LogInfo("Distribution could not be found in storage.", zap.String("hash", distro))
		ctx.SetStatusCode(404)
		return nil
	}
	contents, err := handler.App.Storage.GetDistro(distro)
	if err != nil {
		utils.LogError("Distribution could not be retrieved from storage.", zap.String("hash", distro), zap.Error(err))
		return err
	}
	items, err := json.Marshal(contents)
	if err != nil {
		utils.LogError("Distribution could not be deserialized from storage.", zap.String("hash", distro), zap.Error(err))
		return err
	}
	ctx.SetBody(items)
	utils.LogDebug("Distribution loaded successfully.", zap.String("hash", distro))
	return nil
}

func (handler *DistroHandler) handleHead(ctx *fasthttp.RequestCtx) error {
	distro := ctx.UserValue("distro").(string)
	if handler.App.Storage.HasDistro(distro) {
		utils.LogDebug("Distribution found.", zap.String("hash", distro))
		ctx.SetStatusCode(200)
	} else {
		utils.LogDebug("Distribution not found.", zap.String("hash", distro))
		ctx.SetStatusCode(404)
	}
	return nil
}
