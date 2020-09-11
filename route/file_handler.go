package route

import (
	"fmt"
	"strings"

	"github.com/valyala/fasthttp"
)

type FileHandler struct {
	App *App
}

func NewFileHandler(app *App) *FileHandler {
	return &FileHandler{App: app}
}

func getDistroAndPath(path string) (string, string, error) {
	if path == "" {
		return "", "", fmt.Errorf("Invalid path. Must contain at least distro (%s).", path)
	}
	var sb strings.Builder
	distro := ""
	for _, char := range path[1:] {
		if distro == "" && char == '/' {
			distro = sb.String()
			sb.Reset()
		}
		sb.WriteRune(char)
	}
	if distro == "" {
		return "", "", fmt.Errorf("Invalid path. Must contain at least distro (%s).", path)
	}
	return distro, sb.String(), nil
}

func (handler *FileHandler) handleGet(ctx *fasthttp.RequestCtx) {
	distroHash, path, err := getDistroAndPath(string(ctx.Path()))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	distro, err := handler.App.Cache.GetDistro(distroHash)
	if err != nil {
		ctx.SetStatusCode(500)
		ctx.SetBodyString("Error getting distro from cache.")
		return
	}
	if distro == nil {
		distroFiles, err := handler.App.Storage.GetDistro(distroHash)
		if err != nil {
			ctx.SetStatusCode(500)
			ctx.SetBodyString("Error getting distro.")
			return
		}

		distro, err = handler.buildDistro(distroFiles)
		if err != nil {
			ctx.SetStatusCode(500)
			ctx.SetBodyString("Error building distro.")
			return
		}

		err = handler.App.Cache.SetDistro(distroHash, distro)
		if err != nil {
			ctx.SetStatusCode(500)
			ctx.SetBodyString("Error storing distro in cache.")
			return
		}
	}

	path = strings.TrimLeft(path, "/")
	if val, ok := distro[path]; ok {
		ctx.SetBodyString(val)
	} else {
		ctx.SetStatusCode(404)
	}
}

func (h *FileHandler) buildDistro(paths []string) (map[string]string, error) {
	r := map[string]string{}
	for _, p := range paths {
		parts := strings.Split(p, ":")
		path := parts[0]
		hash := parts[1]
		f, err := h.App.Storage.Get(hash)
		if err != nil {
			return nil, err
		}
		r[path] = string(f)
	}
	return r, nil
}
