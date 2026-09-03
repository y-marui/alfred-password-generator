// Package passgencmd dispatches an Alfred passgen query to the matching
// command handler and builds the Script Filter response.
//
// Commands:
//
//	passgen [length] [pattern]             — basic password (default)
//	panc [split] [length] [by] [pattern]   — with punctuation
//	split [length] [by] [pattern]          — split by groups
//	pin [length]                           — numeric PIN (default: 4 digits)
//	code [length]                          — numeric code (default: 6 digits)
//	help                                   — show available commands
package passgencmd

import (
	"os"
	"strconv"
	"strings"

	"github.com/y-marui/alfred-password-generator/internal/passgen"
	"github.com/y-marui/alfred-password-generator/internal/scriptfilter"
)

const (
	defaultLength  = 18
	defaultBy      = 6
	patternBasic   = "A-Za-z0-9"
	patternPanc    = "A-Za-z0-9!-*"
	patternDigits  = "0-9"
	pinLength      = 4
	codeLength     = 6
	numSuggestions = 5

	// defaultMaxAttempts is how many times passgen retries a generated
	// password to get a mix of every character class present in its
	// pattern, when the max_attempts Config Builder variable is unset or
	// invalid. Configurable because a very restrictive custom pattern
	// combined with a short length can make that mix take more attempts.
	defaultMaxAttempts = 100
)

// overviewEntry is one row of the bare-length overview shown by handleBasic
// when the query is empty or a length with no pattern.
type overviewEntry struct {
	name    string
	pattern string
	split   bool
	// fixedLength overrides the overview's requested length when non-zero
	// (used by pin/code, which have their own conventional lengths).
	fixedLength int
}

var overview = []overviewEntry{
	{name: "pin", pattern: patternDigits, fixedLength: pinLength},
	{name: "code", pattern: patternDigits, fixedLength: codeLength},
	{name: "basic", pattern: patternBasic},
	{name: "panc", pattern: patternPanc},
	{name: "split", pattern: patternBasic, split: true},
	{name: "panc split", pattern: patternPanc, split: true},
}

// Dispatch parses the raw Alfred query and routes it to the matching
// command, falling back to the basic passgen command when the leading
// token isn't a registered command name (e.g. "18 A-Z" or "24").
func Dispatch(query string) scriptfilter.Response {
	command, args := splitFirst(query)
	switch strings.ToLower(command) {
	case "panc":
		return handlePanc(args)
	case "split":
		return handleSplit(args)
	case "pin":
		return handlePin(args)
	case "code":
		return handleCode(args)
	case "help":
		return handleHelp()
	case "passgen":
		return handleBasic(args)
	default:
		return handleBasic(strings.TrimSpace(query))
	}
}

func splitFirst(query string) (first, rest string) {
	fields := strings.Fields(query)
	if len(fields) == 0 {
		return "", ""
	}
	return fields[0], strings.Join(fields[1:], " ")
}

// handleBasic implements "passgen [length] [pattern]" — overview or custom
// password.
func handleBasic(args string) scriptfilter.Response {
	parts := strings.Fields(args)
	if len(parts) == 0 {
		return showOverview(defaultLength)
	}
	length, err := strconv.Atoi(parts[0])
	if err != nil {
		// A non-integer leading token means this isn't "[length]
		// [pattern]" overview syntax; fall through to a single custom
		// result, same as when a pattern follows a valid length below.
		return runBasic(args, patternBasic, defaultLength, 1)
	}
	if len(parts) == 1 {
		return showOverview(length)
	}
	return runBasic(args, patternBasic, defaultLength, 1)
}

// handlePanc implements "passgen panc [...]" — with punctuation; optionally
// split.
func handlePanc(args string) scriptfilter.Response {
	first, rest := splitFirst(args)
	if strings.ToLower(first) == "split" {
		return runSplit(rest, patternPanc)
	}
	return runBasic(args, patternPanc, defaultLength, numSuggestions)
}

// handleSplit implements "passgen split [...]" — split password without
// punctuation.
func handleSplit(args string) scriptfilter.Response {
	return runSplit(args, patternBasic)
}

// handlePin implements "passgen pin [length] [pattern]" — a numeric PIN,
// 4 digits by default.
func handlePin(args string) scriptfilter.Response {
	return runBasic(args, patternDigits, pinLength, numSuggestions)
}

// handleCode implements "passgen code [length] [pattern]" — a numeric
// code, 6 digits by default.
func handleCode(args string) scriptfilter.Response {
	return runBasic(args, patternDigits, codeLength, numSuggestions)
}

