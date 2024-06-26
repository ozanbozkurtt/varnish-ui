package varnish

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

type VarnishStats struct {
	ClientReq     int `json:"client_req"`
	CacheHit      int `json:"cache_hit"`
	CacheHitGrace int `json:"cache_hit_grace"`
	CacheMiss     int `json:"cache_miss"`
	Uncacheable   int `json:"uncacheable"`
}

func GetVarnishStats() ([]byte, error) {
	cmd := exec.Command("varnishstat", "-1")
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("error running varnishstat: %v", err)
	}

	lines := strings.Split(out.String(), "\n")
	stats := &VarnishStats{}

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		value := fields[1]
		switch fields[0] {
		case "MAIN.client_req":
			fmt.Sscanf(value, "%d", &stats.ClientReq)
		case "MAIN.cache_hit":
			fmt.Sscanf(value, "%d", &stats.CacheHit)
		case "MAIN.cache_hit_grace":
			fmt.Sscanf(value, "%d", &stats.CacheHitGrace)
		case "MAIN.cache_miss":
			fmt.Sscanf(value, "%d", &stats.CacheMiss)
		case "MAIN.beresp_uncacheable":
			fmt.Sscanf(value, "%d", &stats.Uncacheable)
		}
	}

	return json.Marshal(stats)
}
