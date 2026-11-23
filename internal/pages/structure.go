package pages

import (
	"fmt"
	"html/template"
	"sort"
	"strings"

	"github.com/bcr-charly/siteweb-bds-ensiie/internal/kernel"
)

// style contient tout l'habillage du site. Aucune ressource distante :
// ni police, ni CDN, ni tracker. Les couleurs de la palette sont injectées
// comme variables CSS depuis la configuration.
const style = `
:root {
	--fond: {{.Palette.Fond}};
	--elevation: {{.Palette.Elevation}};
	--clair: {{.Palette.Clair}};
	--accent: {{.Palette.Accent}};
	--bordure: {{.Palette.Bordure}};
	--texte-clair: #0B0A07;
	--texte-sombre: #FFFFFF;
	--texte-secondaire: #C0C0C0;
	--blanc: #FFFFFF;
}
* { box-sizing: border-box; margin: 0; padding: 0; }
html { scroll-behavior: smooth; }
body {
	background: var(--fond);
	color: var(--texte-sombre);
	font-family: system-ui, -apple-system, "Segoe UI", sans-serif;
	min-height: 100vh;
	display: flex;
	flex-direction: column;
}
a { color: inherit; text-decoration: none; }
main { flex: 1; }

/* ------------------------------------------------------------------ */
/* Logo permanent : bouton « menu home », en haut à gauche, visible en */
/* toute circonstance, indépendant du Header.                          */
/* ------------------------------------------------------------------ */
.logo-flottant {
	position: fixed;
	top: 0.5rem;
	left: 1rem;
	z-index: 200;
	display: inline-flex;
	align-items: center;
	padding: 0.3rem 0.6rem 0.3rem 0.3rem;
	border-radius: 999px;
	transition: transform 0.15s ease;
}
.logo-flottant:hover { transform: scale(1.08); }
.logo-flottant:active { transform: scale(1.03); }
.logo-petit {
	height: 2.6rem;
	width: auto;
	display: block;
}

/* ------------------------------------------------------------------ */
/* Header : nav, onglet de langue. Sur l'accueil il est transparent sur */
/* le hero ; sur les pages internes il devient une barre pleine, de la */
/* même taille dans les deux cas (plus de « header buggé »).           */
/* ------------------------------------------------------------------ */
.header {
	position: fixed;
	top: 0; left: 0; right: 0;
	z-index: 100;
	display: flex;
	justify-content: center;
	align-items: center;
	gap: 1rem;
	padding: 1rem 13rem 1rem 6rem;
	transition: transform 0.25s ease, background-color 0.25s ease, border-color 0.25s ease;
	color: var(--texte-sombre);
	background: transparent;
	border-bottom: 1px solid transparent;
}
.page-interne .header {
	background: var(--elevation);
	border-bottom-color: var(--bordure);
}
.header.cachee { transform: translateY(-100%); }
.nav-barre {
	display: flex;
	gap: 0.3rem;
	list-style: none;
	justify-content: center;
	align-items: center;
	flex-wrap: wrap;
}
.nav-onglet {
	padding: 0.6rem 1.15rem;
	border-radius: 999px;
	font-size: 1.15rem;
	transition: transform 0.15s ease, background-color 0.15s ease;
}
.nav-onglet:hover, .nav-onglet:active {
	background: var(--bordure);
	transform: scale(1.08);
}
.page-accueil .nav-onglet:hover, .page-accueil .nav-onglet:active {
	background: var(--elevation);
}
.nav-onglet[aria-current="page"] { outline: 1px solid var(--clair); }
.page-interne .nav-onglet[aria-current="page"] { outline-color: var(--accent); }

/* ------------------------------------------------------------------ */
/* Onglet de langue, tout à droite : pilule arrondie, menu au style    */
/* de la palette (élévation, bordure braise, coins arrondis).          */
/* ------------------------------------------------------------------ */
.langue-zone {
	position: absolute;
	top: 0.6rem;
	right: 1.2rem;
}
.langue-bouton {
	display: inline-flex;
	align-items: center;
	gap: 0.4rem;
	padding: 0.55rem 1.05rem;
	border-radius: 999px;
	border: 1px solid var(--bordure);
	background: transparent;
	color: inherit;
	font-size: 1rem;
	cursor: pointer;
	transition: transform 0.15s ease, background-color 0.15s ease;
}
.langue-bouton:hover, .langue-bouton:active { transform: scale(1.05); background: var(--elevation); }
.langue-menu {
	position: absolute;
	top: calc(100% + 0.5rem);
	right: 0;
	min-width: 10rem;
	background: var(--elevation);
	border: 1px solid var(--bordure);
	border-radius: 14px;
	padding: 0.4rem;
	list-style: none;
	box-shadow: 0 8px 24px rgba(0, 0, 0, 0.45);
	z-index: 300;
}
.langue-menu li a, .langue-menu li span {
	display: block;
	padding: 0.5rem 0.8rem;
	border-radius: 10px;
	font-size: 0.95rem;
	color: var(--texte-sombre);
}
.langue-menu li a:hover { background: var(--bordure); color: var(--texte-clair); }
.langue-menu li span { color: var(--accent); }
.langue-bloc[open] > summary { border-color: var(--clair); }
.langue-bloc summary::-webkit-details-marker { display: none; }
.langue-bloc summary { list-style: none; }

@media (max-width: 48rem) {
	.header { padding: 0.7rem 4.5rem 0.7rem 4.5rem; }
	.nav-onglet { padding: 0.5rem 0.8rem; font-size: 1rem; }
}

.accueil-hero {
	min-height: 100vh;
	display: flex;
	flex-direction: column;
	align-items: center;
	justify-content: center;
	text-align: center;
	padding: 5rem 1rem 2.5rem;
	background-image: linear-gradient(rgba(11, 7, 5, 0.55), rgba(11, 7, 5, 0.75)), url("/assets/fond-accueil.jpg");
	background-size: cover;
	background-position: center;
	background-repeat: no-repeat;
}
.accueil-contenu { position: relative; z-index: 2; display: flex; flex-direction: column; align-items: center; gap: 1.4rem; }
.logo-grand { height: 20rem; width: auto; }
.accueil-slogan { font-size: 1.35rem; max-width: 40rem; color: var(--texte-sombre); }
.accueil-boutons { display: flex; flex-wrap: wrap; gap: 0.9rem; justify-content: center; }

/* Sections de la page d'accueil, sur le modèle du BDS Télécom Paris. */
.accueil-section {
	max-width: 62rem;
	margin: 0 auto;
	padding: 4.5rem 1.25rem;
	text-align: center;
	display: flex;
	flex-direction: column;
	align-items: center;
	gap: 1.4rem;
}
.accueil-section h2 {
	font-size: clamp(1.6rem, 3.5vw, 2.3rem);
	color: var(--clair);
}
.accueil-section-alt { background: var(--elevation); }
.accueil-texte { max-width: 46rem; line-height: 1.65; color: var(--texte-sombre); }
.association-grille {
	display: grid;
	grid-template-columns: 1fr 1.2fr;
	gap: 2.5rem;
	align-items: center;
	width: 100%;
	max-width: 56rem;
	text-align: left;
}
.association-photo {
	aspect-ratio: 4 / 3;
	border-radius: 18px;
	border: 2px dashed var(--bordure);
	background:
		linear-gradient(135deg, rgba(255, 122, 0, 0.10), rgba(255, 165, 0, 0.05)),
		var(--elevation);
	display: flex;
	align-items: center;
	justify-content: center;
	overflow: hidden;
}
.association-photo span { color: var(--accent); font-size: 0.95rem; }
.association-photo img { width: 100%; height: 100%; object-fit: cover; display: block; }
.association-texte { display: flex; flex-direction: column; gap: 1.1rem; align-items: flex-start; }
.association-texte p { line-height: 1.65; color: var(--texte-sombre); }
.stats-grille {
	display: grid;
	grid-template-columns: repeat(4, 1fr);
	gap: 1.4rem;
	width: 100%;
	max-width: 46rem;
}
.stat {
	background: var(--elevation);
	border: 1px solid var(--bordure);
	border-radius: 14px;
	padding: 1.6rem 0.8rem;
	display: flex;
	flex-direction: column;
	gap: 0.4rem;
}
.accueil-section-alt .stat { background: var(--fond); }
.stat-chiffre { font-size: 2.2rem; font-weight: 800; color: var(--accent); line-height: 1; }
.stat-libelle { font-size: 0.95rem; color: var(--texte-sombre); }

.bouton {
	display: inline-block;
	padding: 0.7rem 1.8rem;
	border-radius: 999px;
	background: var(--accent);
	color: var(--fond);
	border: 1px solid var(--bordure);
	font-size: 1.05rem;
	font-weight: 600;
	transition: transform 0.15s ease, filter 0.2s ease;
}
.bouton:hover { transform: scale(1.06); filter: brightness(1.12); }
.bouton:active { transform: scale(1.03); }
.bouton-borde {
	background: transparent;
	color: var(--texte-sombre);
	border: 1px solid var(--bordure);
}

.bande-titre {
	padding: 5rem 1rem 1.6rem;
	text-align: center;
}
.bande-titre h1 { font-size: clamp(2rem, 5vw, 3.2rem); }

.conteneur { max-width: 62rem; margin: 0 auto; padding: 0 1rem 3rem; }
.colonne-actus { display: flex; flex-direction: column; gap: 1.4rem; }
.carte {
	background: var(--elevation);
	border: 1px solid var(--bordure);
	border-radius: 14px;
	overflow: hidden;
	transition: transform 0.15s ease;
}
.carte:hover { transform: scale(1.01); }
.carte-image { width: 100%; max-height: 22rem; object-fit: cover; display: block; }
.carte-corps { padding: 1.1rem 1.4rem; }
.carte-date { color: var(--accent); font-size: 0.9rem; margin-bottom: 0.3rem; }
.carte p { white-space: pre-line; margin-top: 0.4rem; }


/* Grille équipe : chaque rangée est centrée, la disposition reste
symétrique quelle que soit l'effectif (rangée bureau puis rangées
régulières). */
.grille-equipe {
	display: flex;
	flex-direction: column;
	gap: 2rem;
	width: 100%;
	align-items: stretch;
}
.rangee-equipe {
	display: grid;
	grid-template-columns: repeat(4, 1fr);
	gap: 1.6rem 1rem;
	width: 100%;
	justify-items: center;
}
.icone-ronde { text-align: center; transition: transform 0.2s ease; }
.icone-ronde:hover, .icone-ronde:active { transform: scale(1.08); }
.icone-image {
	width: 12rem; height: 12rem;
	border-radius: 50%;
	object-fit: cover;
	border: 2px solid var(--bordure);
}
.icone-nom { margin-top: 0.5rem; font-weight: 600; }
.icone-role { color: var(--accent); font-size: 0.9rem; }

.connexion-form {
	max-width: 24rem;
	margin: 0 auto;
	background: var(--elevation);
	border: 1px solid var(--bordure);
	border-radius: 14px;
	padding: 2rem;
	display: flex;
	flex-direction: column;
	gap: 1rem;
}
.connexion-form label { display: flex; flex-direction: column; gap: 0.3rem; }
.connexion-form input {
	background: var(--fond);
	border: 1px solid var(--bordure);
	color: var(--texte-sombre);
	padding: 0.6rem;
	border-radius: 8px;
	font-size: 1rem;
}
.erreur { color: var(--clair); }

/* ------------------------------------------------------------------ */
/* Pied de page : identité, réseaux, liens rapides, informations,      */
/* contact — sur chaque page.                                          */
/* ------------------------------------------------------------------ */
.pied {
	background: var(--elevation);
	border-top: 1px solid var(--bordure);
	margin-top: 0;
}
.pied-grille {
	max-width: 62rem;
	margin: 0 auto;
	padding: 2.5rem 1rem 1.5rem;
	display: grid;
	grid-template-columns: 2fr 1fr 1fr 1fr;
	gap: 2rem;
}
.pied-nom { font-size: 1.3rem; font-weight: 700; color: var(--clair); margin-bottom: 0.5rem; }
.pied-description { color: var(--texte-sombre); font-size: 0.95rem; }
.pied-sous-titre { color: var(--accent); font-size: 0.95rem; font-weight: 600; margin-bottom: 0.6rem; }
.pied-liste { list-style: none; display: flex; flex-direction: column; gap: 0.4rem; font-size: 0.95rem; }
.pied-liste a:hover { color: var(--clair); }
.pied-liste .desactive { color: var(--accent); }
.pied-reseaux { display: flex; gap: 0.8rem; margin-top: 1rem; align-items: center; }
.pied-reseaux a, .pied-reseaux span {
	display: inline-flex;
	padding: 0.45rem;
	border-radius: 50%;
	border: 1px solid var(--bordure);
	transition: transform 0.15s ease, background-color 0.15s ease;
}
.pied-reseaux a:hover { transform: scale(1.1); background: var(--bordure); }
.pied-reseaux .desactive { opacity: 0.45; }
.pied-reseaux svg { width: 1.4rem; height: 1.4rem; display: block; }
.pied-mentionlegale {
	border-top: 1px solid var(--bordure);
	color: var(--accent);
	font-size: 0.85rem;
	padding: 1rem;
	text-align: center;
}

/* ------------------------------------------------------------------ */
/* DA braise : sections d'accueil (déjà ci-dessus) et pages internes.  */
/* ------------------------------------------------------------------ */
.accueil-titre { font-size: clamp(1.6rem, 4vw, 2.6rem); max-width: 46rem; }
.page-intro { text-align: center; color: var(--texte-secondaire); max-width: 46rem; margin: 0 auto 1.8rem; }
.retour-haut { text-align: center; margin-top: 1.2rem; }
.accueil-section-stats { padding-top: 2rem; padding-bottom: 2rem; background: var(--elevation); }
.association-photo.grande { aspect-ratio: 16 / 9; }
.evenement-grille { display: grid; grid-template-columns: 1fr 1.2fr; gap: 2.5rem; align-items: center; width: 100%; max-width: 56rem; text-align: left; }

.sports-compte { text-align: center; color: var(--accent); font-weight: 700; font-size: 1.2rem; margin-bottom: 1rem; }
.sports-filtres { display: flex; flex-wrap: wrap; gap: 0.6rem; justify-content: center; margin-bottom: 2rem; }
.filtre {
	padding: 0.45rem 1.1rem;
	border-radius: 999px;
	border: 1px solid var(--bordure);
	background: transparent;
	color: var(--texte-sombre);
	font-size: 0.95rem;
	cursor: pointer;
	transition: transform 0.15s ease, background-color 0.15s ease, color 0.15s ease;
}
.filtre:hover { transform: scale(1.05); background: var(--elevation); }
.filtre.actif { background: var(--accent); color: var(--fond); border-color: var(--accent); }
.sports-grille { display: grid; grid-template-columns: repeat(3, 1fr); gap: 1.6rem; }
@media (max-width: 64rem) { .sports-grille { grid-template-columns: repeat(2, 1fr); } }
@media (max-width: 40rem) { .sports-grille { grid-template-columns: 1fr; } }
.sport-carte { display: flex; flex-direction: column; }
.sport-carte .carte-corps { display: flex; flex-direction: column; gap: 0.8rem; align-items: flex-start; }
.sport-carte:hover { transform: scale(1.02); }
.sport-carte .bouton { pointer-events: none; }
@media (max-width: 48rem) { .evenement-grille { grid-template-columns: 1fr; } }

.cartes-colonne { display: flex; flex-direction: column; gap: 1.6rem; }
.clubs-grille { display: grid; grid-template-columns: repeat(2, 1fr); gap: 1.6rem; max-width: 46rem; margin: 0 auto; }
@media (max-width: 48rem) { .clubs-grille { grid-template-columns: 1fr; } }
.club-carte .carte-corps { display: flex; flex-direction: column; gap: 0.5rem; align-items: flex-start; }
.partenaire-puces { list-style: none; display: flex; flex-direction: column; gap: 0.3rem; margin-top: 0.5rem; }
.partenaire-puces li { padding-left: 1.1rem; position: relative; }
.partenaire-puces li::before { content: ""; position: absolute; left: 0; top: 0.55em; width: 0.4em; height: 0.4em; border-radius: 50%; background: var(--accent); }
.partenaire-conditions { color: var(--texte-secondaire); font-size: 0.85rem; margin-top: 0.6rem; }

.actus-recherche { display: flex; justify-content: center; margin-bottom: 2rem; }
.actus-recherche input {
	width: min(30rem, 100%);
	padding: 0.6rem 1rem;
	border-radius: 999px;
	border: 1px solid var(--bordure);
	background: var(--elevation);
	color: var(--texte-sombre);
	font-size: 1rem;
}
.actus-recherche input:focus { outline: 1px solid var(--clair); }}

.adhesion { max-width: 46rem; }
.adhesion-section { margin-bottom: 2.5rem; }
.adhesion-section h2 { color: var(--clair); margin-bottom: 0.7rem; }
.adhesion-section p { line-height: 1.65; white-space: pre-line; }
.adhesion-puces { list-style: none; display: flex; flex-direction: column; gap: 0.4rem; }
.adhesion-puces li { padding-left: 1.2rem; position: relative; line-height: 1.5; }
.adhesion-puces li::before { content: ""; position: absolute; left: 0; top: 0.5em; width: 0.45em; height: 0.45em; border-radius: 50%; background: var(--accent); }
#adhesion-payer { margin-top: 1rem; }

@media (max-width: 48rem) {
	.association-grille { grid-template-columns: 1fr; text-align: center; }
	.association-texte { align-items: center; }
	.stats-grille { grid-template-columns: 1fr 1fr; }
	.logo-grand { height: 13rem; }
	.accueil-slogan { font-size: 1.05rem; }
	.header { padding: 0.5rem 3.5rem 0.5rem 4.5rem; }
	.icone-image { width: 8.5rem; height: 8.5rem; }
	.pied-grille { grid-template-columns: 1fr 1fr; }
}
`

