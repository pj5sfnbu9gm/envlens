// Package pivotter transposes a multi-target environment comparison into
// a key-centric view.
//
// Instead of "target X differs from baseline in these keys", pivotter
// answers "for key K, what value does each target have?"
//
// This is useful for spotting configuration drift across many deployment
// targets at a glance.
package pivotter
