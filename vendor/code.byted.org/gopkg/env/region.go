package env

import (
	"encoding/json"
	"io/ioutil"
	"sync/atomic"
	"time"
)

const (
	// UnknownRegion .
	UnknownRegion = "-"
	R_CN          = "CN"
	R_SG          = "SG"
	R_US          = "US"
	R_MALIVA      = "MALIVA"
	R_ALISG       = "ALISG" // Singapore Aliyun
	R_CA          = "CA"    // West America
	R_BOE         = "BOE"
	R_SUITEKA     = "SUITEKA"
	R_SUITEKA2    = "SUITEKA2"
	R_SUITEKA3    = "SUITEKA3"
)

const regionFileDefault = "/opt/tiger/chadc/region.json"
const regionRefreshDur = 5 * time.Minute

var (
	region     atomic.Value // local region (string)
	regionFile atomic.Value // file name (string)
	idcRegion  atomic.Value // key is IDC, value is region (map[string]string)
)

// Region .
func Region() string {
	if v := region.Load(); v != nil {
		return v.(string)
	}

	idc := IDC()
	regionResult := GetRegionFromIDC(idc)
	region.Store(regionResult)
	return regionResult
}

// GetRegionFromIDC .
func GetRegionFromIDC(idc string) string {
	if m, ok := idcRegion.Load().(map[string]string); ok {
		if r, ok := m[idc]; ok {
			return r
		}
	}
	return UnknownRegion
}

// SetRegionFile .
func SetRegionFile(file string) {
	regionFile.Store(file)
	updateRegionIDCs()
}

// hasIDC return true if idc is support
func hasIDC(idc string) bool {
	m, _ := idcRegion.Load().(map[string]string)
	_, ok := m[idc]
	return ok
}

// idcList return IDCs which supported
func idcList() []string {
	m, _ := idcRegion.Load().(map[string]string)
	idcList := make([]string, 0, len(m))
	for dc, _ := range m {
		idcList = append(idcList, dc)
	}
	return idcList
}

func init() {
	regionFile.Store(regionFileDefault)
	refreshRegion()
}

// ATTENTION: IT COMES WITH A LOOP, DON'T CALL IT AGAIN.
func refreshRegion() {
	defer time.AfterFunc(regionRefreshDur, refreshRegion)
	updateRegionIDCs()
}

func updateRegionIDCs() {
	// newRegionIDCs, key is region, value is []IDC
	newRegionIDCs := make(map[string][]string)
	file, _ := regionFile.Load().(string)
	content, err := ioutil.ReadFile(file)
	if err == nil {
		err = json.Unmarshal(content, &newRegionIDCs)
	}
	// if error, skip update
	if err != nil {
		return
	}
	newIDCRegion := make(map[string]string)
	for region, dcs := range newRegionIDCs {
		for _, dc := range dcs {
			newIDCRegion[dc] = region
		}
	}
	idcRegion.Store(newIDCRegion)
}
