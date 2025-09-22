package reporter

import (
	"fmt"
	"html/template"
	"log"
	"os"

	"github.com/wcy-dt/ponghub/internal/common"
	"github.com/wcy-dt/ponghub/internal/types/structures/checker"
	"github.com/wcy-dt/ponghub/internal/types/structures/reporter"
	"github.com/wcy-dt/ponghub/internal/types/types/chk_result"
	"github.com/wcy-dt/ponghub/internal/types/types/default_config"
)

// GetReport generates a report based on the check results and log data
func GetReport(checkResult []checker.Checker, logPath string) (reporter.Reporter, error) {
	logResult, err := common.ReadLogs(logPath)
	if err != nil {
		log.Printf("Error loading log data from %s: %v", logPath, err)
		return nil, err
	}

	reportResult := reporter.ParseLogResult(logResult)

	// calculate availability
	for serviceName, serviceReport := range reportResult {
		if len(serviceReport.ServiceHistory) == 0 {
			continue
		}
		statusAllEntryNum := 0
		for _, entry := range serviceReport.ServiceHistory {
			if chk_result.IsALL(entry.Status) {
				statusAllEntryNum++
			}
		}
		availability := float64(statusAllEntryNum) / float64(len(serviceReport.ServiceHistory))
		tmp := reportResult[serviceName]
		tmp.Availability = availability
		reportResult[serviceName] = tmp
	}

	// calculate cert status
	for _, serviceResult := range checkResult {
		serviceName := serviceResult.Name
		for _, endpointResult := range serviceResult.Endpoints {
			url := endpointResult.URL

			tmp := reportResult[serviceName].Endpoints[url]
			tmp.IsHTTPS = endpointResult.IsHTTPS
			tmp.CertRemainingDays = endpointResult.CertRemainingDays
			tmp.IsCertExpired = endpointResult.IsCertExpired
			tmp.DisplayURL = endpointResult.DisplayURL
			tmp.HighlightSegments = endpointResult.HighlightSegments
			reportResult[serviceName].Endpoints[url] = tmp
		}
	}

	return reportResult, nil
}

// WriteReport generates an HTML report from the provided log data
func WriteReport(reportResult reporter.Reporter, reportPath string) error {
	tmpl, err := template.New("report.html").
		Funcs(createTemplateFunc()).
		ParseFiles(default_config.GetTemplatePath())
	if err != nil {
		return fmt.Errorf("template parsing failed: %w", err)
	}

	reportFile, err := os.Create(reportPath)
	if err != nil {
		return fmt.Errorf("file creation failed: %w", err)
	}
	defer func(reportFile *os.File) {
		if err := reportFile.Close(); err != nil {
			fmt.Printf("Error closing output file: %v\n", err)
		}
	}(reportFile)

	if err := tmpl.Execute(reportFile, map[string]any{
		"ReportResult": reportResult,
		"UpdateTime":   getLatestTime(reportResult),
		"DisplayNum":   default_config.GetDisplayNum(),
	}); err != nil {
		return fmt.Errorf("template execution failed: %w", err)
	}

	return nil
}

// getLatestTime retrieves the latest time from the log data
func getLatestTime(reportResult reporter.Reporter) string {
	var latestTime string

	for _, serviceResult := range reportResult {
		for _, serviceHistoryEntry := range serviceResult.ServiceHistory {
			if latestTime == "" {
				latestTime = serviceHistoryEntry.Time
			} else if serviceHistoryEntry.Time > latestTime {
				latestTime = serviceHistoryEntry.Time
			}
		}
	}

	return latestTime
}

// createTemplateFunc defines custom template functions for the report
func createTemplateFunc() template.FuncMap {
	return template.FuncMap{
		"add": func(a, b int) int { return a + b },
		"sub": func(a, b int) int { return a - b },
		"mul": func(a, b float64) float64 { return a * b },
		"div": func(a, b int) float64 {
			if b == 0 {
				return 0
			}
			return float64(a) / float64(b)
		},
		"until": func(n int) []int {
			result := make([]int, n)
			for i := range n {
				result[i] = i
			}
			return result
		},
	}
}
