FROM golang:1.22-alpine AS build
WORKDIR /src
COPY . .
RUN go build -o /out/server ./cmd/server

FROM alpine:3.19
WORKDIR /app
COPY --from=build /out/server /app/server
EXPOSE 8080
CMD ["/app/server"]


