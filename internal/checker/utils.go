package checker

import (
	"log"
	"os"
	"strings"

	"github.com/wcy-dt/ponghub/internal/types/types/chk_result"
)

// isTestMode checks if the current execution is in test mode
func isTestMode() bool {
	// Check if any command line arguments contain "test"
	for _, arg := range os.Args {
		if strings.Contains(arg, "test") || strings.Contains(arg, ".test") {
			return true
		}
	}

	// Check if we're running with go test
	if len(os.Args) > 0 && strings.HasSuffix(os.Args[0], ".test") {
		return true
	}

	return false
}

// logIfTest logs the message only if we're in test mode
func logIfTest(format string, args ...interface{}) {
	if isTestMode() {
		log.Printf(format, args...)
	}
}

// getTestResult determines the test result based on the success count and actual attempts
func getTestResult(successNum, attemptNum int) chk_result.CheckResult {
	switch successNum {
	case attemptNum:
		return chk_result.ALL
	case 0:
		return chk_result.NONE
	default:
		return chk_result.PART
	}
}
