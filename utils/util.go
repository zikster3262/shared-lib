package utils

import (
	"encoding/json"
	"regexp"
	"strings"

	"github.com/rs/zerolog/log"
)

func StructToJson(data interface{}) []byte {
	var bt []byte
	var err error
	bt, err = json.Marshal(data)
	if err != nil {
		log.Error().Msg("couldn't unmarshall the request")
	}

	return bt
}

func GetIDFromChapterURL(url string) string {

	re := regexp.MustCompile("[0-9]+")
	res := re.FindAllString(url, -1)

	if len(res) > 1 {
		return res[len(res)-1]
	} else {
		return res[0]
	}

}

func GetFileName(url string) string {
	result := strings.LastIndex(url, "/")
	return url[result+1:]
}
