FROM golang:1.23.2-bookworm AS builder
WORKDIR /app
ARG TAG=dev
ARG COMMITID=0000000
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o target/turbo-mailer \
    -ldflags "-s -w -X turbo-mailer-server/version.Version=${TAG} -X turbo-mailer-server/version.CommitID=$(expr substr "$COMMITID" 1 7)" main.go


FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y ca-certificates tzdata && \
    rm -rf /var/lib/apt/lists/*
RUN addgroup --gid 1001 appuser && \
    adduser  --uid 1001 --ingroup appuser --disabled-password --gecos "" appuser
RUN mkdir -p /app && chown -R appuser:appuser /app
USER appuser
ENV TZ=Asia/Shanghai
WORKDIR /app
COPY --from=builder --chown=appuser:appuser /app/target/ /app/
COPY --from=builder --chown=appuser:appuser /app/docs/ /app/docs/
CMD ["./turbo-mailer", "serve"]
