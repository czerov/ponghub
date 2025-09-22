package reporter

import "github.com/wcy-dt/ponghub/internal/types/types/highlight"

// Data structures for logging and reporting
type (
	// HistoryEntry represents a single history entry
	HistoryEntry struct {
		Time         string
		Status       string
		ResponseTime int
	}

	History []HistoryEntry

	Endpoint struct {
		EndpointHistory   History
		IsHTTPS           bool
		IsCertExpired     bool
		CertRemainingDays int
		DisplayURL        string              // Resolved URL for display
		HighlightSegments []highlight.Segment // Segments with highlight info
	}

	Endpoints map[string]Endpoint

	// Service represents the result of checking a service
	Service struct {
		ServiceHistory History
		Availability   float64
		Endpoints      Endpoints
	}

	// Reporter represents the result of checking services
	Reporter map[string]Service
)
