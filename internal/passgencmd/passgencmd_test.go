package passgencmd

import (
	"strings"
	"testing"

	"github.com/y-marui/alfred-password-generator/internal/scriptfilter"
)

func TestEmptyQueryShowsOverview(t *testing.T) {
	resp := Dispatch("")
	if len(resp.Items) != len(overview) {
		t.Fatalf("got %d items, want %d", len(resp.Items), len(overview))
	}
	joined := subtitles(resp)
	for _, want := range []string{"basic", "panc", "split"} {
		if !strings.Contains(joined, want) {
			t.Errorf("expected a subtitle containing %q, got %v", want, joined)
		}
	}
}

func TestLengthOnlyShowsOverview(t *testing.T) {
	resp := Dispatch("24")
	if len(resp.Items) != len(overview) {
		t.Fatalf("got %d items, want %d", len(resp.Items), len(overview))
	}
	for _, it := range resp.Items {
		if !strings.Contains(it.Subtitle, "24") {
			t.Errorf("subtitle %q does not mention length 24", it.Subtitle)
		}
	}
}

func TestNonDivisibleLengthSkipsSplit(t *testing.T) {
	// 20 is not divisible by 6, so split variants are skipped.
	resp := Dispatch("20")
	for _, it := range resp.Items {
		if strings.Contains(it.Subtitle, "split") {
			t.Errorf("subtitle %q should not mention split", it.Subtitle)
		}
	}
}

func TestCustomPatternReturnsSingleResult(t *testing.T) {
	resp := Dispatch("18 A-Z")
	if len(resp.Items) != 1 {
		t.Fatalf("got %d items, want 1", len(resp.Items))
	}
	chars := resp.Items[0].Title
	if len([]rune(chars)) != 18 {
		t.Errorf("got length %d, want 18", len([]rune(chars)))
	}
	for _, c := range chars {
		if c < 'A' || c > 'Z' {
			t.Errorf("password %q contains non-uppercase character %q", chars, c)
		}
	}
}

func TestPancIncludesPunctuationByDefault(t *testing.T) {
	resp := Dispatch("panc")
	if len(resp.Items) != numSuggestions {
		t.Fatalf("got %d items, want %d", len(resp.Items), numSuggestions)
	}
	var all strings.Builder
	for _, it := range resp.Items {
		all.WriteString(it.Title)
	}
	// 5 * 18 = 90 chars from a 68-char charset; probability of zero punctuation is negligible.
	if !strings.ContainsAny(all.String(), "!@#^&*") {
		t.Errorf("expected at least one punctuation character across results, got %q", all.String())
	}
}

func TestPancSplit(t *testing.T) {
	resp := Dispatch("panc split 18 6")
	for _, it := range resp.Items {
		parts := strings.Split(it.Title, "-")
		if len(parts) != 3 {
			t.Fatalf("got %d parts, want 3 (title=%q)", len(parts), it.Title)
		}
		for _, p := range parts {
			if len([]rune(p)) != 6 {
				t.Errorf("part %q has length %d, want 6", p, len([]rune(p)))
			}
		}
	}
}

func TestPancSplitNoArgs(t *testing.T) {
	resp := Dispatch("panc split")
	if len(resp.Items) != numSuggestions {
		t.Fatalf("got %d items, want %d", len(resp.Items), numSuggestions)
	}
	for _, it := range resp.Items {
		parts := strings.Split(it.Title, "-")
		if len(parts) != defaultLength/defaultBy {
			t.Errorf("got %d parts, want %d", len(parts), defaultLength/defaultBy)
		}
	}
}

func TestSplitCommand(t *testing.T) {
	resp := Dispatch("split 18 6")
	for _, it := range resp.Items {
		parts := strings.Split(it.Title, "-")
		if len(parts) != 3 {
			t.Fatalf("got %d parts, want 3", len(parts))
		}
		for _, p := range parts {
			if len([]rune(p)) != 6 {
				t.Errorf("part %q has length %d, want 6", p, len([]rune(p)))
			}
		}
	}
}

