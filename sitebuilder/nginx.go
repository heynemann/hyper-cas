package sitebuilder

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/spf13/viper"
)

type NginxSiteBuilder struct {
	serverName string
	sitesPath  string
	template   *template.Template
}

func (sb *NginxSiteBuilder) Generate(label, root string) (string, error) {
	serverName := fmt.Sprintf("%s.%s", label, sb.serverName)
	useBrotli := viper.GetBool("nginx.useBrotli")

	brotli := ""
	if useBrotli {
		brotli = `
	# brotli
	brotli on;
	brotli_comp_level 6;
	brotli_types text/xml image/svg+xml application/x-font-ttf image/vnd.microsoft.icon application/x-font-opentype application/json font/eot application/vnd.ms-fontobject application/javascript font/otf application/xml application/xhtml+xml text/javascript  application/x-javascript text/plain application/x-font-truetype application/xml+rss image/x-icon font/opentype text/css image/x-win-bitmap;`
	}

	data := struct {
		Label      string
		Hash       string
		RootPath   string
		ServerName string
		Brotli     string
	}{
		Label:      label,
		Hash:       root,
		RootPath:   fmt.Sprintf("/app/sites/%s", root),
		ServerName: serverName,
		Brotli:     brotli,
	}

	var tpl bytes.Buffer
	err := sb.template.Execute(&tpl, data)
	if err != nil {
		return "", err
	}
	return tpl.String(), nil
}

func getConfTemplate() (*template.Template, error) {
	const tmpl = `
server {
	root {{.RootPath}};
	index index.html index.htm;
	server_name {{.ServerName}};
	add_header Hyper-Cas-Label {{.Label}};
	add_header Hyper-Cas-Hash {{.Hash}};
	add_header Vary Hyper-Cas-Hash;
	
	gzip on;
	gzip_vary on;
	gzip_proxied any;
	gzip_comp_level 6;
	gzip_types text/plain text/css text/xml application/json application/javascript application/xml+rss application/atom+xml image/svg+xml;
	{{.Brotli}}

    location / {
        try_files $uri $uri/ /index.html;
    }

	include {{.RootPath}}/nginx/*.conf;

    error_page 404 /404.html;
    error_page 500 502 503 504 /50x.html;
}
`
	return template.New("conf").Parse(tmpl)
}

func NewNginxSiteBuilder() (*NginxSiteBuilder, error) {
	viper.SetDefault("useBrotli", false)
	sitesPath := viper.GetString("storage.sitesPath")
	serverName := viper.GetString("nginx.serverName")
	tmpl, err := getConfTemplate()
	if err != nil {
		return nil, err
	}
	return &NginxSiteBuilder{
		sitesPath:  sitesPath,
		serverName: serverName,
		template:   tmpl,
	}, nil
}
