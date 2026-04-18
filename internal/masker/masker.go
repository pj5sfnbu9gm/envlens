package masker

import "strings"

// MaskOptions controls how sensitive values are masked.
type MaskOptions struct {
	// ShowChars is the number of trailing characters to reveal (0 = hide all).
	ShowChars int
	// Placeholder replaces the hidden portion of the value.
	Placeholder string
}

// DefaultMaskOptions returns sensible masking defaults.
func DefaultMaskOptions() MaskOptions {
	return MaskOptions{
		ShowChars:   4,
		Placeholder: "****",
	}
}

// sensitivePatterns are substrings that indicate a key holds a secret.
var sensitivePatterns = []string{
	"SECRET", "PASSWORD", "PASSWD", "TOKEN", "API_KEY",
	"PRIVATE", "CREDENTIAL", "AUTH", "CERT", "KEY",
}

// IsSensitive returns true when the key name matches a known secret pattern.
func IsSensitive(key string) bool {
	upper := strings.ToUpper(key)
	for _, p := range sensitivePatterns {
		if strings.Contains(upper, p) {
			return true
		}
	}
	return false
}

// Mask redacts a value according to the given options.
// If the value is shorter than or equal to ShowChars the entire value is replaced.
func Mask(value string, opts MaskOptions) string {
	if opts.ShowChars <= 0 || len(value) <= opts.ShowChars {
		return opts.Placeholder
	}
	visible := value[len(value)-opts.ShowChars:]
	return opts.Placeholder + visible
}

// MaskEnv returns a copy of the env map with sensitive values masked.
func MaskEnv(env map[string]string, opts MaskOptions) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		if IsSensitive(k) {
			out[k] = Mask(v, opts)
		} else {
			out[k] = v
		}
	}
	return out
}

// SensitiveKeys returns the list of keys from env that are considered sensitive.
func SensitiveKeys(env map[string]string) []string {
	keys := make([]string, 0)
	for k := range env {
		if IsSensitive(k) {
			keys = append(keys, k)
		}
	}
	return keys
}