// script porte les comportements du Header :
//   - ordinateur : le Header disparaît en défilant vers le bas ; il réapparaît
//     quand la souris approche du haut de l'écran ; il disparaît quand la
//     souris s'éloigne ;
//   - tactile : glissement du doigt vers le haut en haut d'écran fait
//     apparaître le Header, glissement inverse le fait disparaître ; un appui
//   - sur la page d'accueil non défilée, le Header est affiché sans fond.
//
// Le logo permanent, séparé du Header, reste visible en toute circonstance.
const script = `
(function () {
	"use strict";
	var header = document.querySelector(".header");
	if (!header) return;
	var caché = false;
	var dernierY = window.scrollY;

	function masquer(v) { header.classList.toggle("cachee", v); }

	window.addEventListener("scroll", function () {
		var y = window.scrollY;
		if (y > dernierY && y > 60) {
			masquer(true);
			caché = true;
		} else if (y < dernierY && y <= 60) {
			masquer(false);
			caché = false;
		}
		dernierY = y;
	}, { passive: true });

	document.addEventListener("mousemove", function (e) {
		if (e.clientY <= 24) {
			if (caché) { masquer(false); caché = false; }
		} else if (e.clientY > 120 && !header.contains(e.target)) {
			if (!caché && window.scrollY > 60) { masquer(true); caché = true; }
		}
	});

	var debutY = null;
	document.addEventListener("touchstart", function (e) {
		debutY = e.touches[0].clientY;
	}, { passive: true });
	document.addEventListener("touchend", function (e) {
		if (debutY === null) return;
		var fin = e.changedTouches[0];
		var delta = fin.clientY - debutY;
		if (fin.clientY < 140 && delta < -40) {
			masquer(false);
			caché = false;
		} else if (fin.clientY < 140 && delta > 40) {
			masquer(true);
			caché = true;
		}
		debutY = null;
	}, { passive: true });

	document.addEventListener("click", function (e) {
		if (caché) return;
		if (window.scrollY > 60 && !header.contains(e.target)) {
			masquer(true);
			caché = true;
		}
	});

})();
`

