package env

import (
	"os"
	"regexp"
	"strings"
)

var env string

func init() {
	env = os.Getenv("SERVICE_ENV")
	if env == "" {
		env = os.Getenv("TCE_ENV")
	}
	if env == "" {
		env = "prod"
	}
}

// https://bytedance.feishu.cn/docs/ws0Kg2Q539Rzqc18VFc98f
func Env() string {
	return env
}

// https://bytedance.feishu.cn/docs/doccn2zIJMueBhHi4kjrYKqaVYc#BTZPX3
var (
	isValidEnv = regexp.MustCompile("^[_a-zA-Z0-9][-_.a-zA-Z0-9]*$").MatchString
)

func IsValidEnv(env string) bool {
	if len(env) > 128 {
		return false
	}
	env = strings.TrimSpace(env)
	return len(env) == 0 || isValidEnv(env)
}
