ARG build_tag

FROM golang:1.14-alpine AS build_base

ENV GOLANG_VERSION 1.14
ENV GO111MODULE=on

RUN apk add bash ca-certificates git gcc g++ libc-dev

RUN mkdir /src
WORKDIR /src

COPY go.mod .
COPY go.sum .

# Get dependancies - will also be cached if we won't change mod/sum
RUN go mod download
# COPY the source code as the last step

FROM build_base AS server_builder

COPY . .

ENV GOPROXY=direct

# Build the binary
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -a -trimpath -o /go/bin/server -race -ldflags "-X 'main.build=${build_tag}'"

FROM alpine
RUN apk add ca-certificates

COPY template/* /go/bin/template/
COPY --from=server_builder /go/bin/server /go/bin/server

EXPOSE 8080
ENTRYPOINT /go/bin/server --port=8080 --template="/go/bin/template"
