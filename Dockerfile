FROM golang:1.18.0-bullseye
WORKDIR /build
COPY ./ ./
WORKDIR /build
RUN make

FROM busybox:uclibc
WORKDIR /runtime
RUN mkdir -p /config
COPY --from=0 /build/arangoinit .
CMD /runtime/arangoinit -config /config/initarango.yaml -user $ARANGODB_ROOTUSER -pass $ARANGODB_ROOTPASS -endpoint $ARANGODB_ENDPOINT -timeout $ARANGODB_TIMEOUT -retry $ARANGODB_RETRY
