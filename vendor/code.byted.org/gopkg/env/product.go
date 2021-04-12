package env

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"strings"
)

// IsTesting if in testing env
func IsTesting() bool {
	return os.Getenv("TESTING_PREFIX") == "offline"
}

// IsBoe return true if it's in BOE by the IDC name
func IsBoe() bool {
	idc := IDC()
	return idc == "boe" || idc == "boei18n"
}

// IsBoeCN return true if the idc is in China
func IsBoeCN() bool {
	idc := IDC()
	return idc == "boe"
}

// IsBoeI18N return true if the idc is oversea
func IsBoeI18N() bool {
	idc := IDC()
	return idc == "boei18n"
}

func IsPPE() bool {
	tceHostEnv := os.Getenv("TCE_HOST_ENV")
	return tceHostEnv == "ppe"
}

const processNameFormat = "/proc/%d/comm"

// IsProduct return true if current service is running on product enviroment else false
func IsProduct() bool {
	if IsTesting() || IsBoe() {
		return false
	}

	// please see: https://wiki.bytedance.net/pages/viewpage.action?pageId=63229064
	if os.Getenv("IS_PROD_RUNTIME") != "" {
		return true
	}

	if os.Getenv("SERVICE_ENV") != "" {
		return true
	}

	u, err := user.Current()
	if err != nil {
		return false
	}

	pn, err := parentProcName()
	if err != nil {
		return false
	}
	if u.Username == "tiger" && (pn == "supervise" || pn == "systemd") {
		return true
	}
	return false
}

// father's service name
func parentProcName() (string, error) {
	ppid := os.Getppid()
	bs, err := ioutil.ReadFile(fmt.Sprintf(processNameFormat, ppid))
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(bs)), nil
}
