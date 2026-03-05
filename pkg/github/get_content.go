package github

import (
	"context"
	"fmt"
)

func (c *Client) GetContent(ctx context.Context, owner, repo, path, ref string) (string, error) {
	var opts *RepositoryContentGetOptions
	if ref != "" {
		opts = &RepositoryContentGetOptions{
			Ref: ref,
		}
	}
	file, _, _, err := c.repo.GetContents(ctx, owner, repo, path, opts)
	if err != nil {
		return "", fmt.Errorf("get a file content by GitHub API: %w", err)
	}
	s, err := file.GetContent()
	if err != nil {
		return "", fmt.Errorf("read a file content: %w", err)
	}
	return s, nil
}
