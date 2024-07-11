package varnish

import (
	"bufio"
	"encoding/json"
	"os"
	"sort"
	"strings"
)

type EndpointStats struct {
	TotalRequests int            `json:"total_requests"`
	TopEndpoints  map[string]int `json:"top_endpoints"`
}

func GetVarnishEndpointStats() ([]byte, error) {
	// Log dosyasını oku ve endpoint istatistiklerini parse et
	endpointStats := parseEndpointStats("/var/log/varnish/varnish.log")
	return json.Marshal(endpointStats)
}

func parseEndpointStats(logFile string) EndpointStats {
	stats := EndpointStats{
		TopEndpoints: make(map[string]int),
	}
	lines, err := readLines(logFile)
	if err != nil {
		return stats
	}

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 7 {
			continue
		}
		endpoint := fields[6]
		stats.TotalRequests++
		stats.TopEndpoints[endpoint]++
	}

	stats.TopEndpoints = getTopEndpoints(stats.TopEndpoints, 5)

	return stats
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func getTopEndpoints(endpointMap map[string]int, topN int) map[string]int {
	type kv struct {
		Key   string
		Value int
	}

	var sortedEndpoints []kv
	for k, v := range endpointMap {
		sortedEndpoints = append(sortedEndpoints, kv{k, v})
	}

	sort.Slice(sortedEndpoints, func(i, j int) bool {
		return sortedEndpoints[i].Value > sortedEndpoints[j].Value
	})

	topEndpoints := make(map[string]int)
	for i := 0; i < topN && i < len(sortedEndpoints); i++ {
		topEndpoints[sortedEndpoints[i].Key] = sortedEndpoints[i].Value
	}

	return topEndpoints
}
