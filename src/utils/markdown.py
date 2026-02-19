import re

import markdown
from bs4 import BeautifulSoup


def markdown_to_telegram_html(md_text):
    raw_html = markdown.markdown(md_text, extensions=["extra", "codehilite"])
    soup = BeautifulSoup(raw_html, "html.parser")
    allowed_tags = ["b", "strong", "i", "em", "u", "ins", "s", "strike", "del", "a", "code", "pre"]

    # replace unsupportted tags
    for tag in soup.find_all(True):
        if tag.name not in allowed_tags:
            # Replace headers with bold text
            if tag.name in ["h1", "h2", "h3", "h4", "h5", "h6"]:
                tag.name = "strong"
                tag.insert_before("\n")
            elif tag.name == "li":
                tag.insert_before("- ")
                tag.unwrap()
            else:
                tag.unwrap()

    result = str(soup).strip()
    # Replace multiple consecutive newlines with a single one
    result = re.sub(r"\n{3,}", "\n\n", result)
    return result