func TestSplitDefaultArgs(t *testing.T) {
	resp := Dispatch("split")
	if len(resp.Items) != numSuggestions {
		t.Fatalf("got %d items, want %d", len(resp.Items), numSuggestions)
	}
}

func TestSplitCustomPattern(t *testing.T) {
	resp := Dispatch("split 12 4 A-Z")
	for _, it := range resp.Items {
		chars := strings.ReplaceAll(it.Title, "-", "")
		if len([]rune(chars)) != 12 {
			t.Errorf("got length %d, want 12", len([]rune(chars)))
		}
		for _, c := range chars {
			if c < 'A' || c > 'Z' {
				t.Errorf("password %q contains non-uppercase character %q", chars, c)
			}
		}
	}
}

func TestInvalidSplitRatioReturnsError(t *testing.T) {
	resp := Dispatch("split 18 7")
	assertSingleError(t, resp)
}

func TestZeroByReturnsError(t *testing.T) {
	resp := Dispatch("split 18 0")
	assertSingleError(t, resp)
}

func TestNegativeByReturnsError(t *testing.T) {
	resp := Dispatch("split 18 -6")
	assertSingleError(t, resp)
}

func TestInvalidPatternReturnsError(t *testing.T) {
	resp := Dispatch("18 z-a")
	assertSingleError(t, resp)
}

func TestPancInvalidPatternReturnsError(t *testing.T) {
	resp := Dispatch("panc 18 z-a")
	assertSingleError(t, resp)
}

func TestTrailingDashPatternReturnsError(t *testing.T) {
	resp := Dispatch("18 A-Z-")
	assertSingleError(t, resp)
}

func TestSplitPartialInvalidByUsesDefaults(t *testing.T) {
	// "abc" fails to parse as by → falls back to default by=6; length=18 is parsed.
	resp := Dispatch("split 18 abc")
	for _, it := range resp.Items {
		parts := strings.Split(it.Title, "-")
		if len(parts) != 3 {
			t.Fatalf("got %d parts, want 3", len(parts))
		}
		for _, p := range parts {
			if len([]rune(p)) != 6 {
				t.Errorf("part %q has length %d, want 6", p, len([]rune(p)))
			}
		}
	}
}

func TestItemsAreCopyable(t *testing.T) {
	resp := Dispatch("")
	for _, it := range resp.Items {
		if it.Valid == nil || !*it.Valid {
			t.Errorf("item %+v should be valid", it)
		}
		if it.Arg != it.Title {
			t.Errorf("arg %q should equal title %q", it.Arg, it.Title)
		}
	}
}

func TestSkipKnowledgeTrue(t *testing.T) {
	resp := Dispatch("")
	if !resp.SkipKnowledge {
		t.Error("expected SkipKnowledge to be true")
	}
}

func TestHelpShowsAllCommands(t *testing.T) {
	resp := Dispatch("help")
	if len(resp.Items) != len(helpCommands) {
		t.Fatalf("got %d items, want %d", len(resp.Items), len(helpCommands))
	}
}

func TestHelpAllItemsInvalid(t *testing.T) {
	resp := Dispatch("help")
	for _, it := range resp.Items {
		if it.Valid == nil || *it.Valid {
			t.Errorf("item %+v should be invalid", it)
		}
	}
}

func subtitles(resp scriptfilter.Response) string {
	var b strings.Builder
	for _, it := range resp.Items {
		b.WriteString(it.Subtitle)
		b.WriteString("\n")
	}
	return b.String()
}

func assertSingleError(t *testing.T, resp scriptfilter.Response) {
	t.Helper()
	if len(resp.Items) != 1 {
		t.Fatalf("got %d items, want 1", len(resp.Items))
	}
	if !strings.Contains(resp.Items[0].Title, "Error") {
		t.Errorf("title %q should contain 'Error'", resp.Items[0].Title)
	}
}
