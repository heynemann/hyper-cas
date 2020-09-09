package serve

import (
	"bufio"
	"bytes"
	"crypto"
	"fmt"
	"strings"

	"github.com/heynemann/hyper-cas/content"
	routing "github.com/qiangxue/fasthttp-routing"
)

type DistroHandler struct {
	App *App
}

func NewDistroHandler(app *App) *DistroHandler {
	return &DistroHandler{App: app}
}

func (handler *DistroHandler) handlePut(ctx *routing.Context) error {
	value := ctx.Request.Body()
	scanner := bufio.NewScanner(bytes.NewReader(value))

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
	}

	tree, err := content.NewTreeWithHashes(items, crypto.SHA256)
	if err != nil {
		return err
	}
	root := tree.Root()
	ctx.SetBody([]byte(fmt.Sprintf("%x", root.Hash)))

	return nil
}

// func (handler *DistroHandler) handleGet(ctx *routing.Context) error {
// hash := ctx.Param("hash")
// cached, err := handler.App.Cache.Get(hash)
// if err != nil {
// return err
// }
// if cached != nil {
// ctx.SetBody(cached)
// return nil
// }
// contents, err := handler.App.Storage.Get(hash)
// if err != nil {
// return err
// }
// err = handler.App.Cache.Set(hash, contents)
// if err != nil {
// return err
// }
// ctx.SetBody(contents)
// return nil
// }