// DonneesPage est le contexte passé à la structure commune de chaque page.
type DonneesPage struct {
	Palette kernel.Palette
	NomSite string
	Locale  string
	Classe  string
	Logo    template.HTML
	Entete  template.HTML
	Contenu template.HTML
	Pied    template.HTML
}

// iconeSVG rend une icône de réseau social en SVG local (aucun appel distant).
func iconeSVG(nom string) template.HTML {
	switch nom {
	case "instagram":
		return template.HTML(`<svg viewBox="0 0 24 24" fill="currentColor" aria-hidden="true"><rect x="3" y="3" width="18" height="18" rx="5" ry="5" fill="none" stroke="currentColor" stroke-width="2"/><circle cx="12" cy="12" r="4" fill="none" stroke="currentColor" stroke-width="2"/><circle cx="17.2" cy="6.8" r="1.3"/></svg>`)
	case "linkedin":
		return template.HTML(`<svg viewBox="0 0 24 24" fill="currentColor" aria-hidden="true"><rect x="3" y="3" width="18" height="18" rx="3" ry="3" fill="none" stroke="currentColor" stroke-width="2"/><rect x="7" y="10" width="2.6" height="7" rx="1"/><circle cx="8.3" cy="7.4" r="1.5"/><path d="M12.5 17v-4.2c0-1.3 1-2.3 2.3-2.3s2.2 1 2.2 2.3V17" fill="none" stroke="currentColor" stroke-width="2.6" stroke-linecap="round"/></svg>`)
	default: // ensiie : blason simplifié de l'école
		return template.HTML(`<svg viewBox="0 0 24 24" fill="currentColor" aria-hidden="true"><circle cx="12" cy="12" r="10" fill="none" stroke="currentColor" stroke-width="2"/><path d="M8 15V9l4 4 4-4v6" fill="none" stroke="currentColor" stroke-width="2" stroke-linejoin="round" stroke-linecap="round"/></svg>`)
	}
}

