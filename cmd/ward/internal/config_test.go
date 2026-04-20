package internal

import "testing"

func TestIsQuietDefault(t *testing.T) {
	quietMode.Store(false)
	if IsQuiet() {
		t.Error("expected false")
	}
}

func TestIsQuietEnabled(t *testing.T) {
	defer quietMode.Store(false)
	quietMode.Store(true)
	if !IsQuiet() {
		t.Error("expected true")
	}
}

func TestIsDebugDefault(t *testing.T) {
	debugMode.Store(false)
	if IsDebug() {
		t.Error("expected false")
	}
}

func TestIsDebugEnabled(t *testing.T) {
	defer debugMode.Store(false)
	debugMode.Store(true)
	if !IsDebug() {
		t.Error("expected true")
	}
}
