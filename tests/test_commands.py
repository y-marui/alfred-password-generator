"""Tests for command handlers."""

from __future__ import annotations

import json

from app.commands import config_cmd, help_cmd, passgen_cmd


class TestPassgenCommand:
    def _items(self, capsys) -> list:
        return json.loads(capsys.readouterr().out)["items"]

    def test_empty_query_shows_overview(self, capsys):
        passgen_cmd.handle_basic("")
        items = self._items(capsys)
        assert len(items) == len(passgen_cmd._OVERVIEW)
        subtitles = [it["subtitle"] for it in items]
        assert any("basic" in s for s in subtitles)
        assert any("panc" in s for s in subtitles)
        assert any("split" in s for s in subtitles)

    def test_length_only_shows_overview(self, capsys):
        passgen_cmd.handle_basic("24")
        items = self._items(capsys)
        assert len(items) == len(passgen_cmd._OVERVIEW)
        for it in items:
            assert "24" in it["subtitle"]

    def test_non_divisible_length_skips_split(self, capsys):
        # 20 is not divisible by 6, so split variants are skipped
        passgen_cmd.handle_basic("20")
        items = self._items(capsys)
        subtitles = [it["subtitle"] for it in items]
        assert all("split" not in s for s in subtitles)

    def test_custom_pattern_returns_single_result(self, capsys):
        passgen_cmd.handle_basic("18 A-Z")
        items = self._items(capsys)
        assert len(items) == 1
        chars = items[0]["title"]
        assert all(c.isupper() for c in chars)
        assert len(chars) == 18

    def test_panc_includes_punctuation_by_default(self, capsys):
        passgen_cmd.handle_panc("")
        items = self._items(capsys)
        assert len(items) == passgen_cmd._NUM_SUGGESTIONS
        all_chars = "".join(it["title"] for it in items)
        punct = set("!@#^&*")
        # 5 * 18 = 90 chars from a 68-char charset; probability of zero punct is negligible
        assert any(c in punct for c in all_chars)

    def test_panc_split(self, capsys):
        passgen_cmd.handle_panc("split 18 6")
        items = self._items(capsys)
        for it in items:
            parts = it["title"].split("-")
            assert len(parts) == 3
            assert all(len(p) == 6 for p in parts)

    def test_panc_split_no_args(self, capsys):
        passgen_cmd.handle_panc("split")
        items = self._items(capsys)
        assert len(items) == passgen_cmd._NUM_SUGGESTIONS
        for it in items:
            parts = it["title"].split("-")
            assert len(parts) == passgen_cmd._DEFAULT_LENGTH // passgen_cmd._DEFAULT_BY

    def test_split_command(self, capsys):
        passgen_cmd.handle_split("18 6")
        items = self._items(capsys)
        for it in items:
            parts = it["title"].split("-")
            assert len(parts) == 3
            assert all(len(p) == 6 for p in parts)

    def test_split_default_args(self, capsys):
        passgen_cmd.handle_split("")
        items = self._items(capsys)
        assert len(items) == passgen_cmd._NUM_SUGGESTIONS

    def test_split_custom_pattern(self, capsys):
        passgen_cmd.handle_split("12 4 A-Z")
        items = self._items(capsys)
        for it in items:
            chars = it["title"].replace("-", "")
            assert all(c.isupper() for c in chars)
            assert len(chars) == 12

    def test_invalid_split_ratio_returns_error(self, capsys):
        passgen_cmd.handle_split("18 7")
        items = self._items(capsys)
        assert len(items) == 1
        assert "Error" in items[0]["title"]

    def test_zero_by_returns_error(self, capsys):
        passgen_cmd.handle_split("18 0")
        items = self._items(capsys)
        assert len(items) == 1
        assert "Error" in items[0]["title"]

    def test_negative_by_returns_error(self, capsys):
        passgen_cmd.handle_split("18 -6")
        items = self._items(capsys)
        assert len(items) == 1
        assert "Error" in items[0]["title"]

    def test_invalid_pattern_returns_error(self, capsys):
        passgen_cmd.handle_basic("18 z-a")
        items = self._items(capsys)
        assert len(items) == 1
        assert "Error" in items[0]["title"]

    def test_panc_invalid_pattern_returns_error(self, capsys):
        passgen_cmd.handle_panc("18 z-a")
        items = self._items(capsys)
        assert len(items) == 1
        assert "Error" in items[0]["title"]

    def test_trailing_dash_pattern_returns_error(self, capsys):
        passgen_cmd.handle_basic("18 A-Z-")
        items = self._items(capsys)
        assert len(items) == 1
        assert "Error" in items[0]["title"]

    def test_split_partial_invalid_by_uses_defaults(self, capsys):
        # "abc" fails int() for by → falls back to default by=6; length=18 is parsed
        passgen_cmd.handle_split("18 abc")
        items = self._items(capsys)
        for it in items:
            parts = it["title"].split("-")
            assert len(parts) == 3
            assert all(len(p) == 6 for p in parts)

    def test_items_are_copyable(self, capsys):
        passgen_cmd.handle_basic("")
        items = self._items(capsys)
        for it in items:
            assert it["valid"] is True
            assert it["arg"] == it["title"]

    def test_skip_knowledge_true(self, capsys):
        passgen_cmd.handle_basic("")
        data = json.loads(capsys.readouterr().out)
        assert data.get("skipknowledge") is True


class TestConfigCommand:
    def test_empty_config_shows_no_settings(self, capsys):
        config_cmd.handle("")
        data = json.loads(capsys.readouterr().out)
        titles = [it["title"] for it in data["items"]]
        assert any("No settings" in t for t in titles)

    def test_reset_clears_config(self, capsys):
        config_cmd._config.set("key", "value")
        config_cmd.handle("reset")
        data = json.loads(capsys.readouterr().out)
        assert "reset" in data["items"][0]["title"].lower()
        assert config_cmd._config.all() == {}

    def test_shows_existing_settings(self, capsys):
        config_cmd._config.set("log_level", "DEBUG")
        config_cmd.handle("")
        data = json.loads(capsys.readouterr().out)
        titles = [it["title"] for it in data["items"]]
        assert any("log_level" in t for t in titles)


class TestHelpCommand:
    def test_shows_all_commands(self, capsys):
        help_cmd.handle("")
        data = json.loads(capsys.readouterr().out)
        assert len(data["items"]) == len(help_cmd._COMMANDS)

    def test_all_items_invalid(self, capsys):
        help_cmd.handle("")
        data = json.loads(capsys.readouterr().out)
        assert all(not it["valid"] for it in data["items"])
