package img

import (
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/zikster3262/shared-lib/utils"
)

type Image struct {
	Title    string `json:"title"`
	URL      string `json:"url"`
	Chapter  string `json:"chapter"`
	Filename string `json:"filename"`
}

func (i Image) DownloadFile() ([]byte, error) {
	ctx := context.Background()

	var body io.Reader

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, i.URL, body)
	if err != nil {
		return nil, errors.Unwrap(err)
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, errors.Unwrap(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)

		if err != nil {
			utils.FailOnCmpError("shared", "img-res-code", err)
		}

		return bodyBytes, nil
	}

	return nil, errors.Unwrap(err)
}
