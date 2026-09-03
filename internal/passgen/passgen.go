// Package passgen generates cryptographically secure passwords from a
// character-range pattern (e.g. "A-Za-z0-9" or "A-Za-z0-9!-*").
package passgen

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
)

// Punctuation is the character set the "!-*" range expands to.
const Punctuation = "!@#^&*"

// numClasses is len(ranges), fixed as a constant so it can size arrays.
const numClasses = 4

// ranges are the character classes a "-" range may span within. A range
// cannot mix characters from two different classes (e.g. "A-z" is invalid).
// A character outside all of these (e.g. the "-" joiner GenerateSplit
// inserts) belongs to no class and never affects the diversity check below.
var ranges = []string{
	"abcdefghijklmnopqrstuvwxyz",
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ",
	"0123456789",
	Punctuation,
}

// ExpandPattern expands a character-range pattern into the set of
// individual characters it selects from. Characters can be listed directly
// (e.g. "ABCabc012") or as ranges (e.g. "A-Z"); ranges within different
// character classes cannot be mixed (e.g. "A-z" is invalid).
func ExpandPattern(pattern string) ([]rune, error) {
	var out []rune
	var last rune
	haveLast := false
	pendingDash := false

	for _, c := range pattern {
		if c == '-' {
			if pendingDash {
				return nil, fmt.Errorf("-- is not a valid pattern")
			}
			pendingDash = true
			continue
		}

		flag := pendingDash
		pendingDash = false
		prevLast, prevHaveLast := last, haveLast
		last, haveLast = c, true

		if !flag {
			out = append(out, c)
			continue
		}

		if !prevHaveLast {
			return nil, fmt.Errorf("-%c is not a valid pattern", c)
		}

		source := rangeContaining(prevLast, c)
		if source == "" {
			return nil, fmt.Errorf("%c-%c is not a valid range", prevLast, c)
		}
		sourceRunes := []rune(source)
		iStart := strings.IndexRune(source, prevLast)
		iEnd := strings.IndexRune(source, c)
		if iStart > iEnd {
			return nil, fmt.Errorf("%c-%c is not a valid range", prevLast, c)
		}
		out = append(out, sourceRunes[iStart+1:iEnd+1]...)
	}

	if pendingDash {
		return nil, fmt.Errorf("pattern ends with '-'; trailing dash is not valid")
	}

	return out, nil
}

// classOf returns the index into ranges that c belongs to, or -1 if c is
// outside every defined class (e.g. a literal "-" joiner).
func classOf(c rune) int {
	for i, r := range ranges {
		if strings.ContainsRune(r, c) {
			return i
		}
	}
	return -1
}

func rangeContaining(a, b rune) string {
	i := classOf(a)
	if i < 0 || !strings.ContainsRune(ranges[i], b) {
		return ""
	}
	return ranges[i]
}

// requiredClasses reports, for each class in ranges, whether chars contains
// at least one character from it — i.e. which classes a generated password
// must represent to be considered diverse.
func requiredClasses(chars []rune) [numClasses]bool {
	var required [numClasses]bool
	for _, c := range chars {
		if i := classOf(c); i >= 0 {
			required[i] = true
		}
	}
	return required
}

// satisfiesClasses reports whether s contains at least one character from
// every class flagged in required.
func satisfiesClasses(s string, required [numClasses]bool) bool {
	var present [numClasses]bool
	for _, c := range s {
		if i := classOf(c); i >= 0 {
			present[i] = true
		}
	}
	for i, need := range required {
		if need && !present[i] {
			return false
		}
	}
	return true
}

// Generate returns a cryptographically secure random password of the given
// length, drawn uniformly from pattern's expanded character set. If that
// character set spans multiple classes (lowercase/uppercase/digit/punctuation),
// Generate retries — up to maxAttempts times — until the result contains at
// least one character from each class, falling back to the last attempt if
// none qualifies within that budget (e.g. because length is too short for
// the number of classes involved). maxAttempts below 1 is treated as 1.
func Generate(pattern string, length, maxAttempts int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("length must be a positive integer, got %d", length)
	}
	chars, err := ExpandPattern(pattern)
	if err != nil {
		return "", err
	}
	if len(chars) == 0 {
		return "", fmt.Errorf("pattern '%s' produces no characters", pattern)
	}
	required := requiredClasses(chars)

	var pwd string
	for attempt := 0; attempt < max(maxAttempts, 1); attempt++ {
		pwd, err = generateOnce(chars, length)
		if err != nil {
			return "", err
		}
		if satisfiesClasses(pwd, required) {
			break
		}
	}
	return pwd, nil
}

// GenerateSplit returns a secure password split into groups of by
// characters, joined by hyphens (e.g. "xxxxxx-xxxxxx-xxxxxx"). The same
// class-diversity retry as Generate applies to the full password (the "-"
// joiners are never part of any class, so they never affect the check).
func GenerateSplit(pattern string, length, by, maxAttempts int) (string, error) {
	if by <= 0 {
		return "", fmt.Errorf("by must be a positive integer, got %d", by)
	}
	if length <= 0 {
		return "", fmt.Errorf("length must be a positive integer, got %d", length)
	}
	if length < by {
		return "", fmt.Errorf("length must be >= by: length=%d, by=%d", length, by)
	}
	if length%by != 0 {
		return "", fmt.Errorf("length must be a multiple of by: length=%d, by=%d", length, by)
	}
	chars, err := ExpandPattern(pattern)
	if err != nil {
		return "", err
	}
	if len(chars) == 0 {
		return "", fmt.Errorf("pattern '%s' produces no characters", pattern)
	}
	required := requiredClasses(chars)

	var result string
	for attempt := 0; attempt < max(maxAttempts, 1); attempt++ {
		groups := make([]string, length/by)
		for g := range groups {
			part, err := generateOnce(chars, by)
			if err != nil {
				return "", err
			}
			groups[g] = part
		}
		result = strings.Join(groups, "-")
		if satisfiesClasses(result, required) {
			break
		}
	}
	return result, nil
}

// generateOnce draws length characters uniformly at random from chars.
func generateOnce(chars []rune, length int) (string, error) {
	out := make([]rune, length)
	for i := range out {
		c, err := secureChoice(chars)
		if err != nil {
			return "", err
		}
		out[i] = c
	}
	return string(out), nil
}

// secureChoice picks a uniformly random rune from chars using a
// cryptographically secure source, matching the security level of Python's
// secrets.choice.
func secureChoice(chars []rune) (rune, error) {
	idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
	if err != nil {
		return 0, fmt.Errorf("generating random index: %w", err)
	}
	return chars[idx.Int64()], nil
}
