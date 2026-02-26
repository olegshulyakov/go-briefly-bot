"""Module for capturing yt-dlp log output."""


class YtDlpCaptureLogger:
    """Captures yt-dlp log output for debugging."""

    def __init__(self) -> None:
        """Initialize the capture logger."""
        self.messages: list[str] = []

    def debug(self, msg: str) -> None:
        """Capture debug message."""
        self.messages.append(f"[debug] {msg}")

    def info(self, msg: str) -> None:
        """Capture info message."""
        self.messages.append(f"[info] {msg}")

    def warning(self, msg: str) -> None:
        """Capture warning message."""
        self.messages.append(f"[warning] {msg}")

    def error(self, msg: str) -> None:
        """Capture error message."""
        self.messages.append(f"[error] {msg}")
