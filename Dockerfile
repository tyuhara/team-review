FROM golang:1.17

WORKDIR /project

COPY ./go.* ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go install -v ./server

# Build Docker with Only Server Binary
FROM alpine:latest

RUN apk --no-cache add ca-certificates

COPY --from=0 /go/bin/server /bin/server

RUN addgroup -g 1001 http && adduser -D -G http -u 1001 http

USER 1001

CMD ["/bin/server"]
