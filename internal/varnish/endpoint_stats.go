package varnish

import (
	"bufio"
	"encoding/json"
	"os/exec"
	"sort"
	"strings"
)

type EndpointStats struct {
	TotalRequests int            `json:"total_requests"`
	TopEndpoints  map[string]int `json:"top_endpoints"`
}

func GetEndpointStats() ([]byte, error) {
	cmd := exec.Command("varnishncsa", "-w", "/dev/stdout")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	endpoints := make(map[string]int)
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) > 6 {
			endpoint := fields[6]
			endpoints[endpoint]++
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	totalRequests := 0
	for _, count := range endpoints {
		totalRequests += count
	}

	type kv struct {
		Key   string
		Value int
	}

	var ss []kv
	for k, v := range endpoints {
		ss = append(ss, kv{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})

	topEndpoints := make(map[string]int)
	for i, kv := range ss {
		if i >= 5 {
			break
		}
		topEndpoints[kv.Key] = kv.Value
	}

	stats := EndpointStats{
		TotalRequests: totalRequests,
		TopEndpoints:  topEndpoints,
	}

	return json.Marshal(stats)
}
