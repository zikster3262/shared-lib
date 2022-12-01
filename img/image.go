package img

import (
	"io"
	"net/http"

	"github.com/zikster3262/shared-lib/utils"
)

type Image struct {
	Title    string `json:"title"`
	Url      string `json:"url"`
	Chapter  string `json:"chapter"`
	Filename string `json:"filename"`
}

func (i Image) DownloadFile() []byte {

	var client http.Client
	resp, err := client.Get(i.Url)
	if err != nil {
		utils.FailOnCmpError("shared", "img-get", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			utils.FailOnCmpError("shared", "img-res-code", err)
		}

		return bodyBytes
	}
	return nil
}
