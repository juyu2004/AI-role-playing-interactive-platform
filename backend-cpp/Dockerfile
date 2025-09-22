FROM alpine:3.19 as build
RUN apk add --no-cache g++ cmake make asio-dev
WORKDIR /src
COPY . .
RUN g++ -std=c++17 -O2 main.cpp -lasio -o /out/server || g++ -std=c++17 -O2 main.cpp -o /out/server

FROM alpine:3.19
WORKDIR /app
COPY --from=build /out/server /app/server
EXPOSE 9090
CMD ["/app/server"]


