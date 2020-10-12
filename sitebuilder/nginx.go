package sitebuilder

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/spf13/viper"
)

type NginxSiteBuilder struct {
	sitesPath string
	template  *template.Template
}

func (sb *NginxSiteBuilder) Generate(label, root string) (string, error) {
	serverName := fmt.Sprintf("%s.hyper-cas.org", label)

	data := struct {
		RootPath   string
		ServerName string
	}{
		RootPath:   fmt.Sprintf("/app/sites/%s", root),
		ServerName: serverName,
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

    location / {
        try_files $uri $uri/ /index.html;
    }

    error_page 404 /404.html;
    error_page 500 502 503 504 /50x.html;
}
`
	return template.New("conf").Parse(tmpl)
}

func NewNginxSiteBuilder() (*NginxSiteBuilder, error) {
	sitesPath := viper.GetString("storage.sitesPath")
	tmpl, err := getConfTemplate()
	if err != nil {
		return nil, err
	}
	return &NginxSiteBuilder{
		sitesPath: sitesPath,
		template:  tmpl,
	}, nil
}
