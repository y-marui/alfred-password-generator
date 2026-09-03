package passgen

import (
	"strings"
	"testing"
)

const (
	asciiLower = "abcdefghijklmnopqrstuvwxyz"
	asciiUpper = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits     = "0123456789"

	// testMaxAttempts is generous enough that the class-diversity retry
	// always succeeds for the patterns/lengths used below, so tests can
	// assert on length/charset without worrying about the retry budget.
	testMaxAttempts = 200
)

func allIn(s, allowed string) bool {
	for _, c := range s {
		if !strings.ContainsRune(allowed, c) {
			return false
		}
	}
	return true
}

func hasClass(s, class string) bool {
	for _, c := range s {
		if strings.ContainsRune(class, c) {
			return true
		}
	}
	return false
}

func TestGenerateReturnsCorrectLength(t *testing.T) {
	pwd, err := Generate("A-Za-z0-9", 18, testMaxAttempts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len([]rune(pwd)) != 18 {
		t.Errorf("got length %d, want 18", len([]rune(pwd)))
	}
}

func TestGenerateCharactersInPattern(t *testing.T) {
	pwd, err := Generate("A-Za-z0-9", 100, testMaxAttempts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !allIn(pwd, asciiLower+asciiUpper+digits) {
		t.Errorf("password %q contains characters outside pattern", pwd)
	}
}

func TestGenerateExplicitChars(t *testing.T) {
	pwd, err := Generate("abc", 10, testMaxAttempts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !allIn(pwd, "abc") {
		t.Errorf("password %q contains characters outside 'abc'", pwd)
	}
}

func TestGenerateCustomLength(t *testing.T) {
	pwd, err := Generate("A-Z", 32, testMaxAttempts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len([]rune(pwd)) != 32 {
		t.Errorf("got length %d, want 32", len([]rune(pwd)))
	}
}

func TestGeneratePunctuationPattern(t *testing.T) {
	pwd, err := Generate("!-*", 20, testMaxAttempts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !allIn(pwd, "!@#^&*") {
		t.Errorf("password %q contains characters outside '!@#^&*'", pwd)
	}
}

func TestGenerateZeroLengthRaises(t *testing.T) {
	_, err := Generate("A-Z", 0, testMaxAttempts)
	if err == nil || !strings.Contains(err.Error(), "positive") {
		t.Fatalf("expected 'positive' error, got %v", err)
	}
}

func TestGenerateNegativeLengthRaises(t *testing.T) {
	_, err := Generate("A-Z", -1, testMaxAttempts)
	if err == nil || !strings.Contains(err.Error(), "positive") {
		t.Fatalf("expected 'positive' error, got %v", err)
	}
}

func TestGenerateInvalidRangeRaises(t *testing.T) {
	if _, err := Generate("z-a", 10, testMaxAttempts); err == nil {
		t.Fatal("expected error for reversed range 'z-a'")
	}
}

func TestGenerateDoubleDashRaises(t *testing.T) {
	if _, err := Generate("a--z", 10, testMaxAttempts); err == nil {
		t.Fatal("expected error for double dash 'a--z'")
	}
}

func TestGenerateTrailingDashRaises(t *testing.T) {
	_, err := Generate("A-Z-", 10, testMaxAttempts)
	if err == nil || !strings.Contains(err.Error(), "trailing") {
		t.Fatalf("expected 'trailing' error, got %v", err)
	}
}

func TestGenerateCrossClassRangeRaises(t *testing.T) {
	if _, err := Generate("A-z", 10, testMaxAttempts); err == nil {
		t.Fatal("expected error for cross-class range 'A-z'")
	}
}

func TestGenerateEnsuresClassDiversity(t *testing.T) {
	// A short length relative to the number of classes makes an
	// undiversified result likely on any single attempt, so this exercises
	// the retry loop rather than passing by chance.
	for i := 0; i < 50; i++ {
		pwd, err := Generate("A-Za-z0-9", 6, testMaxAttempts)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !hasClass(pwd, asciiLower) || !hasClass(pwd, asciiUpper) || !hasClass(pwd, digits) {
			t.Fatalf("password %q does not mix lower/upper/digit", pwd)
		}
	}
}

func TestGenerateDiversityBestEffortWhenAttemptsExhausted(t *testing.T) {
	// maxAttempts < 1 is clamped to 1: a single attempt must never error or
	// hang, even though it isn't guaranteed to be diverse.
	pwd, err := Generate("A-Za-z0-9", 3, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len([]rune(pwd)) != 3 {
		t.Errorf("got length %d, want 3", len([]rune(pwd)))
	}
}

func TestGenerateSplitCorrectFormat(t *testing.T) {
	pwd, err := GenerateSplit("A-Za-z0-9", 18, 6, testMaxAttempts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	parts := strings.Split(pwd, "-")
	if len(parts) != 3 {
		t.Fatalf("got %d parts, want 3", len(parts))
	}
	for _, p := range parts {
		if len([]rune(p)) != 6 {
			t.Errorf("part %q has length %d, want 6", p, len([]rune(p)))
		}
	}
}

func TestGenerateSplitCustomBy(t *testing.T) {
	pwd, err := GenerateSplit("A-Z", 12, 4, testMaxAttempts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	parts := strings.Split(pwd, "-")
	if len(parts) != 3 {
		t.Fatalf("got %d parts, want 3", len(parts))
	}
	for _, p := range parts {
		if len([]rune(p)) != 4 {
			t.Errorf("part %q has length %d, want 4", p, len([]rune(p)))
		}
	}
}

func TestGenerateSplitSingleGroup(t *testing.T) {
	pwd, err := GenerateSplit("A-Z", 6, 6, testMaxAttempts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(pwd, "-") {
		t.Errorf("password %q should not contain '-'", pwd)
	}
	if len([]rune(pwd)) != 6 {
		t.Errorf("got length %d, want 6", len([]rune(pwd)))
	}
}

func TestGenerateSplitLengthNotMultipleOfByRaises(t *testing.T) {
	_, err := GenerateSplit("A-Z", 18, 7, testMaxAttempts)
	if err == nil || !strings.Contains(err.Error(), "multiple") {
		t.Fatalf("expected 'multiple' error, got %v", err)
	}
}

func TestGenerateSplitLengthLessThanByRaises(t *testing.T) {
	_, err := GenerateSplit("A-Z", 3, 6, testMaxAttempts)
	if err == nil || !strings.Contains(err.Error(), ">=") {
		t.Fatalf("expected '>=' error, got %v", err)
	}
}

func TestGenerateSplitZeroByRaises(t *testing.T) {
	_, err := GenerateSplit("A-Z", 18, 0, testMaxAttempts)
	if err == nil || !strings.Contains(err.Error(), "positive") {
		t.Fatalf("expected 'positive' error, got %v", err)
	}
}

func TestGenerateSplitNegativeByRaises(t *testing.T) {
	_, err := GenerateSplit("A-Z", 18, -6, testMaxAttempts)
	if err == nil || !strings.Contains(err.Error(), "positive") {
		t.Fatalf("expected 'positive' error, got %v", err)
	}
}

func TestGenerateSplitZeroLengthRaises(t *testing.T) {
	_, err := GenerateSplit("A-Z", 0, 6, testMaxAttempts)
	if err == nil || !strings.Contains(err.Error(), "positive") {
		t.Fatalf("expected 'positive' error, got %v", err)
	}
}

func TestGenerateSplitCharactersInPattern(t *testing.T) {
	pwd, err := GenerateSplit("A-Za-z0-9", 18, 6, testMaxAttempts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	chars := strings.ReplaceAll(pwd, "-", "")
	if !allIn(chars, asciiLower+asciiUpper+digits) {
		t.Errorf("password %q contains characters outside pattern", pwd)
	}
}

func TestGenerateSplitEnsuresClassDiversityIgnoringHyphens(t *testing.T) {
	for i := 0; i < 50; i++ {
		pwd, err := GenerateSplit("A-Za-z0-9!-*", 12, 3, testMaxAttempts)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		chars := strings.ReplaceAll(pwd, "-", "")
		if !hasClass(chars, asciiLower) || !hasClass(chars, asciiUpper) ||
			!hasClass(chars, digits) || !hasClass(chars, "!@#^&*") {
			t.Fatalf("password %q does not mix all four classes", pwd)
		}
	}
}
