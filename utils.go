package sophoscentral

import "encoding/json"

func PrettyPrint(i interface{}) string {
	b, err := json.MarshalIndent(i, "", "  ")
	if err != nil {
		return err.Error()
	}
	return string(b)
}