// entete rend le Header : barre de navigation selon les pages actives,
// dans l'ordre du site de référence : L'association (accueil), Nos sports,
// Compétitions, Clubs, Actualités, Adhésion, Équipe. Connexion et
// déclose selon l'état de session. Le site est monolingue : pas d'onglet

func entete(n *kernel.Noyau, locale, courant string, connecte bool) template.HTML {
	ordrenav := []string{"sports", "competitions", "clubs", "equipe"}
	cles := map[string]string{
		"sports":       "nav.sports",
		"competitions": "nav.competitions",
		"clubs":        "nav.clubs",
		"equipe":       "nav.equipe",
	}
	var onglets strings.Builder
	for _, nom := range ordrenav {
		if active, ok := n.Config.Pages[nom]; !ok || !active {
			continue
		}
		courantAttr := ""
		if nom == courant {
			courantAttr = ` aria-current="page"`
		}
		fmt.Fprintf(&onglets,
			`<li><a class="nav-onglet" href="/%s"%s>%s</a></li>`,
			nom, courantAttr,
			template.HTMLEscapeString(kernel.Traduit(locale, cles[nom])))
	}
	fmt.Fprintf(&onglets, `<li><a class="nav-onglet" href="/connexion">%s</a></li>`,
		template.HTMLEscapeString(kernel.Traduit(locale, "nav.connexion")))
	if connecte {
		fmt.Fprintf(&onglets, `<li><a class="nav-onglet" href="/deconnexion">%s</a></li>`,
			template.HTMLEscapeString(kernel.Traduit(locale, "nav.deconnexion")))
	}

	locales := kernel.LocalesDisponibles()
	sort.Strings(locales)
	var options strings.Builder
	for _, l := range locales {
		if l == locale {
			fmt.Fprintf(&options, `<li><span aria-current="true">%s</span></li>`,
				template.HTMLEscapeString(kernel.Traduit(locale, "langue.nom."+l)))
			continue
		}
		fmt.Fprintf(&options, `<li><a href="/langue?vers=%s&amp;retour=%s">%s</a></li>`,
			template.HTMLEscapeString(l),
			template.HTMLEscapeString(cheminCourant(courant)),
			template.HTMLEscapeString(kernel.Traduit(locale, "langue.nom."+l)))
	}

	sortie := strings.Join([]string{
		`<header class="header" id="header">`,
		`<nav aria-label="Navigation"><ul class="nav-barre">` + onglets.String() + `</ul></nav>`,
		`<div class="langue-zone"><details class="langue-bloc">` +
			`<summary class="langue-bouton" aria-label="` +
			template.HTMLEscapeString(kernel.Traduit(locale, "nav.langue")) + `">🌐 ` +
			template.HTMLEscapeString(locale) + `</summary>` +
			`<ul class="langue-menu">` + options.String() + `</ul></details></div>`,
		`</header>`,
	}, "")
	return template.HTML(sortie)
}

