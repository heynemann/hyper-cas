package serve

import (
	"fmt"

	routing "github.com/qiangxue/fasthttp-routing"
)

type FileHandler struct {
	App *App
}

func (handler *FileHandler) handlePut(ctx *routing.Context) error {
	key := ctx.FormValue("key")
	value := ctx.FormValue("value")
	fmt.Fprintf(ctx, "%s ==> %s\n", key, value)
	return nil
}
