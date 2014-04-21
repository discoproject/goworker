package jobutil

import (
	"os"
)

var localDict map[string]string

//TODO actually query the settings
func Setting(str string) string {
	if val, ok := localDict[str]; ok {
		return val
	}
	return os.Getenv(str)
}

func SetKeyValue(key string, value string) {
	if localDict == nil {
		localDict = make(map[string]string)
	}
	localDict[key] = value
}
