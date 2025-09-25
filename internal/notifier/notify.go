package notifier

import (
	"log"
	"os"

	"github.com/wcy-dt/ponghub/internal/types/structures/checker"
	"github.com/wcy-dt/ponghub/internal/types/types/chk_result"
	"github.com/wcy-dt/ponghub/internal/types/types/default_config"
)

// WriteNotifications sends notifications based on the service check results
func WriteNotifications(checkResult []checker.Service, certNotifyDays int) {
	// find all endpoints with status NONE
	statusNoneEndpoints := make(map[string][]string)
	for _, serviceResult := range checkResult {
		for _, endpointResult := range serviceResult.Endpoints {
			if endpointResult.Status == chk_result.NONE {
				statusNoneEndpoints[serviceResult.Name] = append(statusNoneEndpoints[serviceResult.Name], endpointResult.URL)
			}
		}
	}

	// find all endpoints whose certificates are expired or has less than 7 days remaining
	certProblemEndpoints := make(map[string][]string)
	for _, serviceResult := range checkResult {
		for _, endpointResult := range serviceResult.Endpoints {
			if endpointResult.IsHTTPS && (endpointResult.IsCertExpired || endpointResult.CertRemainingDays <= certNotifyDays) {
				certProblemEndpoints[serviceResult.Name] = append(certProblemEndpoints[serviceResult.Name], endpointResult.URL)
			}
		}
	}

	// output to file default_config.GetNotifyPath()
	notifyPath := default_config.GetNotifyPath()
	if err := os.Remove(notifyPath); err != nil && !os.IsNotExist(err) {
		log.Println("Error removing notify file:", err)
		return
	}
	if len(statusNoneEndpoints) == 0 {
		// if no endpointURLs are down, do nothing
		return
	}

	// new notify file
	f, err := os.Create(notifyPath)
	if err != nil {
		log.Println("Error creating notify file:", err)
		return
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Println("Error closing notify file:", err)
		}
	}()

	// write none status endpoints to file
	for serviceName, endpointURLs := range statusNoneEndpoints {
		if _, err := f.WriteString(serviceName + "\n"); err != nil {
			log.Println("Error writing to notify file:", err)
			return
		}
		for _, endpointURL := range endpointURLs {
			if _, err := f.WriteString("\t" + endpointURL + " is unavailable.\n"); err != nil {
				log.Println("Error writing to notify file:", err)
				return
			}
		}
	}

	// write cert problem endpoints to file
	for serviceName, endpointURLs := range certProblemEndpoints {
		if _, err := f.WriteString(serviceName + "\n"); err != nil {
			log.Println("Error writing to notify file:", err)
			return
		}
		for _, endpointURL := range endpointURLs {
			if _, err := f.WriteString("\t" + endpointURL + " has certificate issues.\n"); err != nil {
				log.Println("Error writing to notify file:", err)
				return
			}
		}
	}
}
