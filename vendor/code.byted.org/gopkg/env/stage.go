package env

import "os"

const (
	UnknownStage  = "-"
	StageStaging  = "staging"
	StageCanary   = "canary"
	StageSingleDC = "single_dc"
	StageAllDC    = "all_dc"
)

var stage = UnknownStage

func init() {
	if os.Getenv("IS_TCE_DOCKER_ENV") == "1" {
		tceStage := os.Getenv("TCE_STAGE")
		if tceStage != "" {
			stage = tceStage
		}
	}
}

// Stage .
func Stage() string {
	return stage
}
