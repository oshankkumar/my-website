FROM golang:1.7.3 AS build
WORKDIR /go/src/github.com/oshankkumar/my-website/
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=build /go/src/github.com/oshankkumar/my-website/app .
ENV PORT=8080
CMD ["./app"]