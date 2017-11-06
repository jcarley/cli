package logs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

func chooseQueryGenerator(version string) queryGenerator {
	generator := generateES5Query
	if strings.HasPrefix(version, "1.") {
		generator = generateES1Query
	} else if strings.HasPrefix(version, "2.") {
		generator = generateES2Query
	}
	return generator
}

func generateES5Query(queryString, appLogsIdentifier, appLogsValue string, timestamp time.Time, from int, hostNames []string, fileName string) ([]byte, error) {
	hostFilter, fileFilter := createFilters(hostNames, fileName)
	query := `{
	"_source": ["@timestamp", "message", "` + appLogsIdentifier + `"],
	"query": {
		"bool": {
			"must": [
				{"wildcard": {"message": "` + queryString + `"}},
				{"term": {"` + appLogsIdentifier + `": "` + appLogsValue + `"}},` + fileFilter + `
				{"range": {"@timestamp": {"gt": "` + fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02dZ", timestamp.Year(), timestamp.Month(), timestamp.Day(), timestamp.Hour(), timestamp.Minute(), timestamp.Second()) + `"}}}
			]` + hostFilter + `
		}
	},
	"sort": [
		{
			"@timestamp": {
				"order": "asc",
				"unmapped_type":"boolean"
			}
		},
		{
			"message.raw": {
				"order": "asc",
				"unmapped_type":"boolean"
			}
		}
	],
	"from": ` + fmt.Sprintf("%d", from) + `,
	"size": ` + fmt.Sprintf("%d", size) + `
	}`
	var buf bytes.Buffer
	err := json.Compact(&buf, []byte(query))
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func generateES2Query(queryString, appLogsIdentifier, appLogsValue string, timestamp time.Time, from int, hostNames []string, fileName string) ([]byte, error) {
	hostFilter, fileFilter := createFilters(hostNames, fileName)
	query := `{
	"fields": ["@timestamp", "message", "` + appLogsIdentifier + `"],
	"query": {
		"wildcard": {
			"message": "` + queryString + `"
		}
	},
	"filter": {
		"query": {
			"bool": {
				"must": [
					{"term": {"` + appLogsIdentifier + `": "` + appLogsValue + `"}},
					` + fileFilter + `
					{"range": {"@timestamp": {"gt": "` + fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02dZ", timestamp.Year(), timestamp.Month(), timestamp.Day(), timestamp.Hour(), timestamp.Minute(), timestamp.Second()) + `"}}}
				]` + hostFilter + `
			}
		}
	},
	"sort": [
		{
			"@timestamp": {
				"order": "asc",
				"unmapped_type":"boolean"
			}
		},
		{
			"message.raw": {
				"order": "asc",
				"unmapped_type":"boolean"
			}
		}
	],
	"from": ` + fmt.Sprintf("%d", from) + `,
	"size": ` + fmt.Sprintf("%d", size) + `
	}`
	var buf bytes.Buffer
	err := json.Compact(&buf, []byte(query))
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func generateES1Query(queryString, appLogsIdentifier, appLogsValue string, timestamp time.Time, from int, hostNames []string, fileName string) ([]byte, error) {
	hostFilter, fileFilter := createFilters(hostNames, fileName)
	query := `{
	"fields": ["@timestamp", "message", "` + appLogsIdentifier + `"],
	"query": {
		"wildcard": {
			"message": "` + queryString + `"
		}
	},
	"filter": {
		"query": {
			"bool": {
				"must": [
					{"term": {"` + appLogsIdentifier + `": "` + appLogsValue + `"}},
					` + fileFilter + `
					{"range": {"@timestamp": {"gt": "` + fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02dZ", timestamp.Year(), timestamp.Month(), timestamp.Day(), timestamp.Hour(), timestamp.Minute(), timestamp.Second()) + `"}}}
				]` + hostFilter + `
			}
		}
	},
	"sort": {
		"@timestamp": {
			"order": "asc"
		},
		"message": {
			"order": "asc"
		}
	},
	"from": ` + fmt.Sprintf("%d", from) + `,
	"size": ` + fmt.Sprintf("%d", size) + `
	}`
	var buf bytes.Buffer
	json.Compact(&buf, []byte(query))
	return buf.Bytes(), nil
}

func createFilters(hostNames []string, fileName string) (string, string) {
	var hostFilter string
	var fileFilter string
	if len(hostNames) > 0 {
		formattedHostNames := make([]string, len(hostNames))
		for i, hostName := range hostNames {
			formattedHostNames[i] = fmt.Sprintf(`"%s"`, hostName)
		}
		hostFilter = `,
			"should": [`
		for _, hostName := range hostNames {
			hostFilter += `
				{"match_phrase": {"host": "` + hostName + `"}},`
		}
		hostFilter = strings.TrimSuffix(hostFilter, ",")
		hostFilter += `
			],
			"minimum_should_match": 1`
	} else if len(fileName) > 0 {
		fileFilter += `
		{"match_phrase": {"file": "` + fileName + `"}},`
	}
	return hostFilter, fileFilter
}
