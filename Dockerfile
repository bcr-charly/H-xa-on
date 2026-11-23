FROM golang:1.24-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY cmd ./cmd
COPY internal ./internal
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o /siteweb-bds-ensiie ./cmd/siteweb-bds-ensiie

FROM alpine:3.20
WORKDIR /app
COPY --from=build /siteweb-bds-ensiie /app/siteweb-bds-ensiie
COPY conf.txt /app/conf.txt
COPY contenu /app/contenu
COPY assets /app/assets
EXPOSE 8080
ENTRYPOINT ["/app/siteweb-bds-ensiie", "-dossier", "/app"]
