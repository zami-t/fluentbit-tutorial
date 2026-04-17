# app/Dockerfile
FROM golang:1.26.2 AS build
WORKDIR /src
COPY main.go .
RUN go build -o /app main.go

FROM alpine:3.19
COPY --from=build /app /app
CMD ["/app"]

