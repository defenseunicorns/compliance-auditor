package network

import (
	"context"

	"github.com/hashicorp/go-getter/v2"
)

// DownloadFile downloads a file to the requested dst (full file path) from src.
// wd is used when resolving relative paths and should be set to the common root
// directory of relative paths (defaults to the directory the input file is in).
// Go-getter supports the following protocols, though not all are relevant in
// file mode (and Lula can extend the Getter interface if we need to add to this
// list):
//
//   - Local files
//   - Git
//   - Mercurial
//   - HTTP
//   - Amazon S3
//   - Google GCP
//   - SMB
//   - Source: https://pkg.go.dev/github.com/hashicorp/go-getter/v2#section-readme
//
// DownloadFile returns the actual download path and error. The error returned
// may be a go-multierror.
func DownloadFile(ctx context.Context, dst, src, wd string) (string, error) {
	client := getter.Client{
		// following symlinks is a security concern
		DisableSymlinks: true,
	}

	rsp, err := client.Get(ctx, &getter.Request{
		Src:     src,
		Dst:     dst,
		Pwd:     wd,
		Copy:    true, // if we are getting a local file, copy it instead of making a symlink
		GetMode: getter.ModeFile,
	})
	if rsp == nil {
		return "", err
	}

	return rsp.Dst, err
}
