FROM golang:1.17 AS build
WORKDIR /build

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./
ENV CGO_ENABLED=0
RUN go build  \
  -ldflags " \
    -s -w \
    -X github.com/dpb587/hubitat-cli/cmd/cmdflags.VersionCommit=$( git rev-parse HEAD ) \
    -X github.com/dpb587/hubitat-cli/cmd/cmdflags.VersionBuilt=$( date -u +%Y-%m-%dT%H:%M:%S+00:00 ) \
  " \
  -o /build/hubitat-cli

FROM gcr.io/distroless/static
WORKDIR /
COPY --from=build /build/hubitat-cli /bin/hubitat-cli
USER nonroot:nonroot
ENTRYPOINT ["/bin/hubitat-cli"]
CMD ["--help"]
