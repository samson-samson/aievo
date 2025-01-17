package json

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/antgroup/aievo/utils/json/gjson"
	"github.com/mitchellh/mapstructure"
)

const JsonParse = "(?s)```json\n(.*?)\n```"

func Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

func Unmarshal(data []byte, v any) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprint(r))
		}
	}()
	result := gjson.Parse(string(data))
	if result.Type == gjson.Null {
		return errors.New("invalid json")
	}
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Metadata: nil,
		Result:   v,
		TagName:  "json",
	})
	if err != nil {
		return err
	}
	err = decoder.Decode(result.Value())
	return err
}

func TrimJsonString(input string) string {
	input = strings.TrimSpace(input)
	compile := regexp.MustCompile(JsonParse)
	submatch := compile.FindAllStringSubmatch(input, -1)
	if len(submatch) != 0 {
		input = submatch[0][1]
	}
	return input
}
