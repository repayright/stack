VERSION 0.7

ARG --global ALPINE_VERSION=3.17
ARG --global VERSION=develop
ARG --global GOOS=linux
ARG --global GOARCH=arm64
ARG --global BUILD_DATE=-
ARG --global REPOSITORY=ghcr.io
ARG --global DOCKER_VERSION=23.0.1
ARG --global GOLANG_VERSION=1.20
ARG --global NODE_VERSION=19
ARG --global GOLANGCI_LINT_VERSION=v1.51.2
ARG --global OPENAPI_GENERATOR_VERSION=v6.4.0

FROM alpine:${ALPINE_VERSION}
WORKDIR /src

DOWNLOAD_DEPENDENCIES:
    COMMAND
    ARG LOCATION
    ENV ACTUAL_LOCATION=$(pwd)
    WORKDIR /src/$LOCATION
    COPY (+sources/$LOCATION/go.* --SOURCE_PATH=$LOCATION) .
    RUN go mod download
    WORKDIR $ACTUAL_LOCATION

LOAD_SOURCES:
    COMMAND
    ARG LOCATION
    COPY (+sources/$LOCATION/ --SOURCE_PATH=$LOCATION) /src/$LOCATION/

LOAD_SOURCE_FILE:
    COMMAND
    ARG LOCATION
    COPY (+sources/$LOCATION --SOURCE_PATH=$LOCATION) /src/$LOCATION

SAVE_IMAGE:
    COMMAND
    ARG COMPONENT
    SAVE IMAGE ${REPOSITORY}/formancehq/${COMPONENT}:${VERSION}

GO_STD_BUILD:
    COMMAND
    ARG --required COMPONENT
    ARG GOOS=
    ARG GOARCH=
    ARG SEGMENT_WRITE_KEY=
    ARG EARTHLY_GIT_HASH
    RUN GOOS=$GOOS GOARCH=$GOARCH go build -o $COMPONENT \
        -ldflags="-X github.com/formancehq/$COMPONENT/cmd.Version=${VERSION} \
        -X github.com/formancehq/$COMPONENT/cmd.BuildDate=$BUILD_DATE \
        -X github.com/formancehq/$COMPONENT/cmd.Commit=$EARTHLY_GIT_HASH \
        -X github.com/formancehq/$COMPONENT/cmd.DefaultSegmentWriteKey=$SEGMENT_WRITE_KEY" ./
    SAVE ARTIFACT ./$COMPONENT

GO_STD_TEST_COMMAND:
    COMMAND
    RUN --mount=type=cache,target=/root/.cache/go-build \
        go test -coverpkg ./... -coverprofile coverage.out -covermode atomic ./...
    SAVE ARTIFACT ./coverage.out AS LOCAL ./coverage.out

# Share common images to avoid version drift
benthos:
    FROM jeffail/benthos:4.11

postgres:
    FROM postgres:15-alpine

opensearch:
    FROM opensearchproject/opensearch:2.2.1

redpanda:
    FROM docker.redpanda.com/vectorized/redpanda:v22.2.2

nats:
    FROM nats

base-build-image:
    FROM docker:23.0.1-dind-alpine3.17
    RUN apk update && apk add jq go nodejs make bash curl
    RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b /bin ${GOLANGCI_LINT_VERSION}
    ENV CGO_ENABLED=0
    ENV GOPROXY http://172.20.0.13:8080
    ENV GOSUMDB=off

    WORKDIR /src

base-final-image:
    RUN apk update && apk add ca-certificates curl
    ARG COMPONENT
    ENV OTEL_SERVICE_NAME $COMPONENT

sources:
    ARG --required SOURCE_PATH
    COPY $SOURCE_PATH /src/$SOURCE_PATH
    SAVE ARTIFACT /src/$SOURCE_PATH $SOURCE_PATH

build-component-binary:
    ARG --required COMPONENT

    FROM ./components/$COMPONENT+sources
    ARG GOOS=linux
    ARG GOARCH=arm64
    ARG SEGMENT_WRITE_KEY=
    ARG EARTHLY_GIT_HASH
    RUN GOOS=$GOOS GOARCH=$GOARCH go build -o $COMPONENT \
        -ldflags="-X github.com/formancehq/$COMPONENT/cmd.Version=${VERSION} \
        -X github.com/formancehq/$COMPONENT/cmd.BuildDate=$BUILD_DATE \
        -X github.com/formancehq/$COMPONENT/cmd.Commit=$EARTHLY_GIT_HASH \
        -X github.com/formancehq/$COMPONENT/cmd.DefaultSegmentWriteKey=$SEGMENT_WRITE_KEY" ./
    SAVE ARTIFACT ./$COMPONENT AS LOCAL ./bin/$COMPONENT-$GOOS-$GOARCH

build-image:
    ARG --required COMPONENT
    BUILD ./components/$COMPONENT+build-image

tests:
    ARG --required LOCATION
    BUILD ./$LOCATION+tests

