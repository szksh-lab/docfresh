package run

import (
	"context"
	"fmt"
	"path"
)

func (c *Controller) getGitHubContent(ctx context.Context, content *GitHubContent) (*TemplateInput, error) {
	s, err := c.gh.GetContent(ctx, content.Owner, content.Repo, content.Path, content.Ref)
	if err != nil {
		return nil, fmt.Errorf("get a file by GitHub Content API: %w", err)
	}
	s, err = extractRange(s, content.Range)
	if err != nil {
		return nil, fmt.Errorf("extract range from github content: %w", err)
	}
	sl := content.Language
	if sl == "" {
		if lang, ok := c.langs[path.Ext(content.Path)]; ok {
			sl = lang.Language
		}
	}
	return &TemplateInput{
		Type:     "github-content",
		Content:  s,
		Owner:    content.Owner,
		Repo:     content.Repo,
		Path:     content.Path,
		Ref:      content.Ref,
		Language: sl,
	}, nil
}
