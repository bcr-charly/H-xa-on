# Mise en service — une page, orientée bureau des sports

## 1. Ce dont vous avez besoin

Un ordinateur (ou un Raspberry Pi, ou un serveur de l'école) et le dossier
`hexaion/` copié tel quel. Rien d'autre : pas de base de données, pas de
service externe, pas de compétence d'administration.

## 2. Premiers réglages (2 minutes)

Ouvrez `conf.txt` avec n'importe quel éditeur de texte :

- `port` : numéro de porte du site (8080 par défaut) ;
- `nom` : nom affiché (ex. « BDS ENSIIE ») ;
- `palette` : les cinq couleurs du site — optionnel, la palette dorée par
  défaut est livrée ;
- `pages` : `ON`/`OFF` pour chaque page.

Sauvegardez, relancez le binaire. C'est tout.

## 3. Changer les contenus (sans recompiler)

- Les actualités vivent dans `contenu/actus.txt` : un bloc `[actu]` par
  actualité, avec `titre`, `date`, `texte` (une ligne indentée = nouvelle
  ligne du texte) et `image` (fichier dans `assets/`). Les plus récentes en
  tête.
- L'équipe vit dans `contenu/equipe.txt` : un bloc `[membre]` par personne,
  avec `nom`, `role`, `photo` (fichier dans `assets/`), `lien` optionnel et
  `rang = bureau` pour les trois premiers postes (président, VP, trésorier),
  affichés dans la rangée d'honneur. Les autres membres se rangent en
  grille symétrique, centrée, sans recompilation du contenu.
- La sous-page présentation de l'accueil vit dans
  `contenu/accueil.txt` : `texte` (multi-lignes en indentant),
  `president` (légende sous la bulle), `photo` (fichier dans `assets/`).
- Les sports affichés dans la sous-page « Sports » vivent dans
  `contenu/sports.txt` : un bloc `[sport]` par sport avec `nom` et `image`.
  Cliquer sur une case ouvre la page dédiée du sport (`/sport/<nom>`),
  alimentée par les mêmes champs : `presentation` (multi-lignes),
  `competitions` (une par ligne indentée), `respo` et `contact`,
  `photo`/`photo2`, `horaires` (affichés à gauche) et `lieu` (à droite) ;
  le bouton « S'inscrire au ... » mène à l'onglet inscription.
- Les compétitions affichées sur l'accueil vivent dans
  `contenu/competitions.txt` : un bloc `[compet]` par compétition avec
  `nom`, `image`, `description` (page dédiée `/competition/<nom>`),
  `slug` optionnel. La bulle « s'inscrire » mène à la page
  `/inscriptions`.
- Les images vont dans `assets/`, le logo est `assets/logo.jpg`.

Après modification, relancez le binaire pour recharger les contenus.

## 4. Comptes et connexion

Une seule fois, en ligne de commande :

```sh
./hexaion -creer-admin login:motdepasse
```

Cela crée le compte administrateur dans `data/comptes.json`. Ensuite, la
création des comptes adhérents se fait par l'administrateur. Les mots de passe
ne sont jamais stockés en clair (hachage lent avec sel).

Un administrateur peut exporter toute collection en CSV :

```
http://votre-site/admin/export?collection=comptes
```

Fin d'année scolaire : la purge intégrale des données est une fonction du
noyau (`PurgeAnnuelle`), à déclencher à la prise en main de la nouvelle
équipe.

## 5. Lancer le site

```sh
./hexaion
```

Le site répond sur `http://localhost:8080` (ou le port configuré). Pour un
lancement au démarrage de la machine ou en arrière-plan, voir
[DEPLOIEMENT.md](DEPLOIEMENT.md) (service systemd, Docker, HTTPS).

## 6. En cas de souci

- Le terminal affiche les avertissements de configuration (clé inconnue, page
  activée mais non embarquée) : rien de cela ne fait planter le site.
- Un contenu mal formé arrête seulement la page concernée.
- Les données vivent dans `data/` : pour tout réinitialiser, videz `data/`.