func showOverview(length int) scriptfilter.Response {
	var items []scriptfilter.Item
	attempts := maxAttempts()
	for i, e := range overview {
		effLength := length
		if e.fixedLength > 0 {
			effLength = e.fixedLength
		}
		var pwd, subtitle string
		var err error
		if e.split {
			pwd, err = passgen.GenerateSplit(e.pattern, effLength, defaultBy, attempts)
			subtitle = e.name + " · " + strconv.Itoa(effLength) + " chars in groups of " + strconv.Itoa(defaultBy)
		} else {
			pwd, err = passgen.Generate(e.pattern, effLength, attempts)
			subtitle = e.name + " · " + strconv.Itoa(effLength) + " chars"
		}
		if err != nil {
			continue
		}
		items = append(items, scriptfilter.Item{
			Title:    pwd,
			Subtitle: subtitle,
			Arg:      pwd,
			UID:      "passgen-" + strconv.Itoa(i),
			Valid:    scriptfilter.BoolPtr(true),
		})
	}
	if len(items) == 0 {
		return scriptfilter.Response{Items: []scriptfilter.Item{errorItem("No valid patterns for length " + strconv.Itoa(length))}}
	}
	return scriptfilter.Response{Items: items, SkipKnowledge: true}
}

func runBasic(args, defaultPattern string, defaultLen, count int) scriptfilter.Response {
	length, pattern := parseBasic(args, defaultPattern, defaultLen)
	subtitle := strconv.Itoa(length) + " chars, pattern: " + pattern
	attempts := maxAttempts()
	return output(func() (string, error) { return passgen.Generate(pattern, length, attempts) }, subtitle, count)
}

func runSplit(args, defaultPattern string) scriptfilter.Response {
	length, by, pattern := parseSplit(args, defaultPattern)
	subtitle := strconv.Itoa(length) + " chars in groups of " + strconv.Itoa(by) + ", pattern: " + pattern
	attempts := maxAttempts()
	return output(func() (string, error) { return passgen.GenerateSplit(pattern, length, by, attempts) }, subtitle, numSuggestions)
}

func parseBasic(args, defaultPattern string, defaultLen int) (length int, pattern string) {
	fields := strings.Fields(args)
	if len(fields) == 0 {
		return defaultLen, defaultPattern
	}
	n, err := strconv.Atoi(fields[0])
	if err != nil {
		return defaultLen, defaultPattern
	}
	pattern = defaultPattern
	if len(fields) > 1 {
		pattern = strings.Join(fields[1:], " ")
	}
	return n, pattern
}

func parseSplit(args, defaultPattern string) (length, by int, pattern string) {
	fields := strings.Fields(args)
	length, by, pattern = defaultLength, defaultBy, defaultPattern

	if len(fields) >= 1 {
		if n, err := strconv.Atoi(fields[0]); err == nil {
			length = n
		}
	}
	if len(fields) >= 2 {
		if n, err := strconv.Atoi(fields[1]); err == nil {
			by = n
		}
	}
	if len(fields) >= 3 {
		pattern = strings.Join(fields[2:], " ")
	}

	return length, by, pattern
}

func output(genFn func() (string, error), subtitle string, count int) scriptfilter.Response {
	items := make([]scriptfilter.Item, 0, count)
	for i := 0; i < count; i++ {
		pwd, err := genFn()
		if err != nil {
			return scriptfilter.Response{Items: []scriptfilter.Item{errorItem(err.Error())}}
		}
		items = append(items, scriptfilter.Item{
			Title:    pwd,
			Subtitle: subtitle,
			Arg:      pwd,
			UID:      "passgen-" + strconv.Itoa(i),
			Valid:    scriptfilter.BoolPtr(true),
		})
	}
	return scriptfilter.Response{Items: items, SkipKnowledge: true}
}

func errorItem(message string) scriptfilter.Item {
	return scriptfilter.Item{
		Title:    "Error: " + message,
		Subtitle: "Press ⌘C to copy the error",
		Arg:      message,
		Valid:    scriptfilter.BoolPtr(false),
	}
}

// maxAttempts reads the max_attempts Config Builder variable, falling back
// to defaultMaxAttempts when it's unset, empty, or not a positive integer.
func maxAttempts() int {
	raw := os.Getenv("max_attempts")
	if raw == "" {
		return defaultMaxAttempts
	}
	n, err := strconv.Atoi(raw)
	if err != nil || n < 1 {
		return defaultMaxAttempts
	}
	return n
}

var helpCommands = []struct {
	cmd, desc, autocomplete string
}{
	{"passgen [length] [pattern]", "Generate password (default: 18 chars, A-Za-z0-9)", ""},
	{"passgen panc [length] [pattern]", "Generate with punctuation (!@#^&*)", "panc "},
	{"passgen split [length] [by] [pattern]", "Generate split password (default: 18/6)", "split "},
	{"passgen panc split [length] [by] [pattern]", "Split with punctuation", "panc split "},
	{"passgen pin [length]", "Generate a numeric PIN (default: 4 digits)", "pin "},
	{"passgen code [length]", "Generate a numeric code (default: 6 digits)", "code "},
	{"passgen help", "Show this help", "help"},
}

func handleHelp() scriptfilter.Response {
	items := make([]scriptfilter.Item, len(helpCommands))
	for i, c := range helpCommands {
		items[i] = scriptfilter.Item{
			Title:        c.cmd,
			Subtitle:     c.desc,
			Arg:          "",
			UID:          "help-" + strconv.Itoa(i),
			Valid:        scriptfilter.BoolPtr(false),
			Autocomplete: c.autocomplete,
		}
	}
	return scriptfilter.Response{Items: items}
}
