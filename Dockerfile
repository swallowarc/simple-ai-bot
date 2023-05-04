# Build Container
FROM golang:1.20 as builder
WORKDIR /go/src/github.com/swallowarc/simple-line-ai-bot
COPY . .

# Set Environment Variable
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
ARG GITHUB_KEY

# Build
RUN make

# runtime image
FROM alpine:3.17.3
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/src/github.com/swallowarc/simple-line-ai-bot/bin /bin
ENTRYPOINT ["/bin/simple-line-ai-bot"]
