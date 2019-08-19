FROM golang:buster

RUN apt-get update && apt-get install --no-install-recommends --no-install-suggests -y \
    libwebp-dev

RUN mkdir -p "${GOPATH}/src/github.com/yano3/okowa"
COPY . /go/src/github.com/yano3/okowa

RUN cd ${GOPATH}/src/github.com/yano3/okowa \
 && go get ./... \
 && go install github.com/yano3/okowa

FROM debian:buster-slim

RUN apt-get update && apt-get install --no-install-recommends --no-install-suggests -y \
    libwebp6 \
    ca-certificates \
 \
 && apt-get clean \
 && rm -rf /var/lib/apt/lists/* \
 \
 && mkdir -p "/go/bin"

COPY --from=0 /go/bin/okowa /go/bin

ENV GOPATH /go
ENV PATH $GOPATH/bin:$PATH

EXPOSE 8080

CMD ["okowa"]
