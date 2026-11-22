FROM golang:1.24-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY cmd ./cmd
COPY internal ./internal
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o /hexaion ./cmd/hexaion

FROM alpine:3.20
WORKDIR /app
COPY --from=build /hexaion /app/hexaion
COPY conf.txt /app/conf.txt
COPY contenu /app/contenu
COPY assets /app/assets
EXPOSE 8080
ENTRYPOINT ["/app/hexaion", "-dossier", "/app"]
