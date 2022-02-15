package hub

import (
	"context"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"strconv"

	"github.com/pkg/errors"
)

type BackupFile struct {
	Name string
	Size int
}

func (c *Client) DownloadBackupFile(ctx context.Context, name string) (BackupFile, io.ReadCloser, error) {
	c.log.V(2).Info("requesting backup file", "name", name)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("/hub/backupDB?fileName=%s", url.QueryEscape(name)), nil)
	if err != nil {
		return BackupFile{}, nil, errors.Wrap(err, "creating request")
	}

	res, err := c.Do(req)
	if err != nil {
		return BackupFile{}, nil, errors.Wrap(err, "sending request")
	} else if res.StatusCode != http.StatusOK {
		if cerr := res.Body.Close(); cerr != nil {
			c.log.Error(cerr, "closing errored response")
		}

		return BackupFile{}, nil, fmt.Errorf("unexpected response status code: %d", res.StatusCode)
	}

	var meta BackupFile

	if contentDisposition := res.Header.Get("content-disposition"); len(contentDisposition) > 0 {
		_, contentDispositionParams, err := mime.ParseMediaType(contentDisposition)
		if err != nil {
			return BackupFile{}, nil, errors.Wrap(err, "parsing content-disposition")
		}

		meta.Name = contentDispositionParams["filename"]
	}

	if contentLength := res.Header.Get("content-length"); len(contentLength) > 0 {
		contentLength, err := strconv.Atoi(contentLength)
		if err != nil {
			return BackupFile{}, nil, errors.Wrap(err, "parsing content-length")
		}

		meta.Size = contentLength
	}

	c.log.V(1).Info("requested backup file", "name", name)

	return meta, res.Body, nil
}
