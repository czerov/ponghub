package common

import (
	"crypto/rand"
	"fmt"
	"math/big"
	mathrand "math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/wcy-dt/ponghub/internal/types/types/highlight"
)

// ParameterResolver handles dynamic parameter resolution in configuration
type ParameterResolver struct {
	currentTime time.Time
	randSource  *mathrand.Rand
}

// NewParameterResolver creates a new parameter resolver with current time
func NewParameterResolver() *ParameterResolver {
	return &ParameterResolver{
		currentTime: time.Now(),
		randSource:  mathrand.New(mathrand.NewSource(time.Now().UnixNano())),
	}
}

// NewParameterResolverWithTime creates a new parameter resolver with specified time
func NewParameterResolverWithTime(t time.Time) *ParameterResolver {
	return &ParameterResolver{
		currentTime: t,
		randSource:  mathrand.New(mathrand.NewSource(t.UnixNano())),
	}
}

// generateRandomString generates a random string of specified length
func (pr *ParameterResolver) generateRandomString(length int, charset string) string {
	if charset == "" {
		charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	}

	result := make([]byte, length)
	for i := range result {
		result[i] = charset[pr.randSource.Intn(len(charset))]
	}
	return string(result)
}

// generateSecureRandomString generates a cryptographically secure random string
func (pr *ParameterResolver) generateSecureRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		result[i] = charset[num.Int64()]
	}
	return string(result)
}

// formatTimeWithPattern converts Go time format patterns to actual values
func (pr *ParameterResolver) formatTimeWithPattern(pattern string) string {
	// Convert strftime-like patterns to Go time format
	replacements := map[string]string{
		"%Y": "2006",    // 4-digit year
		"%y": "06",      // 2-digit year
		"%m": "01",      // month (01-12)
		"%d": "02",      // day (01-31)
		"%H": "15",      // hour (00-23)
		"%M": "04",      // minute (00-59)
		"%S": "05",      // second (00-59)
		"%B": "January", // full month name
		"%b": "Jan",     // abbreviated month name
		"%A": "Monday",  // full weekday name
		"%a": "Mon",     // abbreviated weekday name
		"%j": "002",     // day of year (001-366)
		"%U": "",        // week of year (placeholder)
		"%W": "",        // week of year (placeholder)
		"%w": "",        // weekday (placeholder)
		"%Z": "MST",     // timezone name
		"%z": "-0700",   // timezone offset
		"%s": "",        // Unix timestamp (placeholder)
	}

	goFormat := pattern
	for strftime, goFmt := range replacements {
		if goFmt != "" {
			goFormat = strings.ReplaceAll(goFormat, strftime, goFmt)
		}
	}

	// Handle special cases that don't have direct Go equivalents
	if strings.Contains(pattern, "%U") || strings.Contains(pattern, "%W") {
		_, week := pr.currentTime.ISOWeek()
		weekStr := fmt.Sprintf("%02d", week)
		goFormat = strings.ReplaceAll(goFormat, "%U", weekStr)
		goFormat = strings.ReplaceAll(goFormat, "%W", weekStr)
	}

	if strings.Contains(pattern, "%w") {
		weekday := int(pr.currentTime.Weekday())
		goFormat = strings.ReplaceAll(goFormat, "%w", fmt.Sprintf("%d", weekday))
	}

	if strings.Contains(pattern, "%s") {
		timestamp := fmt.Sprintf("%d", pr.currentTime.Unix())
		goFormat = strings.ReplaceAll(goFormat, "%s", timestamp)
	}

	return pr.currentTime.Format(goFormat)
}

