package config

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	wantPort := 3333
	t.Setenv("BACKEND_PORT", fmt.Sprint(wantPort))

	got, err := New()
	if err != nil {
		t.Fatalf("cannot create config: %v", err)
	}
	if got.BckendPort != wantPort {
		t.Errorf("want %d, but %d", wantPort, got.BckendPort)
	}

	wantEnv := "dev"
	if got.BackendEnv != wantEnv {
		t.Errorf("want %s, but %s", wantEnv, got.BackendEnv)
	}
}
