"""help command - display available commands.

Usage in Alfred:  passgen help
"""

from __future__ import annotations

from alfred.response import item, output

_COMMANDS = [
    ("passgen [length] [pattern]", "Generate password (default: 18 chars, A-Za-z0-9)", ""),
    ("passgen panc [length] [pattern]", "Generate with punctuation (!@#^&*)", "panc "),
    ("passgen split [length] [by] [pattern]", "Generate split password (default: 18/6)", "split "),
    ("passgen panc split [length] [by] [pattern]", "Split with punctuation", "panc split "),
    ("passgen config", "View or reset configuration", "config"),
    ("passgen help", "Show this help", "help"),
]


def handle(_args: str) -> None:
    """Display all available commands."""
    output(
        [
            item(
                title=cmd,
                subtitle=desc,
                arg="",
                uid=f"help-{i}",
                valid=False,
                autocomplete=autocomplete,
            )
            for i, (cmd, desc, autocomplete) in enumerate(_COMMANDS)
        ]
    )
