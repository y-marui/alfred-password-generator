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

// ranges are the character classes a "-" range may span within. A range
// cannot mix characters from two different classes (e.g. "A-z" is invalid).
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

func rangeContaining(a, b rune) string {
	for _, r := range ranges {
		if strings.ContainsRune(r, a) && strings.ContainsRune(r, b) {
			return r
		}
	}
	return ""
}

// Generate returns a cryptographically secure random password of the given
// length, drawn uniformly from pattern's expanded character set.
func Generate(pattern string, length int) (string, error) {
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

// GenerateSplit returns a secure password split into groups of by
// characters, joined by hyphens (e.g. "xxxxxx-xxxxxx-xxxxxx").
func GenerateSplit(pattern string, length, by int) (string, error) {
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

	groups := make([]string, length/by)
	for g := range groups {
		group := make([]rune, by)
		for i := range group {
			c, err := secureChoice(chars)
			if err != nil {
				return "", err
			}
			group[i] = c
		}
		groups[g] = string(group)
	}
	return strings.Join(groups, "-"), nil
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
