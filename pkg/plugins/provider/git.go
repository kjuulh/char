package provider

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"git.front.kjuulh.io/kjuulh/char/pkg/schema"
	"golang.org/x/sync/errgroup"
)

type GitPluginProvider struct{}

func NewGitPluginProvider() *GitPluginProvider {
	return &GitPluginProvider{}
}

func (gpp *GitPluginProvider) FetchPlugins(ctx context.Context, registry string, plugins schema.CharSchemaPlugins) error {
	errgroup, ctx := errgroup.WithContext(ctx)
	baseDir := ".char/plugins"
	if os.Getenv("CHAR_DEV_MODE") == "true" {
		if err := os.RemoveAll(baseDir); err != nil {
			return err
		}
	}
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return fmt.Errorf("path already exists cannot create: %w", err)
	}

	for n, plugin := range plugins {
		n, plugin := n, plugin
		errgroup.Go(func() error {
			dest := fmt.Sprintf(
				"%s/%s",
				strings.TrimRight(baseDir, "/"), n.Hash(),
			)
			if _, err := os.Stat(dest); errors.Is(err, os.ErrNotExist) {
				log.Printf("fetching git plugin repo: %s", n)
				return gpp.FetchPlugin(
					ctx,
					registry,
					plugin,
					dest,
				)
			}

			return nil

		})
	}

	if err := errgroup.Wait(); err != nil {
		return err
	}

	return nil
}

func (gpp *GitPluginProvider) FetchPlugin(ctx context.Context, registry string, plugin *schema.CharSchemaPlugin, dest string) error {
	cloneUrl, err := plugin.Opts.GetCloneUrl(ctx, registry, &schema.CloneUrlOpt{
		Protocol: schema.GitProtocolSsh,
		SshUser:  "git",
	},
	)
	if err != nil {
		return err
	}

	output, err := exec.Command(
		"git",
		"clone",
		"--depth=1",
		cloneUrl,
		dest,
	).CombinedOutput()
	if len(output) > 0 {
		log.Print(string(output))
	}
	if err != nil {
		return err
	}
	return nil
}
