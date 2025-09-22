package reporter

import (
	"log"
	"sort"

	"github.com/wcy-dt/ponghub/internal/types/structures/logger"
)

// convertToHistory converts logger history entries to reporter history format,
// sorts them by time and returns only the last displayNum entries.
func convertToHistory(logEntries []logger.HistoryEntry, displayNum int) History {
	var history History
	for _, entry := range logEntries {
		history = append(history, HistoryEntry{
			Time:         entry.Time,
			Status:       entry.Status,
			ResponseTime: entry.ResponseTime,
		})
	}

	// Sort the history by time, the most recent entries first
	sort.Slice(history, func(i, j int) bool {
		return history[i].Time < history[j].Time
	})

	// Get only the last `displayNum` entries
	if len(history) > displayNum {
		history = history[len(history)-displayNum:]
	}

	return history
}

// ParseLogResult converts logger.Logger data into a reporter.Reporter format.
func ParseLogResult(logResult logger.Logger, displayNum int) Reporter {
	report := make(Reporter)
	for serviceName, serviceLog := range logResult {
		if len(serviceLog.ServiceHistory) == 0 {
			log.Printf("No history data for service %s", serviceName)
			continue // Skip services with no history data
		}

		// Convert logger.Endpoints to reporter.Endpoints
		endpoints := make(Endpoints)
		for url, endpointLog := range serviceLog.Endpoints {
			endpointHistory := convertToHistory(endpointLog, displayNum)
			endpoints[url] = Endpoint{
				EndpointHistory: endpointHistory,
			}
		}

		// convert logger.ServiceHistory to reporter.ServiceHistory
		serviceHistory := convertToHistory(serviceLog.ServiceHistory, displayNum)

		newService := Service{
			ServiceHistory: serviceHistory,
			Endpoints:      endpoints,
		}
		report[serviceName] = newService
	}
	return report
}