// resolveSpecialParameter resolves non-datetime special parameters
func (pr *ParameterResolver) resolveSpecialParameter(param string) string {
	// Handle different types of special parameters
	switch {
	// UUID generation
	case param == "uuid":
		return uuid.New().String()
	case param == "uuid_short":
		return strings.ReplaceAll(uuid.New().String(), "-", "")

	// Random numbers
	case param == "rand":
		return fmt.Sprintf("%d", pr.randSource.Intn(1000000))
	case param == "rand_int":
		return fmt.Sprintf("%d", pr.randSource.Intn(2147483647))
	case strings.HasPrefix(param, "rand(") && strings.HasSuffix(param, ")"):
		// Handle rand(min,max) format
		content := param[5 : len(param)-1]
		parts := strings.Split(content, ",")
		if len(parts) == 2 {
			paramMin, err1 := strconv.Atoi(strings.TrimSpace(parts[0]))
			paramMax, err2 := strconv.Atoi(strings.TrimSpace(parts[1]))
			if err1 == nil && err2 == nil && paramMax > paramMin {
				return fmt.Sprintf("%d", pr.randSource.Intn(paramMax-paramMin)+paramMin)
			}
		}
		return fmt.Sprintf("%d", pr.randSource.Intn(1000000))

	// Random strings
	case param == "rand_str":
		return pr.generateRandomString(8, "")
	case param == "rand_str_secure":
		return pr.generateSecureRandomString(16)
	case strings.HasPrefix(param, "rand_str(") && strings.HasSuffix(param, ")"):
		// Handle rand_str(length) format
		lengthStr := param[9 : len(param)-1]
		if length, err := strconv.Atoi(strings.TrimSpace(lengthStr)); err == nil && length > 0 {
			return pr.generateRandomString(length, "")
		}
		return pr.generateRandomString(8, "")
	case strings.HasPrefix(param, "rand_hex(") && strings.HasSuffix(param, ")"):
		// Handle rand_hex(length) format
		lengthStr := param[9 : len(param)-1]
		if length, err := strconv.Atoi(strings.TrimSpace(lengthStr)); err == nil && length > 0 {
			return pr.generateRandomString(length, "0123456789abcdef")
		}
		return pr.generateRandomString(8, "0123456789abcdef")

	// Environment variables
	case strings.HasPrefix(param, "env(") && strings.HasSuffix(param, ")"):
		envVar := param[4 : len(param)-1]
		if value := os.Getenv(envVar); value != "" {
			return value
		}
		return ""

	// Sequence numbers (based on current time)
	case param == "seq":
		return fmt.Sprintf("%d", pr.currentTime.UnixNano()%1000000)
	case param == "seq_daily":
		// Daily sequence: seconds since midnight
		midnight := time.Date(pr.currentTime.Year(), pr.currentTime.Month(), pr.currentTime.Day(), 0, 0, 0, 0, pr.currentTime.Location())
		return fmt.Sprintf("%d", int(pr.currentTime.Sub(midnight).Seconds()))

	// Hash-like values
	case param == "hash_short":
		return fmt.Sprintf("%x", pr.currentTime.UnixNano()%0xFFFFFF)
	case param == "hash_md5_like":
		return fmt.Sprintf("%032x", pr.currentTime.UnixNano())

	default:
		// If it's a time format, try to format it
		if strings.Contains(param, "%") {
			return pr.formatTimeWithPattern(param)
		}
		return param // Return as-is if not recognized
	}
}

// resolveSpecialParameterForDisplay resolves parameters for display with sensitive data masking
func (pr *ParameterResolver) resolveSpecialParameterForDisplay(param string) string {
	// Handle different types of special parameters
	switch {
	// Environment variables - mask sensitive values
	case strings.HasPrefix(param, "env(") && strings.HasSuffix(param, ")"):
		envVar := param[4 : len(param)-1]
		if value := os.Getenv(envVar); value != "" {
			return pr.maskSensitiveValue(value, envVar)
		}
		return ""

	// For other parameters, use normal resolution
	default:
		return pr.resolveSpecialParameter(param)
	}
}

