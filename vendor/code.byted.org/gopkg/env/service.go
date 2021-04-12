package env

import "os"

const (
	PSMUnknown     = "-"
	ClusterDefault = "default"
)

var psm string
var cluster string

func init() {
	psm = os.Getenv("LOAD_SERVICE_PSM")
	if psm == "" {
		psm = os.Getenv("PSM")
	}
	if psm == "" {
		psm = os.Getenv("TCE_PSM")
	}
	if psm == "" {
		psm = PSMUnknown
	}

	cluster = os.Getenv("SERVICE_CLUSTER")
	if cluster == "" {
		cluster = ClusterDefault
	}
}

// PSM .
func PSM() string {
	return psm
}

// SetPSM is used for unit test.
func SetPSM(psm_ string) {
	psm = psm_
}

// Cluster .
func Cluster() string {
	return cluster
}
