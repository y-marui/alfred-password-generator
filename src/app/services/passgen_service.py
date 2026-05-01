from __future__ import annotations

import secrets
import string
from collections.abc import Iterator

_PUNCT = "!@#^&*"
_RANGES = (string.ascii_lowercase, string.ascii_uppercase, string.digits, _PUNCT)


def _loader(pattern: str) -> Iterator[str]:
    """Expand a character range pattern (e.g. 'A-Za-z0-9!-*') into individual characters."""
    _last = ""
    _flag = False
    for c in pattern:
        if c == "-":
            if _flag:
                raise ValueError("-- is not a valid pattern")
            _flag = True
        else:
            flag, _flag = _flag, False
            last, _last = _last, c
            if flag:
                if last == "":
                    raise ValueError(f"-{c} is not a valid pattern")
                source = next((s for s in _RANGES if last in s and c in s), None)
                if source is None:
                    raise ValueError(f"{last}-{c} is not a valid range")
                i_s, i_e = source.find(last), source.find(c)
                if i_s > i_e:
                    raise ValueError(f"{last}-{c} is not a valid range")
                yield from source[i_s + 1 : i_e + 1]
            else:
                yield c
    if _flag:
        raise ValueError("pattern ends with '-'; trailing dash is not valid")


def generate(pattern: str, length: int) -> str:
    """Generate a cryptographically secure random password of given length from pattern."""
    if length <= 0:
        raise ValueError(f"length must be a positive integer, got {length}")
    arr = list(_loader(pattern))
    if not arr:
        raise ValueError(f"Pattern '{pattern}' produces no characters")
    return "".join(secrets.choice(arr) for _ in range(length))


def generate_split(pattern: str, length: int, by: int) -> str:
    """Generate a secure password split into groups of `by` characters, joined by hyphens."""
    if by <= 0:
        raise ValueError(f"by must be a positive integer, got {by}")
    if length <= 0:
        raise ValueError(f"length must be a positive integer, got {length}")
    if length < by:
        raise ValueError(f"length must be >= by: length={length}, by={by}")
    if length % by != 0:
        raise ValueError(f"length must be a multiple of by: length={length}, by={by}")
    arr = list(_loader(pattern))
    if not arr:
        raise ValueError(f"Pattern '{pattern}' produces no characters")
    return "-".join("".join(secrets.choice(arr) for _ in range(by)) for _ in range(length // by))
