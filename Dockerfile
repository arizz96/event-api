FROM golang:alpine

ARG LIBRDKAFKA_VERSION="0.11.1-r1"

ENV REPO="github.com/arizz96/event-api"
WORKDIR /go/src/${REPO}

RUN apk add --no-cache \
		build-base \
		cyrus-sasl-dev \
		git \
		librdkafka-dev=${LIBRDKAFKA_VERSION} \
		libressl \
		openssl-dev \
		yajl-dev \
		zlib-dev

COPY . .
RUN cd vendor/github.com/gin-contrib && \
		git clone https://github.com/arizz96/cors.git && \
		cd cors && \
		git checkout 'feature/regexp-origin-match' && \
		go build
RUN make build

ENV GIN_MODE=release
EXPOSE 8080 8081

# Use ENTRYPOINT in production images
CMD ["./event-api", "start"]
