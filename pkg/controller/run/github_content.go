package run

import (
	"context"
	"fmt"
)

func (c *Controller) getGitHubContent(ctx context.Context, content *GitHubContent) (*TemplateInput, error) {
	s, err := c.gh.GetContent(ctx, content.Owner, content.Repo, content.Path, content.Ref)
	if err != nil {
		return nil, fmt.Errorf("get a file by GitHub Content API: %w", err)
	}
	return &TemplateInput{
		Type:    "github-content",
		Content: s,
		Owner:   content.Owner,
		Repo:    content.Repo,
		Path:    content.Path,
		Ref:     content.Ref,
	}, nil
}
