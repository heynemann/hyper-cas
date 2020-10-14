package serve

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"

	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/vtex/hyper-cas/content"
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
	err = handler.App.Storage.StoreDistro(hash, contents)
	if err != nil {
		return err
	}
	ctx.SetBody([]byte(hash))

	return nil
}
