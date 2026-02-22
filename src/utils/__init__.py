"""
Utilities module.

Provides helper functions for text processing and formatting.
"""

from .compression import CompressionMethod, compress, decompress
from .markdown import markdown_to_telegram_html
from .text import to_lexical_chunks

__all__ = [
    "CompressionMethod",
    "compress",
    "decompress",
    "markdown_to_telegram_html",
    "to_lexical_chunks",
]
