package database

import (
	"errors"
	"regexp"
	"strings"

	"github.com/arangodb/go-driver"
)

// DBError represents a parsed database error with user-friendly information
type DBError struct {
	Code    int    // ArangoDB error code
	Field   string // The field that caused the error (if applicable)
	Message string // User-friendly error message
	Raw     error  // Original error
}

func (e *DBError) Error() string {
	return e.Message
}

// Common ArangoDB error codes
const (
	ErrUniqueConstraintViolated = 1210
	ErrDocumentNotFound         = 1202
)

// ParseError analyzes an ArangoDB error and returns a structured DBError
func ParseError(err error) *DBError {
	if err == nil {
		return nil
	}

	dbErr := &DBError{
		Raw:     err,
		Message: "An unexpected error occurred",
	}

	// Check if it's an ArangoDB error
	var arangoErr driver.ArangoError
	if errors.As(err, &arangoErr) {
		dbErr.Code = arangoErr.ErrorNum

		switch arangoErr.ErrorNum {
		case ErrUniqueConstraintViolated:
			dbErr.Field = extractFieldFromUniqueError(arangoErr.ErrorMessage)
			if dbErr.Field != "" {
				dbErr.Message = dbErr.Field + " already exists"
			} else {
				dbErr.Message = "A record with this value already exists"
			}
		case ErrDocumentNotFound:
			dbErr.Message = "Record not found"
		default:
			dbErr.Message = arangoErr.ErrorMessage
		}

		return dbErr
	}

	// Fallback: check error string for known patterns
	errStr := err.Error()

	if strings.Contains(errStr, "unique constraint") || strings.Contains(errStr, "duplicate") {
		dbErr.Code = ErrUniqueConstraintViolated
		dbErr.Field = extractFieldFromUniqueError(errStr)
		if dbErr.Field != "" {
			dbErr.Message = dbErr.Field + " already exists"
		} else {
			dbErr.Message = "A record with this value already exists"
		}
	}

	return dbErr
}

// extractFieldFromUniqueError attempts to extract the field name from unique constraint error messages
// ArangoDB errors look like: "unique constraint violated - in index ... over 'fieldname'" or "over 'field1', 'field2'"
func extractFieldFromUniqueError(errMsg string) string {
	// Pattern: over 'fieldname' or over 'field1', 'field2'
	re := regexp.MustCompile(`over ['"]([^'"]+)['"]`)
	matches := re.FindStringSubmatch(errMsg)
	if len(matches) > 1 {
		return formatFieldName(matches[1])
	}

	// Alternative pattern: field names might appear differently
	// Check for common field names in the error
	commonFields := []string{"email", "username", "slug", "name"}
	lowerErr := strings.ToLower(errMsg)
	for _, field := range commonFields {
		if strings.Contains(lowerErr, field) {
			return formatFieldName(field)
		}
	}

	return ""
}

// formatFieldName converts field names to user-friendly format
// e.g., "email" -> "Email", "user_name" -> "User name"
func formatFieldName(field string) string {
	if field == "" {
		return ""
	}

	// Replace underscores with spaces
	field = strings.ReplaceAll(field, "_", " ")

	// Capitalize first letter
	return strings.ToUpper(field[:1]) + field[1:]
}

// IsUniqueConstraintError checks if the error is a unique constraint violation
func IsUniqueConstraintError(err error) bool {
	if err == nil {
		return false
	}

	var arangoErr driver.ArangoError
	if errors.As(err, &arangoErr) {
		return arangoErr.ErrorNum == ErrUniqueConstraintViolated
	}

	errStr := err.Error()
	return strings.Contains(errStr, "unique constraint") || strings.Contains(errStr, "duplicate")
}

// IsNotFoundError checks if the error is a document not found error
func IsNotFoundError(err error) bool {
	if err == nil {
		return false
	}

	var arangoErr driver.ArangoError
	if errors.As(err, &arangoErr) {
		return arangoErr.ErrorNum == ErrDocumentNotFound
	}

	return strings.Contains(err.Error(), "not found")
}

