package net2

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/user"
	"path/filepath"
	"sync"
	"time"
)

type dnsItem struct {
	UpdatedAt time.Time `json:"u"`
	Hosts     []string  `json:"h"`
}

var (
	tmpfsCacheTimeout = 12 * time.Hour

	dnsMu    sync.RWMutex
	dnscache = make(map[string]dnsItem)
)

var errEmpty = errors.New("returns empty host")

func LookupIPAddr(name string, cachetime time.Duration) ([]string, error) {
	dnsMu.RLock()
	v := dnscache[name]
	dnsMu.RUnlock()
	cachehosts := v.Hosts
	if time.Since(v.UpdatedAt) < cachetime {
		if len(cachehosts) == 0 {
			return nil, errEmpty
		}
		return cachehosts, nil
	}
	dnsMu.Lock()
	defer dnsMu.Unlock()
	if v := dnscache[name]; time.Since(v.UpdatedAt) < cachetime {
		if len(v.Hosts) == 0 {
			return nil, errEmpty
		}
		return v.Hosts, nil
	}
	b, err := ioutil.ReadFile(genIPAddrFilename(name))
	if len(b) > 0 {
		if er := json.Unmarshal(b, &v); er != nil {
			log.Println("[net2] json.Unmarshal", err)
		}
	}
	if err == nil && time.Since(v.UpdatedAt) < tmpfsCacheTimeout && len(v.Hosts) > 0 {
		return v.Hosts, nil
	}
	hosts, err := net.LookupHost(name)
	if err != nil {
		if len(cachehosts) > 0 { // stale
			log.Println("[net2] LookupHost", err)
			return cachehosts, nil
		}
		return nil, err
	}
	if len(hosts) == 0 {
		if len(cachehosts) > 0 {
			log.Println("[net2] LookupHost response empty")
			return cachehosts, nil
		}
		return nil, errEmpty
	}
	item := dnsItem{Hosts: hosts, UpdatedAt: time.Now()}
	if err := save2tempfile(name, item); err != nil {
		log.Println("[net2] save2tempfile", err)
	}
	dnscache[name] = item
	return hosts, nil
}

var cu, _ = user.Current()

func genIPAddrFilename(name string) string {
	if cu != nil {
		return filepath.Join(os.TempDir(), name+"-"+cu.Username+".json")
	}
	return filepath.Join(os.TempDir(), name+".json")
}

func save2tempfile(name string, item dnsItem) error {
	f, err := ioutil.TempFile("", name)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	if err := enc.Encode(item); err != nil {
		return err
	}
	return os.Rename(f.Name(), genIPAddrFilename(name)) // atomic
}
