package schema

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type CharSchemaPluginName string

func (cspn CharSchemaPluginName) Hash() string {
	bytes := sha256.Sum256([]byte(cspn))
	return hex.EncodeToString(bytes[:])
}

type PluginOps struct {
	Org            string
	RepositoryName string
	Path           string
	Version        string
}

type GitProtocol string

const (
	GitProtocolHttps GitProtocol = "https"
	GitProtocolSsh               = "ssh"
)

type CloneUrlOpt struct {
	Protocol GitProtocol
	SshUser  string
}

func (po *PluginOps) GetCloneUrl(ctx context.Context, registry string, opt *CloneUrlOpt) (string, error) {
	if opt == nil {
		return "", errors.New("opt is required")
	}
	switch opt.Protocol {
	case GitProtocolHttps:
		return fmt.Sprintf("https://%s/%s/%s.git", registry, po.Org, po.RepositoryName), nil
	case GitProtocolSsh:
		return fmt.Sprintf("%s@%s:%s/%s.git", opt.SshUser, registry, po.Org, po.RepositoryName), nil
	default:
		return "", errors.New("protocol not allowed")
	}
}

var memo = map[string]*PluginOps{}

func (cspn CharSchemaPluginName) Get() (*PluginOps, error) {
	if m, ok := memo[string(cspn)]; ok {
		return m, nil
	}

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

	if po.Org == "" || po.RepositoryName == "" {
		return nil, errors.New("could not find org or repository name")
	}

	memo[string(cspn)] = po

	return po, nil
}

type CharSchemaPlugins map[CharSchemaPluginName]*CharSchemaPlugin
type CharSchemaPluginVarName string
type CharSchemaPluginVars map[CharSchemaPluginVarName]string

type CharSchemaPlugin struct {
	Opts *PluginOps
	Vars CharSchemaPluginVars `json:"vars"`
}