// slugDe transforme un nom en identifiant d'URL : minuscules, accents
// retirés, séparateurs normalisés.
func slugDe(nom string) string {
	majuscules := map[rune]rune{
		'À': 'a', 'Á': 'a', 'Â': 'a', 'Ã': 'a', 'Ä': 'a', 'Å': 'a',
		'È': 'e', 'É': 'e', 'Ê': 'e', 'Ë': 'e',
		'Ì': 'i', 'Í': 'i', 'Î': 'i', 'Ï': 'i',
		'Ò': 'o', 'Ó': 'o', 'Ô': 'o', 'Õ': 'o', 'Ö': 'o',
		'Ù': 'u', 'Ú': 'u', 'Û': 'u', 'Ü': 'u',
		'Ç': 'c', 'Ñ': 'n',
	}
	var sb strings.Builder
	for _, c := range strings.ToLower(nom) {
		if r, ok := majuscules[c]; ok {
			sb.WriteRune(r)
			continue
		}
		if (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') {
			sb.WriteRune(c)
			continue
		}
		if c == ' ' || c == '-' || c == '_' || c == '/' || c == '\'' {
			sb.WriteRune('-')
		}
	}
	res := strings.Trim(sb.String(), "-")
	for strings.Contains(res, "--") {
		res = strings.ReplaceAll(res, "--", "-")
	}
	if res == "" {
		res = "x"
	}
	return res
}

