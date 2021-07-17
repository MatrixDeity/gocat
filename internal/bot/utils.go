package bot

import (
	"fmt"
	"os"
	"strings"
)

func getEnvVarChecked(varname string) string {
	value := os.Getenv(varname)
	if len(strings.TrimSpace(value)) == 0 {
		panic(fmt.Errorf("%s env var is empty", varname))
	}
	return value
}
