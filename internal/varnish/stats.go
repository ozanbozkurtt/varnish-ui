package varnish

import (
	"encoding/json"
	"os/exec"
	"strconv"
	"strings"
)

type VarnishStats struct {
	ClientReq      int `json:"client_req"`
	CacheHit       int `json:"cache_hit"`
	CacheHitGrace  int `json:"cache_hit_grace"`
	CacheMiss      int `json:"cache_miss"`
	Uncacheable    int `json:"uncacheable"`
	BackendConn    int `json:"backend_conn"`
	BackendReuse   int `json:"backend_reuse"`
	BackendRecycle int `json:"backend_recycle"`
	FetchLength    int `json:"fetch_length"`
	FetchChunked   int `json:"fetch_chunked"`
	SessConn       int `json:"sess_conn"`
	NObject        int `json:"n_object"`
	Expired        int `json:"expired"`
	Threads        int `json:"threads"`
	Bans           int `json:"bans"`
}

func GetVarnishStats() ([]byte, error) {
	cmd := exec.Command("varnishstat", "-1")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	stats := parseVarnishStats(string(output))
	return json.Marshal(stats)
}

func parseVarnishStats(output string) VarnishStats {
	stats := VarnishStats{}
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}
		value := fields[1]

		switch fields[0] {
		case "MAIN.client_req":
			stats.ClientReq = parseValue(value)
		case "MAIN.cache_hit":
			stats.CacheHit = parseValue(value)
		case "MAIN.cache_hit_grace":
			stats.CacheHitGrace = parseValue(value)
		case "MAIN.cache_miss":
			stats.CacheMiss = parseValue(value)
		case "MAIN.beresp_uncacheable":
			stats.Uncacheable = parseValue(value)
		case "MAIN.backend_conn":
			stats.BackendConn = parseValue(value)
		case "MAIN.backend_reuse":
			stats.BackendReuse = parseValue(value)
		case "MAIN.backend_recycle":
			stats.BackendRecycle = parseValue(value)
		case "MAIN.fetch_length":
			stats.FetchLength = parseValue(value)
		case "MAIN.fetch_chunked":
			stats.FetchChunked = parseValue(value)
		case "MAIN.sess_conn":
			stats.SessConn = parseValue(value)
		case "MAIN.n_object":
			stats.NObject = parseValue(value)
		case "MAIN.n_expired":
			stats.Expired = parseValue(value)
		case "MAIN.threads":
			stats.Threads = parseValue(value)
		case "MAIN.bans":
			stats.Bans = parseValue(value)
		}
	}

	return stats
}

func parseValue(value string) int {
	parsedValue, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	return parsedValue
}