// logoPermanent rend le bouton « menu home », toujours visible, qui ramène
// à l'accueil.
func logoPermanent(nomSite string) template.HTML {
	return template.HTML(`<a class="logo-flottant" href="/" aria-label="Accueil">` +
		`<img class="logo-petit" src="/assets/logo.png?v=2026-11-b" alt="` + template.HTMLEscapeString(nomSite) + `"></a>`)
}

// pied rend le pied de page : identité du BDS, réseaux sociaux, liens
// rapides, informations, contact. Les sections « liens rapides » et FAQ
// sont prévues dès maintenant, leurs destinations arriveront avec les pages.
func pied(n *kernel.Noyau, locale string) template.HTML {
	t := kernel.Traduit
	var liensRapides strings.Builder
	type lienRapide struct {
		cle   string
		url   string
		actif bool
	}
	liens := []lienRapide{
		{"pied.sports", "/sports", n.Config.Pages["sports"]},
		{"pied.competitions", "/competitions", n.Config.Pages["competitions"]},
		{"pied.clubs", "/clubs", n.Config.Pages["clubs"]},
		{"pied.equipe", "/equipe", n.Config.Pages["equipe"]},
	}
	for _, l := range liens {
		if l.actif {
			fmt.Fprintf(&liensRapides, `<li><a href="%s">%s</a></li>`, l.url,
				template.HTMLEscapeString(t(locale, l.cle)))
		} else {
			fmt.Fprintf(&liensRapides, `<li class="desactive">%s (%s)</li>`,
				template.HTMLEscapeString(t(locale, l.cle)),
				template.HTMLEscapeString(t(locale, "pied.avenir")))
		}
	}

	instagram := fmt.Sprintf(
		`<a href="https://www.instagram.com/bds_ensiie/" aria-label="Instagram" rel="noopener">%s</a>`,
		iconeSVG("instagram"))
	linkedin := fmt.Sprintf(
		`<span class="desactive" aria-label="LinkedIn — %s" title="LinkedIn — %s">%s</span>`,
		template.HTMLEscapeString(t(locale, "pied.avenir")),
		template.HTMLEscapeString(t(locale, "pied.avenir")),
		iconeSVG("linkedin"))
	ensiie := fmt.Sprintf(
		`<a href="https://www.ensiie.fr/" aria-label="%s" rel="noopener">%s</a>`,
		template.HTMLEscapeString(t(locale, "pied.siteecole")),
		iconeSVG("ensiie"))

	sortie := strings.Join([]string{
		`<div class="pied-grille">`,
		`<div>` +
			`<div class="pied-nom">` + template.HTMLEscapeString(n.Config.Nom) + `</div>` +
			`<p class="pied-description">` + template.HTMLEscapeString(t(locale, "pied.description")) + `</p>` +
			`<div class="pied-reseaux">` + instagram + linkedin + ensiie + `</div>` +
			`</div>`,
		`<div><div class="pied-sous-titre">` + template.HTMLEscapeString(t(locale, "pied.liensrapides")) + `</div>` +
			`<ul class="pied-liste">` + liensRapides.String() + `</ul></div>`,
		`<div><div class="pied-sous-titre">` + template.HTMLEscapeString(t(locale, "pied.informations")) + `</div>` +
			`<ul class="pied-liste">` +
			`<li class="desactive">FAQ (` + template.HTMLEscapeString(t(locale, "pied.avenir")) + `)</li>` +
			`<li><a href="#apropos">` + template.HTMLEscapeString(t(locale, "pied.apropos")) + `</a></li>` +
			`</ul></div>`,
		`<div><div class="pied-sous-titre">` + template.HTMLEscapeString(t(locale, "pied.contact")) + `</div>` +
			`<address class="pied-liste" style="font-style:normal">` +
			`<li>1 square de la Résistance<br>91025 Évry-Courcouronnes</li>` +
			`<li><a href="mailto:info@ensiie.fr">info@ensiie.fr</a></li>` +
			`<li><a href="tel:+33169367350">+33 (0)1 69 36 73 50</a></li>` +
			`</address></div>`,
		`</div>`,
		`<div class="pied-mentionlegale">© ` + fmt.Sprint(anneeCourant()) + ` ` +
			template.HTMLEscapeString(n.Config.Nom) + ` — ` +
			template.HTMLEscapeString(t(locale, "pied.description")) + `</div>`,
	}, "")
	return template.HTML(sortie)
}

