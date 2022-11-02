package schema

import (
	"regexp"
	"strings"
)

type CharSchemaPluginName string

type PluginOps struct {
	Org            string
	RepositoryName string
	Path           string
	Version        string
}

func (cspn CharSchemaPluginName) Get() *PluginOps {
	po := &PluginOps{}
	reg := regexp.MustCompile(
		`(?P<org>[\d\w\-_\.]+)\/(?P<repo>[\d\w\-_\.]+)(?P<path>#[\d\w\-_\.\/]+)?(?P<version>@[\d\w\-_\.\/]+)?(?P<path>#[\d\w\-_\.\/]+)?`,
	)
	matches := reg.FindStringSubmatch(string(cspn))
	tags := reg.SubexpNames()

	matchTags := make(map[string]string, len(matches))
	for i, match := range matches {
		tag := tags[i]
		if existingTag, ok := matchTags[tag]; !ok || existingTag == "" {
			matchTags[tag] = match
		}
	}

	if org, ok := matchTags["org"]; ok {
		po.Org = org
	}
	if repo, ok := matchTags["repo"]; ok {
		po.RepositoryName = repo
	}
	if path, ok := matchTags["path"]; ok {
		po.Path = strings.TrimLeft(path, "#")
	}
	if version, ok := matchTags["version"]; ok {
		po.Version = strings.TrimLeft(version, "@")
	}

	return po
}

type CharSchemaPlugins map[CharSchemaPluginName]CharSchemaPlugin

type CharSchemaPlugin struct {
}
