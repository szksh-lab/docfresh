// Package run implements the business logic for the 'docfresh run' command.
package run

import (
	"context"
	"io"
	"net/http"
	"os"

	"github.com/spf13/afero"
)

// Controller manages the initialization of docfresh configuration.
// It provides methods to create configuration files with appropriate permissions.
type Controller struct {
	fs         afero.Fs
	httpClient *http.Client
	gh         GitHub
	stderr     io.Writer
}

type GitHub interface {
	GetContent(ctx context.Context, owner, repo, path, ref string) (string, error)
}

// New creates a new Controller instance with the provided filesystem and environment.
// The filesystem is used for all file operations, allowing for easy testing with mock filesystems.
func New(fs afero.Fs, gh GitHub) *Controller {
	return &Controller{
		fs:         fs,
		httpClient: http.DefaultClient, // TODO Change
		gh:         gh,
		stderr:     os.Stderr,
	}
}