func cheminCourant(nom string) string {
	switch nom {
	case "actus", "equipe", "connexion", "inscription", "sports", "competitions", "clubs":
		return "/" + nom
	default:
		return "/"
	}
}

var anneeCourant = func() int { return kernel.DateNow().Year() }

// page rend la page complète : doctype, styles, logo permanent, Header,
// contenu, pied, script.
func page(n *kernel.Noyau, locale, nomPage, titre string, contenu template.HTML, connecte bool) string {
	classe := "page-interne"
	if nomPage == "accueil" {
		classe = "page-accueil"
	}
	dp := DonneesPage{
		Palette: n.Config.Palette,
		NomSite: n.Config.Nom,
		Locale:  locale,
		Classe:  classe,
		Logo:    logoPermanent(n.Config.Nom),
		Entete:  entete(n, locale, nomPage, connecte),
		Contenu: contenu,
		Pied:    pied(n, locale),
	}
	t := template.Must(template.New("structure").Parse(
		`<!DOCTYPE html>
<html lang="{{.Locale}}">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>{{.NomSite}}</title>
<style>` + style + `</style>
</head>
<body class="{{.Classe}}">
{{.Logo}}
{{.Entete}}
<main>{{.Contenu}}</main>
{{.Pied}}
<script>` + script + `</script>
</body>
</html>`))
	var sb strings.Builder
	if err := t.Execute(&sb, dp); err != nil {
		return "<!DOCTYPE html><html><body>erreur interne</body></html>"
	}
	return sb.String()
}