lint:
    ARG --required LOCATION
    FROM ./$LOCATION+sources
    COPY .golangci.yml /tmp/.golangci.yml
    RUN --mount=type=cache,target=/root/.cache/go-build \
        golangci-lint -v run --fix --timeout=10m --verbose -c /tmp/.golangci.yml
    SAVE ARTIFACT ./* AS LOCAL ./$LOCATION/

build-component-binaries:
    LOCALLY
    ARG GOOS=
    ARG GOARCH=
    FOR component IN $(ls components)
        BUILD +build-component-binary --COMPONENT=$component --GOOS=$GOOS --GOARCH=$GOARCH
    END

build-component-binaries-on-all-architectures:
    LOCALLY
    BUILD +build-component-binaries --GOOS=linux --GOARCH=amd64
    BUILD +build-component-binaries --GOOS=linux --GOARCH=arm64
    BUILD +build-component-binaries --GOOS=darwin --GOARCH=amd64
    BUILD +build-component-binaries --GOOS=darwin --GOARCH=arm64
    BUILD +build-component-binaries --GOOS=windows --GOARCH=amd64

build-all-images:
    LOCALLY
    FOR component IN $(ls components)
        BUILD +build-image --COMPONENT=$component
    END

tests-all:
    LOCALLY
    BUILD +tests --LOCATION=libs/go-libs
    FOR component IN $(ls components)
        BUILD +tests --LOCATION=components/$component
    END

lint-all:
    LOCALLY
    BUILD +lint --LOCATION=libs/go-libs
    FOR component IN $(ls components)
        BUILD +lint --LOCATION=components/$component
    END

sdk-generation:
    FROM node:${NODE_VERSION}-alpine${ALPINE_VERSION}
    WORKDIR /src/openapi
    COPY openapi/package.* .
    RUN npm install
    SAVE ARTIFACT ./package.*
    SAVE ARTIFACT ./node_modules

generate-openapi-spec:
    FROM node:${NODE_VERSION}-alpine${ALPINE_VERSION}
    WORKDIR /src/openapi
    COPY openapi/package.* .
    RUN npm install
    COPY openapi /src/openapi
    RUN npm run build
    RUN sed -i -e "s/SDK_VERSION/$VERSION/g" build/generate.json
    SAVE ARTIFACT ./build/generate.json

generate-sdk:
    ARG --required LANG
    FROM openapitools/openapi-generator-cli:v6.4.0
    COPY (+generate-openapi-spec/generate.json) ./openapi/build/generate.json
    RUN docker-entrypoint.sh generate \
        -i ./openapi/build/generate.json \
        -g $(echo $LANG | cut -d - -f1) \
        -c ./openapi/configs/$LANG.yaml \
        -o ./sdks/$LANG \
        --git-user-id=formancehq \
        --git-repo-id=formance-sdk-$LANG \
        -p packageVersion=$VERSION \
        -p apiVersion=$VERSION
    SAVE ARTIFACT ./sdks/$LANG AS LOCAL ./sdks/$LANG

generate-all-sdks:
    LOCALLY
    FOR lang IN $(ls openapi/configs)
        BUILD +generate-sdk --LANG=$lang
    END

tests-integrations:
    FROM +base-build-image
    DO +LOAD_SOURCES --LOCATION=sdks/go
    DO +LOAD_SOURCES --LOCATION=libs/go-libs
    DO +LOAD_SOURCES --LOCATION=libs/events
    DO +LOAD_SOURCES --LOCATION=components
    DO +LOAD_SOURCE_FILE --LOCATION=go.mod
    DO +LOAD_SOURCE_FILE --LOCATION=go.sum
    RUN go mod download
    DO +DOWNLOAD_DEPENDENCIES --LOCATION=tests/integration
    DO +LOAD_SOURCES --LOCATION=tests/integration

    WORKDIR /src/tests/integration
    # Keep here and don't use a fixed version. Ginkgo is using go.mod version.
    RUN go install github.com/onsi/ginkgo/v2/ginkgo
    RUN cp $(go env GOPATH)/bin/ginkgo /bin

    ARG STRIPE_API_KEY
    WITH DOCKER --compose docker-compose.yml \
        --allow-privileged \
        --load jeffail/benthos:4.11=+benthos \
        --load postgres:15-alpine=+postgres \
        --load opensearchproject/opensearch:2.2.1=+opensearch \
        --load nats=+nats
        RUN --mount=type=cache,target=/root/.cache/go-build \
            export DOCKER_HOSTNAME=$(ip addr show eth0 | grep brd | grep inet | cut -d \  -f 6 | cut -d / -f 1) && \
            echo "Found gateway IP: $DOCKER_HOSTNAME" && \
            ginkgo -v -p --cover --coverpkg github.com/formancehq/... --coverprofile coverage.out ./suite
    END
    SAVE ARTIFACT ./coverage.out AS LOCAL ./coverage.out

all:
    WAIT
        BUILD +generate-all-sdks
        BUILD +lint-all
        BUILD +tests-all
        BUILD +tests-integrations
    END
    BUILD +build-component-binaries-on-all-architectures
    BUILD +build-all-images
