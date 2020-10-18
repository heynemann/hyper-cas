package serve

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/valyala/fasthttp"
	"github.com/vtex/hyper-cas/content"
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
			return fmt.Errorf("The body should be composed of lines with {filepath}:{content hash} only.")
		}
		items = append(items, content.NodeItem{
			parts[0],
			[]byte(parts[1]),
		})
		contents = append(contents, fmt.Sprintf("%s:%s", parts[0], parts[1]))
	}

	tree, err := content.NewTreeWithHashes(items)
	if err != nil {
		return err
	}
	root := tree.Root()
	hash := fmt.Sprintf("%x", root.Hash)

	if handler.App.Storage.HasDistro(hash) {
		ctx.SetStatusCode(200)
		ctx.SetBodyString(hash)
		return nil
	}
	err = handler.App.Storage.StoreDistro(hash, contents)
	if err != nil {
		return err
	}
	ctx.SetBodyString(hash)

	return nil
}

func (handler *DistroHandler) handleGet(ctx *fasthttp.RequestCtx) error {
	distro := ctx.UserValue("distro").(string)
	if !handler.App.Storage.HasDistro(distro) {
		ctx.SetStatusCode(404)
		return nil
	}
	contents, err := handler.App.Storage.GetDistro(distro)
	if err != nil {
		return err
	}
	items, err := json.Marshal(contents)
	if err != nil {
		return err
	}
	ctx.SetBody(items)
	return nil
}

func (handler *DistroHandler) handleHead(ctx *fasthttp.RequestCtx) error {
	distro := ctx.UserValue("distro").(string)
	if handler.App.Storage.HasDistro(distro) {
		ctx.SetStatusCode(200)
	} else {
		ctx.SetStatusCode(404)
	}
	return nil
}
