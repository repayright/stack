FROM golang:1.19-alpine AS builder
WORKDIR /src
COPY components/ledger components/ledger
COPY libs libs
WORKDIR /src/components/ledger
ARG CGO_ENABLED=0
ARG APP_SHA
ARG VERSION
RUN --mount=type=cache,id=gobuild,target=/root/.cache/go-build --mount=type=cache,id=gomodcache,target=/go/pkg/mod go build -o ledger \
    -ldflags="-X $(cat go.mod |head -1|cut -d \  -f2)/cmd.Version=${VERSION} \
    -X $(cat go.mod |head -1|cut -d \  -f2)/cmd.BuildDate=$(date +%s) \
    -X $(cat go.mod |head -1|cut -d \  -f2)/cmd.Commit=${APP_SHA}" ./

FROM alpine:3.16
RUN apk update && apk add ca-certificates curl
COPY --from=builder /src/components/ledger/ledger /main
ENV OTEL_SERVICE_NAME ledger
ENTRYPOINT ["/main"]
CMD ["serve"]
