# Build the manager binary
ARG VERSION=1.19-alpine

FROM golang:${VERSION} as builder
WORKDIR /workspace
ENV CGO_ENABLED=0
ENV GOOS=linux
COPY libs/go-libs libs/go-libs
COPY components/payments components/payments
COPY components/search components/search
COPY components/operator components/operator
COPY go.mod .
COPY go.sum .
RUN go mod download
RUN go install -v -installsuffix cgo -a std
RUN go mod vendor
RUN go build -v -a -o manager main.go

FROM golang:${VERSION} as reloader
RUN go install github.com/cosmtrek/air@latest

# # Use distroless as minimal base image to package the manager binary
# # Refer to https://github.com/GoogleContainerTools/distroless for more details
# FROM gcr.io/distroless/static:nonroot as release
# LABEL org.opencontainers.image.source=https://github.com/formancehq/operator
# WORKDIR /
# COPY --from=builder /workspace/manager /usr/bin/operator
# USER 65532:65532
# ENTRYPOINT ["/usr/bin/operator"]
