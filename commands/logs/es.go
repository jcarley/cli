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
	additionalFields, hostFilter, fileFilter := createFeildsAndFilters(hostNames, fileName, 2)
	query := `{
	"stored_fields": ["@timestamp", "message", "` + appLogsIdentifier + `"` + additionalFields + `],
	"query": {
		"wildcard": {
			"message": "` + queryString + `"
		}
	},
	"post_filter": {
		"bool": {
			"must": [
				{"term": {"` + appLogsIdentifier + `": "` + appLogsValue + `"}},
				{"range": {"@timestamp": {"gt": "` + fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02dZ", timestamp.Year(), timestamp.Month(), timestamp.Day(), timestamp.Hour(), timestamp.Minute(), timestamp.Second()) + `"}}}
				` + hostFilter + fileFilter + `
			]
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
	additionalFields, hostFilter, fileFilter := createFeildsAndFilters(hostNames, fileName, 2)
	query := `{
	"fields": ["@timestamp", "message", "` + appLogsIdentifier + `"` + additionalFields + `],
	"query": {
		"wildcard": {
			"message": "` + queryString + `"
		}
	},
	"filter": {
		"bool": {
			"must": [
				{"term": {"` + appLogsIdentifier + `": "` + appLogsValue + `"}},
				{"range": {"@timestamp": {"gt": "` + fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02dZ", timestamp.Year(), timestamp.Month(), timestamp.Day(), timestamp.Hour(), timestamp.Minute(), timestamp.Second()) + `"}}}
				` + hostFilter + fileFilter + `
			]
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
	additionalFields, hostFilter, fileFilter := createFeildsAndFilters(hostNames, fileName, 1)
	query := `{
	"fields": ["@timestamp", "message", "` + appLogsIdentifier + `"` + additionalFields + `],
	"query": {
		"wildcard": {
			"message": "` + queryString + `"
		}
	},
	"filter": {
		"bool": {
			"must": [
				{"term": {"` + appLogsIdentifier + `": "` + appLogsValue + `"}},
				{"range": {"@timestamp": {"gt": "` + fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02dZ", timestamp.Year(), timestamp.Month(), timestamp.Day(), timestamp.Hour(), timestamp.Minute(), timestamp.Second()) + `"}}}
				` + fileFilter + `
			]` + hostFilter + `
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

func createFeildsAndFilters(hostNames []string, fileName string, version int) (string, string, string) {
	var hostFilter string
	var fileFilter string
	var additionalFields string
	if len(hostNames) > 0 {
		if version >= 2 {
			additionalFields += `, "host"`
			for i, hostName := range hostNames {
				hostNames[i] = fmt.Sprintf(`"%s"`, hostName)
			}
			hostFilter = fmt.Sprintf(`, {"terms": {"host": [%s]}}`, strings.Join(hostNames, ", "))
		} else {
			additionalFields += `, "host"`
			hostFilter = `,
			"should": [`
			for _, hostName := range hostNames {
				hostFilter += `
				{"term": {"host": "` + hostName + `"}},`
			}
			hostFilter = strings.TrimSuffix(hostFilter, ",")
			hostFilter += `
			]`
		}
	} else if len(fileName) > 0 {
		additionalFields += `, "file"`
		fileFilter += `, {"term": {"file": "` + fileName + `"}}`
	}
	return additionalFields, hostFilter, fileFilter
}
