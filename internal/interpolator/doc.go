// Package interpolator provides variable interpolation for environment maps.
//
// It expands references of the form ${VAR} or $VAR found in environment values
// by substituting them with the corresponding value from the same map or,
// optionally, from the host process environment.
//
// Example usage:
//
//	env := map[string]string{
//		"BASE_URL": "https://example.com",
//		"API_URL":  "${BASE_URL}/api/v1",
//	}
//
//	resolved, err := interpolator.Interpolate(env, interpolator.DefaultOptions())
//	// resolved["API_URL"] == "https://example.com/api/v1"
package interpolator
