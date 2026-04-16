package typecast

import (
	"testing"
)

func TestCast_StringDefault(t *testing.T) {
	env := map[string]string{"NAME": "alice"}
	results, err := Cast(env, DefaultOptions())
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 || results[0].Value != "alice" {
		t.Fatalf("expected alice, got %v", results)
	}
}

func TestCast_Int(t *testing.T) {
	env := map[string]string{"PORT": "8080"}
	opts := DefaultOptions()
	opts.Hints = map[string]string{"PORT": "int"}
	results, err := Cast(env, opts)
	if err != nil {
		t.Fatal(err)
	}
	if results[0].Value.(int) != 8080 {
		t.Fatalf("expected 8080")
	}
}

func TestCast_Float(t *testing.T) {
	env := map[string]string{"RATIO": "3.14"}
	opts := DefaultOptions()
	opts.Hints = map[string]string{"RATIO": "float"}
	results, err := Cast(env, opts)
	if err != nil {
		t.Fatal(err)
	}
	if results[0].Value.(float64) != 3.14 {
		t.Fatalf("expected 3.14")
	}
}

func TestCast_Bool(t *testing.T) {
	env := map[string]string{"DEBUG": "true"}
	opts := DefaultOptions()
	opts.Hints = map[string]string{"DEBUG": "bool"}
	results, err := Cast(env, opts)
	if err != nil {
		t.Fatal(err)
	}
	if results[0].Value.(bool) != true {
		t.Fatalf("expected true")
	}
}

func TestCast_InvalidInt_Lenient(t *testing.T) {
	env := map[string]string{"PORT": "abc"}
	opts := DefaultOptions()
	opts.Hints = map[string]string{"PORT": "int"}
	results, err := Cast(env, opts)
	if err != nil {
		t.Fatal("expected no error in lenient mode")
	}
	if results[0].Error == "" {
		t.Fatal("expected error message in result")
	}
}

func TestCast_InvalidInt_Strict(t *testing.T) {
	env := map[string]string{"PORT": "abc"}
	opts := DefaultOptions()
	opts.Hints = map[string]string{"PORT": "int"}
	opts.StrictMode = true
	_, err := Cast(env, opts)
	if err == nil {
		t.Fatal("expected error in strict mode")
	}
}

func TestCast_EmptyEnv(t *testing.T) {
	results, err := Cast(map[string]string{}, DefaultOptions())
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 0 {
		t.Fatalf("expected empty results")
	}
}
