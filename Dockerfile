FROM golang:1.21-alpine
WORKDIR /
COPY . .
RUN go build -o app .
CMD ["./app"]