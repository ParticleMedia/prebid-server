FROM ubuntu:20.04 AS build
RUN apt-get update && \
    apt-get -y upgrade && \
    apt-get install wget make git -y && \
    apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*
WORKDIR /tmp
RUN wget https://dl.google.com/go/go1.19.2.linux-amd64.tar.gz && \
    tar -xf go1.19.2.linux-amd64.tar.gz && \
    mv go /usr/local


# Install and configure GO
ENV GOROOT=/usr/local/go
ENV PATH=$GOROOT/bin:$PATH
ENV GOPROXY="https://proxy.golang.org"
RUN go install github.com/actgardner/gogen-avro/v10/cmd/...@latest
ENV PATH=/root/go/bin:$PATH
ENV CGO_ENABLED 0

# Build project
RUN mkdir -p /app/prebid-server/
WORKDIR /app/prebid-server/
COPY ./ ./
RUN go mod tidy
RUN go mod vendor
RUN make avro
ARG TEST="true"
RUN if [ "$TEST" != "false" ]; then ./validate.sh ; fi
RUN go build -mod=vendor -ldflags "-X github.com/prebid/prebid-server/version.Ver=`git describe --tags | sed 's/^v//'` -X github.com/prebid/prebid-server/version.Rev=`git rev-parse HEAD`" .

FROM ubuntu:20.04 AS release
LABEL maintainer="hans.hjort@xandr.com" 
WORKDIR /usr/local/bin/
COPY --from=build /app/prebid-server .
RUN chmod a+xr prebid-server
COPY static static/
COPY stored_requests/data stored_requests/data
RUN chmod -R a+r static/ stored_requests/data
RUN apt-get update && \
    apt-get install -y ca-certificates mtr && \
    apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*
RUN adduser prebid_user
USER prebid_user
EXPOSE 8000
EXPOSE 6060
ENTRYPOINT ["/usr/local/bin/prebid-server"]
CMD ["-v", "1", "-logtostderr"]
