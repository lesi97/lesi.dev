# syntax=docker/dockerfile:1.7

ARG GO_VERSION=1.24.5
ARG ALPINE_VERSION=3.20

FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-alpine AS go-build
RUN apk add --no-cache git

WORKDIR /src

COPY apps/api/go.mod apps/api/go.sum ./apps/api/

RUN --mount=type=cache,id=cache-go-mod,target=/go/pkg/mod cd apps/api && go mod download

COPY apps/api ./apps/api

ARG TARGETARCH
RUN --mount=type=cache,id=cache-go-mod,target=/go/pkg/mod cd apps/api && CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH go build -trimpath -ldflags="-s -w" -o /out/server ./cmd


FROM alpine:${ALPINE_VERSION} AS final
RUN apk add --no-cache ca-certificates tzdata

ARG UID=10001
RUN adduser --disabled-password --gecos "" --home "/nonexistent" --shell "/sbin/nologin" --no-create-home --uid "${UID}" appuser

USER appuser
WORKDIR /app

COPY --from=go-build /out/server /app/server

ENV GO_ENV=production

EXPOSE 8080
ENTRYPOINT ["/app/server"]
