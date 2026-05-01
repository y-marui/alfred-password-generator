"""Application orchestrator.

Wires together the Router and Command handlers.
This is the single entry point called by scripts/entry.py.

Commands:
  passgen [length] [pattern]             — basic password (default)
  panc [split] [length] [by] [pattern]   — with punctuation
  split [length] [by] [pattern]          — split by groups
  config [reset]                         — view or reset configuration
  help                                   — show available commands
"""

from __future__ import annotations

from alfred.router import Router
from app.commands import config_cmd, help_cmd, passgen_cmd

router = Router(default="passgen")
router.register("passgen")(passgen_cmd.handle_basic)
router.register("panc")(passgen_cmd.handle_panc)
router.register("split")(passgen_cmd.handle_split)
router.register("config")(config_cmd.handle)
router.register("help")(help_cmd.handle)


def run(query: str) -> None:
    """Main application entry point.

    Args:
        query: Raw query string from Alfred (e.g. "panc split 18 6").
    """
    router.dispatch(query)
