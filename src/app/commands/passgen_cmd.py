"""passgen command - generate passwords.

Usage in Alfred:
  passgen [length]                      — overview of all patterns
  passgen [length] [pattern]            — custom charset (single result)
  passgen panc [length] [pattern]       — with punctuation
  passgen split [length] [by] [pattern] — split into groups
  passgen panc split [length] [by] [pattern]

Default patterns:
  basic / split:       A-Za-z0-9
  panc / panc split:   A-Za-z0-9!-*  (punctuation: !@#^&*)
"""

from __future__ import annotations

from collections.abc import Callable

from alfred.logger import get_logger
from alfred.response import error_item, item, output
from app.services import passgen_service

log = get_logger(__name__)

_DEFAULT_LENGTH = 18
_DEFAULT_BY = 6
_PATTERN_BASIC = "A-Za-z0-9"
_PATTERN_PANC = "A-Za-z0-9!-*"
_NUM_SUGGESTIONS = 5

_OVERVIEW = [
    ("basic", _PATTERN_BASIC, False),
    ("panc", _PATTERN_PANC, False),
    ("split", _PATTERN_BASIC, True),
    ("panc split", _PATTERN_PANC, True),
]


def handle_basic(args: str) -> None:
    """passgen [length] [pattern] — overview or custom password."""
    log.debug("passgen command: args=%r", args)
    parts = args.strip().split(None, 1)
    if not parts:
        _show_overview(_DEFAULT_LENGTH)
        return
    try:
        length = int(parts[0])
    except ValueError:
        _run_basic(args, _PATTERN_BASIC, count=1)
        return
    if len(parts) == 1:
        _show_overview(length)
    else:
        _run_basic(args, _PATTERN_BASIC, count=1)


def handle_panc(args: str) -> None:
    """passgen panc [...] — with punctuation; optionally split."""
    log.debug("passgen panc command: args=%r", args)
    parts = args.strip().split(None, 1)
    if parts and parts[0].lower() == "split":
        _run_split(parts[1] if len(parts) > 1 else "", _PATTERN_PANC)
    else:
        _run_basic(args, _PATTERN_PANC)


def handle_split(args: str) -> None:
    """passgen split [...] — split password without punctuation."""
    log.debug("passgen split command: args=%r", args)
    _run_split(args, _PATTERN_BASIC)


def _show_overview(length: int) -> None:
    items = []
    for i, (name, pattern, do_split) in enumerate(_OVERVIEW):
        try:
            if do_split:
                pwd = passgen_service.generate_split(pattern, length, _DEFAULT_BY)
                subtitle = f"{name} · {length} chars in groups of {_DEFAULT_BY}"
            else:
                pwd = passgen_service.generate(pattern, length)
                subtitle = f"{name} · {length} chars"
        except ValueError:
            continue
        items.append(item(title=pwd, subtitle=subtitle, arg=pwd, uid=f"passgen-{i}"))
    if not items:
        output([error_item(f"No valid patterns for length {length}")])
        return
    output(items, skip_knowledge=True)


def _run_basic(args: str, default_pattern: str, count: int = _NUM_SUGGESTIONS) -> None:
    length, pattern = _parse_basic(args, default_pattern)
    subtitle = f"{length} chars, pattern: {pattern}"
    _output(lambda: passgen_service.generate(pattern, length), subtitle, count=count)


def _run_split(args: str, default_pattern: str) -> None:
    length, by, pattern = _parse_split(args, default_pattern)
    subtitle = f"{length} chars in groups of {by}, pattern: {pattern}"
    _output(lambda: passgen_service.generate_split(pattern, length, by), subtitle)


def _parse_basic(args: str, default_pattern: str) -> tuple[int, str]:
    parts = args.strip().split(None, 1)
    if not parts:
        return _DEFAULT_LENGTH, default_pattern
    try:
        length = int(parts[0])
    except ValueError:
        log.debug("_parse_basic: non-integer %r, using defaults", parts[0])
        return _DEFAULT_LENGTH, default_pattern
    pattern = parts[1].strip() if len(parts) > 1 else default_pattern
    return length, pattern


def _parse_split(args: str, default_pattern: str) -> tuple[int, int, str]:
    parts = args.strip().split(None, 2)
    length = _DEFAULT_LENGTH
    by = _DEFAULT_BY
    pattern = default_pattern

    if len(parts) >= 1:
        try:
            length = int(parts[0])
        except ValueError:
            log.debug("_parse_split: non-integer length %r, using default", parts[0])

    if len(parts) >= 2:
        try:
            by = int(parts[1])
        except ValueError:
            log.debug("_parse_split: non-integer by %r, using default", parts[1])

    if len(parts) >= 3:
        pattern = parts[2].strip()

    return length, by, pattern


def _output(gen_fn: Callable[[], str], subtitle: str, count: int = _NUM_SUGGESTIONS) -> None:
    items = []
    for i in range(count):
        try:
            pwd = gen_fn()
        except ValueError as e:
            log.error("password generation failed: %s", e)
            output([error_item(str(e))])
            return
        items.append(
            item(
                title=pwd,
                subtitle=subtitle,
                arg=pwd,
                uid=f"passgen-{i}",
            )
        )
    output(items, skip_knowledge=True)