// maskSensitiveValue masks sensitive environment variable values
func (pr *ParameterResolver) maskSensitiveValue(value, envVar string) string {
	// List of sensitive environment variable patterns
	sensitivePatterns := []string{
		"key", "secret", "token", "password", "pass", "pwd",
		"auth", "credential", "private", "api_key", "access",
		"jwt", "bearer", "signature", "hash", "salt",
	}

	envVarLower := strings.ToLower(envVar)

	// Check if this environment variable name suggests it contains sensitive data
	for _, pattern := range sensitivePatterns {
		if strings.Contains(envVarLower, pattern) {
			return pr.maskValue(value)
		}
	}

	// If value looks like a token/key (long alphanumeric string), mask it
	if len(value) > 20 && regexp.MustCompile(`^[a-zA-Z0-9+/=-]+$`).MatchString(value) {
		return pr.maskValue(value)
	}

	// Return original value if not considered sensitive
	return value
}

// maskValue creates a masked version of a sensitive value
func (pr *ParameterResolver) maskValue(value string) string {
	if len(value) == 0 {
		return value
	}

	if len(value) <= 4 {
		return strings.Repeat("*", len(value))
	}

	// Show first 2 and last 2 characters, mask the middle
	visible := 2
	if len(value) < 8 {
		visible = 1
	}

	prefix := value[:visible]
	suffix := value[len(value)-visible:]
	maskLength := len(value) - 2*visible

	return prefix + strings.Repeat("*", maskLength) + suffix
}

// ResolveParameters resolves dynamic parameters in a string
func (pr *ParameterResolver) ResolveParameters(input string) string {
	// Use regex to find and replace parameters in {{...}} format
	re := regexp.MustCompile(`\{\{([^}]+)}}`)

	result := re.ReplaceAllStringFunc(input, func(match string) string {
		// Extract the parameter from {{parameter}}
		param := strings.TrimSpace(re.FindStringSubmatch(match)[1])

		// Handle time format parameters (starting with %)
		if strings.HasPrefix(param, "%") {
			return pr.formatTimeWithPattern(param)
		}

		// Handle other special parameters
		return pr.resolveSpecialParameter(param)
	})

	return result
}

// GetResolvedValue returns the resolved value for display purposes
func (pr *ParameterResolver) GetResolvedValue(original string) string {
	return pr.ResolveParameters(original)
}

// GetOriginalValue returns the original template value for configuration display
func (pr *ParameterResolver) GetOriginalValue(original string) string {
	return original
}

// HighlightChanges creates a highlighted version showing what parts were replaced
func (pr *ParameterResolver) HighlightChanges(originalURL string) (string, []highlight.Segment) {
	re := regexp.MustCompile(`\{\{([^}]+)}}`)
	matches := re.FindAllStringSubmatchIndex(originalURL, -1)

	if len(matches) == 0 {
		return originalURL, nil
	}

	var segments []highlight.Segment
	result := ""
	lastEnd := 0

	for _, match := range matches {
		start, end := match[0], match[1]

		// Add unchanged part before the match
		if start > lastEnd {
			unchanged := originalURL[lastEnd:start]
			result += unchanged
			segments = append(segments, highlight.Segment{
				Text:        unchanged,
				IsHighlight: false,
			})
		}

		// Replace the match with resolved value (with masking for sensitive data)
		param := strings.TrimSpace(originalURL[match[2]:match[3]]) // Extract parameter without {{...}}
		var resolved string

		if strings.HasPrefix(param, "%") {
			resolved = pr.formatTimeWithPattern(param)
		} else {
			// Use display version with masking for environment variables
			resolved = pr.resolveSpecialParameterForDisplay(param)
		}

		result += resolved
		segments = append(segments, highlight.Segment{
			Text:        resolved,
			IsHighlight: true,
		})

		lastEnd = end
	}

	// Add any remaining unchanged part
	if lastEnd < len(originalURL) {
		unchanged := originalURL[lastEnd:]
		result += unchanged
		segments = append(segments, highlight.Segment{
			Text:        unchanged,
			IsHighlight: false,
		})
	}

	return result, segments
}
