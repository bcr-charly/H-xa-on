# Déploiement : HTTPS, service systemd, Docker, arm64

Le binaire sert du **HTTP**. Le HTTPS est terminé par un reverse proxy en
amont (nginx, Caddy, ou le serveur de l'école) : c'est le montage standard,
documenté ici.

## Montage HTTPS standard avec nginx

```nginx
server {
    listen 443 ssl;
    server_name bds.mon-ecole.fr;

    ssl_certificate     /etc/ssl/certs/bds.pem;
    ssl_certificate_key /etc/ssl/private/bds.key;

    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Forwarded-Proto https;
    }
}
```

Équivalent Caddy, en une ligne (Caddy obtient et renouvelle le certificat) :

```
bds.mon-ecole.fr {
    reverse_proxy 127.0.0.1:8080
}
```

## Service systemd

`/etc/systemd/system/siteweb-bds-ensiie.service` :

```ini
[Unit]
Description=Site web du BDS ENSIIE
After=network.target

[Service]
WorkingDirectory=/opt/siteweb-bds-ensiie
ExecStart=/opt/siteweb-bds-ensiie/siteweb-bds-ensiie -dossier /opt/siteweb-bds-ensiie
Restart=on-failure

[Install]
WantedBy=multi-user.target
```

Puis : `systemctl enable --now siteweb-bds-ensiie`.

## Docker (optionnel)

Une image minimale contenant le dossier de déploiement ; le conteneur
n'expose que le port configuré :

```sh
docker build -t siteweb-bds-ensiie .
docker run -d -p 8080:8080 \
  -v $(pwd)/data:/app/data -v $(pwd)/contenu:/app/contenu siteweb-bds-ensiie
```

Le `Dockerfile` est fourni à la racine du projet.

## Cross-compilation arm64

Publier une variante Raspberry Pi ne demande qu'une cible de build
supplémentaire :

```sh
GOOS=linux GOARCH=arm64 go build -o siteweb-bds-ensiie-arm64 ./cmd/siteweb-bds-ensiie
```

## Build reproduit

```sh
CGO_ENABLED=0 go build -ldflags "-s -w" -o siteweb-bds-ensiie ./cmd/siteweb-bds-ensiie
```

`CGO_ENABLED=0` garantit un binaire statique, sans runtime ni dépendance à
installer.
