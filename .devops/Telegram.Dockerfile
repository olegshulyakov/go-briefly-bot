FROM python:3.12-alpine

WORKDIR /app

RUN apk add --no-cache ffmpeg

COPY pyproject.toml ./
RUN python -m pip install --no-cache-dir .

COPY src /app/src
COPY locales /app/locales

RUN mkdir -p /app/data
VOLUME ["/app/data"]

ENTRYPOINT ["python", "-m", "src.main"]
