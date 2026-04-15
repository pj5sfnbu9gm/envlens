// Package masker provides utilities for detecting and redacting sensitive
// environment variable values before they are displayed or written to output.
//
// A key is considered sensitive when its name contains well-known substrings
// such as PASSWORD, SECRET, TOKEN, or API_KEY (case-insensitive).
//
// Example usage:
//
//	env := map[string]string{
//		"DB_PASSWORD": "s3cr3t",
//		"APP_NAME":    "myapp",
//	}
//	masked := masker.MaskEnv(env, masker.DefaultMaskOptions())
//	// masked["DB_PASSWORD"] == "****3cr3t" (trailing 4 chars visible)
//	// masked["APP_NAME"]    == "myapp"     (unchanged)
package masker
