"""Tests for the passgen service."""

from __future__ import annotations

import pytest

from app.services import passgen_service


class TestGenerate:
    def test_returns_correct_length(self):
        pwd = passgen_service.generate("A-Za-z0-9", 18)
        assert len(pwd) == 18

    def test_characters_in_pattern(self):
        import string

        pwd = passgen_service.generate("A-Za-z0-9", 100)
        allowed = set(string.ascii_letters + string.digits)
        assert all(c in allowed for c in pwd)

    def test_explicit_chars(self):
        pwd = passgen_service.generate("abc", 10)
        assert all(c in "abc" for c in pwd)

    def test_custom_length(self):
        assert len(passgen_service.generate("A-Z", 32)) == 32

    def test_punctuation_pattern(self):
        pwd = passgen_service.generate("!-*", 20)
        assert all(c in "!@#^&*" for c in pwd)

    def test_invalid_range_raises(self):
        with pytest.raises(ValueError):
            passgen_service.generate("z-a", 10)

    def test_double_dash_raises(self):
        with pytest.raises(ValueError):
            passgen_service.generate("a--z", 10)


class TestGenerateSplit:
    def test_correct_format(self):
        pwd = passgen_service.generate_split("A-Za-z0-9", 18, 6)
        parts = pwd.split("-")
        assert len(parts) == 3
        assert all(len(p) == 6 for p in parts)

    def test_custom_by(self):
        pwd = passgen_service.generate_split("A-Z", 12, 4)
        parts = pwd.split("-")
        assert len(parts) == 3
        assert all(len(p) == 4 for p in parts)

    def test_length_not_multiple_of_by_raises(self):
        with pytest.raises(ValueError, match="multiple"):
            passgen_service.generate_split("A-Z", 18, 7)

    def test_length_less_than_by_raises(self):
        with pytest.raises(ValueError, match=">="):
            passgen_service.generate_split("A-Z", 3, 6)

    def test_characters_in_pattern(self):
        import string

        pwd = passgen_service.generate_split("A-Za-z0-9", 18, 6)
        chars = pwd.replace("-", "")
        allowed = set(string.ascii_letters + string.digits)
        assert all(c in allowed for c in chars)
