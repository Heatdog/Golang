FROM golang:latest

ENV GOPATH=/

COPY ./ ./

RUN go build -o hw6 ./cmd/redditclone/main.go
CMD ["./hw6"]