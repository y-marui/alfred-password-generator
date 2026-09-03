// Package passgencmd dispatches an Alfred passgen query to the matching
// command handler and builds the Script Filter response.
//
// Commands:
//
//	passgen [length] [pattern]             — basic password (default)
//	panc [split] [length] [by] [pattern]   — with punctuation
//	split [length] [by] [pattern]          — split by groups
//	help                                   — show available commands
package passgencmd

import (
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
	numSuggestions = 5
)

// overviewEntry is one row of the bare-length overview shown by handleBasic
// when the query is empty or a length with no pattern.
type overviewEntry struct {
	name    string
	pattern string
	split   bool
}

var overview = []overviewEntry{
	{"basic", patternBasic, false},
	{"panc", patternPanc, false},
	{"split", patternBasic, true},
	{"panc split", patternPanc, true},
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
		return runBasic(args, patternBasic, 1)
	}
	if len(parts) == 1 {
		return showOverview(length)
	}
	return runBasic(args, patternBasic, 1)
}

// handlePanc implements "passgen panc [...]" — with punctuation; optionally
// split.
func handlePanc(args string) scriptfilter.Response {
	first, rest := splitFirst(args)
	if strings.ToLower(first) == "split" {
		return runSplit(rest, patternPanc)
	}
	return runBasic(args, patternPanc, numSuggestions)
}

// handleSplit implements "passgen split [...]" — split password without
// punctuation.
func handleSplit(args string) scriptfilter.Response {
	return runSplit(args, patternBasic)
}

func showOverview(length int) scriptfilter.Response {
	var items []scriptfilter.Item
	for i, e := range overview {
		var pwd, subtitle string
		var err error
		if e.split {
			pwd, err = passgen.GenerateSplit(e.pattern, length, defaultBy)
			subtitle = e.name + " · " + strconv.Itoa(length) + " chars in groups of " + strconv.Itoa(defaultBy)
		} else {
			pwd, err = passgen.Generate(e.pattern, length)
			subtitle = e.name + " · " + strconv.Itoa(length) + " chars"
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

func runBasic(args, defaultPattern string, count int) scriptfilter.Response {
	length, pattern := parseBasic(args, defaultPattern)
	subtitle := strconv.Itoa(length) + " chars, pattern: " + pattern
	return output(func() (string, error) { return passgen.Generate(pattern, length) }, subtitle, count)
}

func runSplit(args, defaultPattern string) scriptfilter.Response {
	length, by, pattern := parseSplit(args, defaultPattern)
	subtitle := strconv.Itoa(length) + " chars in groups of " + strconv.Itoa(by) + ", pattern: " + pattern
	return output(func() (string, error) { return passgen.GenerateSplit(pattern, length, by) }, subtitle, numSuggestions)
}

func parseBasic(args, defaultPattern string) (length int, pattern string) {
	fields := strings.Fields(args)
	if len(fields) == 0 {
		return defaultLength, defaultPattern
	}
	n, err := strconv.Atoi(fields[0])
	if err != nil {
		return defaultLength, defaultPattern
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

var helpCommands = []struct {
	cmd, desc, autocomplete string
}{
	{"passgen [length] [pattern]", "Generate password (default: 18 chars, A-Za-z0-9)", ""},
	{"passgen panc [length] [pattern]", "Generate with punctuation (!@#^&*)", "panc "},
	{"passgen split [length] [by] [pattern]", "Generate split password (default: 18/6)", "split "},
	{"passgen panc split [length] [by] [pattern]", "Split with punctuation", "panc split "},
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
