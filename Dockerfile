FROM gcr.io/distroless/static
WORKDIR /tmp
ADD tmp/build/hubitat-cli-*-linux-amd64 /bin/hubitat-cli
USER nonroot:nonroot
ENTRYPOINT ["/bin/hubitat-cli"]
CMD ["--help"]
