# Héxaïon

Site web en poche : un serveur web complet, autonome et personnalisable,
distribué sous licence libre pour le BDS de l'ENSIIE et pour toute association
étudiante.

- **Objectif local** : créer et exploiter le site web du BDS de l'ENSIIE.
- **Objectif général** : livrer un produit générique, déployable par
  n'importe quelle association étudiante, sans infrastructure ni compétences
  d'administration.

## Principe

Un binaire Go unique, un dossier copiable, aucune dépendance à installer :

```
hexaion/
├── hexaion            # binaire unique compilé
├── conf.txt           # configuration de l'instance
├── contenu/           # données éditables sans recompilation
│   ├── actus.txt
│   └── equipe.txt
├── assets/            # logo, photos, images des pages
└── data/              # données persistantes (comptes, sessions)
```

- Changer palette, pages actives, port : éditer `conf.txt`, relancer le binaire.
- Changer un contenu (actualité, membre) : éditer le fichier dans `contenu/`,
  relancer le binaire.
- Ajouter une page inédite : écrire un module Go dans `internal/pages/`
  (plus d'éventuels composants dans `internal/components/`), recompiler.

## Déploiement rapide

```sh
go build -o hexaion ./cmd/hexaion
./hexaion -creer-admin login:motdepasse   # une fois, crée le compte admin
./hexaion                                 # sert HTTP sur le port configuré
```

Cross-compilation vers une autre architecture, en une variable
d'environnement :

```sh
GOOS=linux GOARCH=arm64 go build -o hexaion-arm64 ./cmd/hexaion
```

La notice complète, orientée bureau des sports, est dans
[`docs/MISE-EN-SERVICE.md`](docs/MISE-EN-SERVICE.md). Le montage HTTPS par
reverse proxy (nginx/Caddy), le service systemd et Docker sont documentés dans
[`docs/DEPLOIEMENT.md`](docs/DEPLOIEMENT.md).

## Architecture

- `cmd/hexaion` : point d'entrée, lecture de `conf.txt`, création du compte
  administrateur initial.
- `internal/kernel` : noyau — serveur HTTP, routage, configuration, contenus,
  stockage (fichiers plats à réécriture atomique et export CSV), sessions
  (cookies signés HMAC, état dans `data/sessions`), i18n, hachage des mots de
  passe (PBKDF2, bibliothèque crypto standard).
- `internal/components` : bibliothèque de composants, un fichier par composant
  (bouton, carte, icône ronde, bande de titre). Toute page s'exprime avec ce
  vocabulaire.
- `internal/pages` : un module Go par page (Accueil, Actus, Équipe, Connexion),
  plus la structure commune (en-tête, barre de navigation, styles, script).
  Une page appelle des composants et lit ses contenus, jamais ceux d'une autre.

Le noyau ne connaît aucune page en particulier : chaque page s'enregistre dans
une table d'enregistrement (`internal/pages/pages.go`). Une page activée dans
`conf.txt` mais non embarquée produit un avertissement au démarrage et est
ignorée, jamais un plantage.

## Dépendances

**Aucune.** La bibliothèque standard Go est la base : `net/http`,
`html/template`, `crypto/pbkdf2`, `crypto/hmac`, `encoding/json`,
`encoding/csv`. Toute exception future devra être justifiée dans cette
section et son coût de réécriture évalué.

Aucun appel distant (polices, CDN, trackers) dans les pages livrées : chaque
ligne exécutée est maîtrisée. Unique exception prévue à la règle : le module
HelloAsso de la page boutique, non livré à ce jour (section 9 de la
spécification).

## Internationalisation

La mécanique fait partie du noyau dès le MVP (`internal/kernel/i18n.go`) :
toutes les chaînes affichées proviennent de tables de traduction par locale,
aucune n'est écrite en dur dans une page. Seule la traduction française est
livrée ; l'anglaise est fournie à titre de démonstration du câblage.

## Tests

```sh
go test ./...
```

## Licence

MIT — voir [LICENSE](LICENSE).
