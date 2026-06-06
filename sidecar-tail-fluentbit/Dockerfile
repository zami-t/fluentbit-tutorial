# app/Dockerfile
FROM golang:1.26.2 AS build
WORKDIR /src
COPY main.go .
COPY go.mod .
COPY go.sum .
RUN CGO_ENABLED=0 go build -o app main.go

FROM gcr.io/distroless/static-debian13
COPY --from=build /src/app /app
CMD ["/app"]

