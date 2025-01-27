package util

import "strings"

// TruncateWithSuffix function truncates a string with a suffix
func TruncateWithSuffix(s string, maxLength int) string {
	if len(s) > maxLength {
		return s[:maxLength] + "..."
	}
	return s
}

// FindUntilFromRecurrence function finds the UNTIL value from the recurrence string
func FindUntilFromRecurrence(recurrence []string) string {
  if recurrence == nil {
    return ""
  }

  for _, r := range recurrence {
    if strings.HasPrefix(r, "RRULE:") {
      untilTxt := "UNTIL="
      untilLen := len(untilTxt)
      untilIndex := strings.Index(r,untilTxt)
      if untilIndex != -1 {
        // find the next semicolon
        semicolonIndex := strings.Index(r[untilIndex:], ";")
        if semicolonIndex != -1 {
          return r[untilIndex+untilLen:untilIndex+semicolonIndex]
        }
      }
    }
  }

  return ""
}
