package sitebuilder

type SiteBuilder interface {
	Generate(label, root string) (string, error)
}
