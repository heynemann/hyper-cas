package sitebuilder

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"

	"github.com/spf13/viper"
)

type NginxSiteBuilder struct {
	serverName string
	sitesPath  string
	template   *template.Template
}

func (sb *NginxSiteBuilder) Generate(label, root string) (string, error) {
	serverName := fmt.Sprintf("%s.%s", label, sb.serverName)

	data := struct {
		Label      string
		Hash       string
		RootPath   string
		ServerName string
	}{
		Label:      label,
		Hash:       root,
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
	path := viper.GetString("nginx.template")
	tmpl, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return template.New("conf").Parse(string(tmpl))
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
