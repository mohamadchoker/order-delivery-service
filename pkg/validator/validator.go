package validator

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

var (
	// postalCodeRegex validates postal codes (flexible pattern)
	postalCodeRegex = regexp.MustCompile(`^[A-Z0-9\s-]{3,10}$`)
)

// ValidationError represents a validation error with field context
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
}

// ValidationErrors is a collection of validation errors
type ValidationErrors []ValidationError

func (e ValidationErrors) Error() string {
	var msgs []string
	for _, err := range e {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// HasErrors checks if there are any validation errors
func (e ValidationErrors) HasErrors() bool {
	return len(e) > 0
}

// Validator provides validation methods
type Validator struct {
	errors ValidationErrors
}

// New creates a new Validator
func New() *Validator {
	return &Validator{
		errors: make(ValidationErrors, 0),
	}
}

// AddError adds a validation error
func (v *Validator) AddError(field, message string) {
	v.errors = append(v.errors, ValidationError{
		Field:   field,
		Message: message,
	})
}

// HasErrors checks if there are any validation errors
func (v *Validator) HasErrors() bool {
	return v.errors.HasErrors()
}

// Errors returns all validation errors
func (v *Validator) Errors() error {
	if !v.HasErrors() {
		return nil
	}
	return v.errors
}

// ValidateRequired checks if a string field is not empty
func (v *Validator) ValidateRequired(field, value string) {
	if strings.TrimSpace(value) == "" {
		v.AddError(field, "is required")
	}
}

// ValidateStringLength validates string length
func (v *Validator) ValidateStringLength(field, value string, min, max int) {
	length := len(strings.TrimSpace(value))
	if length < min {
		v.AddError(field, fmt.Sprintf("must be at least %d characters", min))
	}
	if max > 0 && length > max {
		v.AddError(field, fmt.Sprintf("must not exceed %d characters", max))
	}
}

// ValidateTimeNotZero checks if time is not zero
func (v *Validator) ValidateTimeNotZero(field string, t time.Time) {
	if t.IsZero() {
		v.AddError(field, "is required")
	}
}

// ValidateTimeAfter checks if time is after another time
func (v *Validator) ValidateTimeAfter(field string, t, after time.Time, afterFieldName string) {
	if !t.IsZero() && !after.IsZero() && !t.After(after) {
		v.AddError(field, fmt.Sprintf("must be after %s", afterFieldName))
	}
}

// ValidateTimeFuture checks if time is in the future
func (v *Validator) ValidateTimeFuture(field string, t time.Time) {
	if !t.IsZero() && t.Before(time.Now()) {
		v.AddError(field, "must be in the future")
	}
}

// ValidateTimeRange validates time is within a reasonable range
func (v *Validator) ValidateTimeRange(field string, t time.Time, minDuration, maxDuration time.Duration) {
	if t.IsZero() {
		return
	}

	now := time.Now()
	diff := t.Sub(now)

	if diff < minDuration {
		v.AddError(field, fmt.Sprintf("must be at least %v from now", minDuration))
	}
	if maxDuration > 0 && diff > maxDuration {
		v.AddError(field, fmt.Sprintf("must not be more than %v from now", maxDuration))
	}
}

// ValidateAddress validates an address entity
func (v *Validator) ValidateAddress(fieldPrefix string, street, city, state, postalCode, country string, latitude, longitude float64) {
	if strings.TrimSpace(street) == "" {
		v.AddError(fieldPrefix+".street", "is required")
	}
	if strings.TrimSpace(city) == "" {
		v.AddError(fieldPrefix+".city", "is required")
	}
	if strings.TrimSpace(state) == "" {
		v.AddError(fieldPrefix+".state", "is required")
	}
	if strings.TrimSpace(postalCode) == "" {
		v.AddError(fieldPrefix+".postal_code", "is required")
	} else if !postalCodeRegex.MatchString(strings.ToUpper(postalCode)) {
		v.AddError(fieldPrefix+".postal_code", "is invalid")
	}
	if strings.TrimSpace(country) == "" {
		v.AddError(fieldPrefix+".country", "is required")
	}

	// Validate coordinates if provided
	if latitude != 0 || longitude != 0 {
		if latitude < -90 || latitude > 90 {
			v.AddError(fieldPrefix+".latitude", "must be between -90 and 90")
		}
		if longitude < -180 || longitude > 180 {
			v.AddError(fieldPrefix+".longitude", "must be between -180 and 180")
		}
	}
}

// ValidateUUID validates a UUID string
func (v *Validator) ValidateUUID(field, value string) {
	if strings.TrimSpace(value) == "" {
		v.AddError(field, "is required")
		return
	}

	// Simple UUID v4 validation pattern
	uuidPattern := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
	if !uuidPattern.MatchString(strings.ToLower(value)) {
		v.AddError(field, "is not a valid UUID")
	}
}

// ValidateEnum validates if value is in allowed list
func (v *Validator) ValidateEnum(field string, value interface{}, allowed []interface{}) {
	for _, a := range allowed {
		if value == a {
			return
		}
	}
	v.AddError(field, fmt.Sprintf("must be one of: %v", allowed))
}

// ValidatePositive checks if number is positive
func (v *Validator) ValidatePositive(field string, value int) {
	if value <= 0 {
		v.AddError(field, "must be positive")
	}
}

// ValidateRange validates if number is within range
func (v *Validator) ValidateRange(field string, value, min, max int) {
	if value < min {
		v.AddError(field, fmt.Sprintf("must be at least %d", min))
	}
	if max > 0 && value > max {
		v.AddError(field, fmt.Sprintf("must not exceed %d", max))
	}
}
