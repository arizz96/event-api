FROM golang:alpine

ENV REPO="github.com/arizz96/event-api"
WORKDIR /go/src/${REPO}

RUN apk add --no-cache \
		build-base \
		cyrus-sasl-dev \
		git \
		libressl \
		openssl-dev \
		yajl-dev \
		zlib-dev \
		bash

COPY . .

RUN git clone https://github.com/edenhill/librdkafka.git && \
		cd librdkafka && \
		git checkout '0.11.1.x' && \
		cp ../patch-librdkafka/src-cpp.Makefile src-cpp/Makefile && \
		cp ../patch-librdkafka/src.Makefile src/Makefile && \
		./configure --clean && \
		./configure --prefix /usr && \
		make && \
		make install

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
