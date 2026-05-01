from __future__ import annotations

import random
import string
from collections.abc import Iterator

_PUNCT = "!@#^&*"


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
            _flag, flag = False, _flag
            _last, last = c, _last
            if flag:
                if last == "":
                    raise ValueError(f"-{c} is not a valid pattern")
                if last in string.ascii_lowercase and c in string.ascii_lowercase:
                    i_s = string.ascii_lowercase.find(last)
                    i_e = string.ascii_lowercase.find(c)
                    if i_s > i_e:
                        raise ValueError(f"{last}-{c} is not a valid range")
                    for _c in string.ascii_lowercase[i_s + 1 : i_e + 1]:
                        yield _c
                elif last in string.ascii_uppercase and c in string.ascii_uppercase:
                    i_s = string.ascii_uppercase.find(last)
                    i_e = string.ascii_uppercase.find(c)
                    if i_s > i_e:
                        raise ValueError(f"{last}-{c} is not a valid range")
                    for _c in string.ascii_uppercase[i_s + 1 : i_e + 1]:
                        yield _c
                elif last in string.digits and c in string.digits:
                    i_s = string.digits.find(last)
                    i_e = string.digits.find(c)
                    if i_s > i_e:
                        raise ValueError(f"{last}-{c} is not a valid range")
                    for _c in string.digits[i_s + 1 : i_e + 1]:
                        yield _c
                elif last in _PUNCT and c in _PUNCT:
                    i_s = _PUNCT.find(last)
                    i_e = _PUNCT.find(c)
                    if i_s > i_e:
                        raise ValueError(f"{last}-{c} is not a valid range")
                    for _c in _PUNCT[i_s + 1 : i_e + 1]:
                        yield _c
                else:
                    raise ValueError(f"{last}-{c} is not a valid range")
            else:
                yield c


def generate(pattern: str, length: int) -> str:
    """Generate a random password of given length from characters in pattern."""
    arr = list(_loader(pattern))
    if not arr:
        raise ValueError(f"Pattern '{pattern}' produces no characters")
    return "".join(random.choices(arr, k=length))


def generate_split(pattern: str, length: int, by: int) -> str:
    """Generate a password split into groups of `by` characters, joined by hyphens."""
    if length < by:
        raise ValueError(f"length must be >= by: length={length}, by={by}")
    if length % by != 0:
        raise ValueError(f"length must be a multiple of by: length={length}, by={by}")
    arr = list(_loader(pattern))
    if not arr:
        raise ValueError(f"Pattern '{pattern}' produces no characters")
    return "-".join(["".join(random.choices(arr, k=by)) for _ in range(length // by)])
