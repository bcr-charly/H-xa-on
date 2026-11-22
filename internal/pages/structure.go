package pages

import (
	"fmt"
	"html/template"
	"strings"

	"github.com/bcr-charly/H-xa-on/internal/kernel"
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
	--texte-clair: #101010;
	--texte-sombre: #F5F1E6;
}
* { box-sizing: border-box; margin: 0; padding: 0; }
html { scroll-behavior: smooth; }
body {
	background: var(--fond);
	color: var(--texte-sombre);
	font-family: system-ui, -apple-system, "Segoe UI", sans-serif;
	min-height: 100vh;
}
a { color: inherit; text-decoration: none; }

.entete {
	position: fixed;
	top: 0; left: 0; right: 0;
	z-index: 100;
	display: grid;
	grid-template-columns: 1fr auto 1fr;
	align-items: center;
	padding: 0.5rem 1rem;
	transition: transform 0.25s ease, background-color 0.25s ease;
}
.entete.cachee { transform: translateY(-100%); }
.entete-survol { background: #FFFFFF; color: #101010; }
.entete-survol .nav-onglet, .entete-survol .logo-bouton, .entete-survol .langue-select { color: #101010; }
.logo-petit {
	height: 2.6rem;
	width: auto;
	display: block;
}
.nav-barre {
	display: flex;
	gap: 0.25rem;
	list-style: none;
	justify-content: center;
	flex-wrap: wrap;
}
.nav-onglet {
	padding: 0.45rem 0.9rem;
	border-radius: 999px;
	font-size: 1rem;
	transition: transform 0.15s ease, background-color 0.15s ease;
}
.nav-onglet:hover, .nav-onglet:active {
	background: var(--bordure);
	transform: scale(1.08);
}
.nav-onglet[aria-current="page"] { outline: 1px solid var(--clair); }
.langue-zone { justify-self: end; }
.langue-select {
	background: transparent;
	border: 1px solid var(--bordure);
	color: inherit;
	border-radius: 8px;
	padding: 0.3rem 0.5rem;
	font-size: 0.95rem;
}

.accueil-hero {
	position: relative;
	min-height: 100vh;
	display: flex;
	flex-direction: column;
	align-items: center;
	justify-content: center;
	overflow: hidden;
	text-align: center;
}
.accueil-fond {
	position: absolute;
	inset: 0;
	background-size: cover;
	background-position: center;
	filter: brightness(0.55);
	z-index: 0;
}
.accueil-voile {
	position: absolute;
	inset: 0;
	background: var(--fond);
	opacity: 0.35;
	z-index: 1;
}
.accueil-contenu { position: relative; z-index: 2; display: flex; flex-direction: column; align-items: center; gap: 1.2rem; padding: 1rem; }
.logo-grand { height: 14rem; width: auto; }
.accueil-slogan { font-size: 1.35rem; max-width: 40rem; color: var(--texte-sombre); }

.bouton {
	display: inline-block;
	padding: 0.7rem 1.8rem;
	border-radius: 999px;
	background: var(--clair);
	color: var(--texte-clair);
	border: 1px solid var(--bordure);
	font-size: 1.05rem;
	font-weight: 600;
	transition: transform 0.15s ease, background-color 0.2s ease;
}
.bouton:hover { transform: scale(1.06); background: var(--accent); }
.bouton:active { transform: scale(1.03); }
.bouton-borde {
	background: transparent;
	color: var(--texte-sombre);
	border: 1px solid var(--bordure);
}

.bande-titre {
	padding: 7rem 1rem 2.5rem;
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

.grille-equipe {
	display: grid;
	grid-template-columns: repeat(auto-fill, minmax(11rem, 1fr));
	gap: 2rem 1rem;
	justify-items: center;
	padding-bottom: 3rem;
}
.icone-ronde { text-align: center; transition: transform 0.2s ease; }
.icone-ronde:hover, .icone-ronde:active { transform: scale(1.08); }
.icone-image {
	width: 9rem; height: 9rem;
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
.pied { padding: 2rem; text-align: center; color: var(--accent); }
@media (max-width: 48rem) {
	.logo-grand { height: 9rem; }
	.accueil-slogan { font-size: 1.05rem; }
}
`

// script porte les comportements de la barre de navigation :
//   - ordinateur : la barre disparaît en défilant vers le bas ; elle réapparaît
//     sur fond blanc opaque quand la souris approche du haut de l'écran ; elle
//     disparaît quand la souris s'éloigne ;
//   - tactile : glissement du doigt vers le haut en haut d'écran fait apparaître
//     la barre, glissement inverse la fait disparaître ; un appui en dehors
//     la referme ;
//   - sur la page d'accueil non défilée, la barre est affichée sans fond.
const script = `
(function () {
	"use strict";
	var entete = document.querySelector(".entete");
	if (!entete) return;
	var cachée = false;
	var dernierY = window.scrollY;

	function metre(v) { entete.classList.toggle("cachee", v); }
	function afficherSurvol() {
		metre(false);
		entete.classList.add("entete-survol");
	}
	function cacherSurvol() {
		entete.classList.remove("entete-survol");
		metre(true);
	}

	window.addEventListener("scroll", function () {
		var y = window.scrollY;
		if (y > dernierY && y > 60 && !entete.classList.contains("entete-survol")) {
			metre(true);
			cachée = true;
		} else if (y < dernierY && !cachée) {
			metre(y === 0 ? false : entete.classList.contains("cachee"));
		}
		dernierY = y;
		if (y === 0) { entete.classList.remove("entete-survol"); metre(false); }
	}, { passive: true });

	document.addEventListener("mousemove", function (e) {
		if (e.clientY <= 24) {
			if (cachée) afficherSurvol();
		} else if (entete.classList.contains("entete-survol") && e.clientY > 90) {
			cacherSurvol();
			cachée = true;
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
		if (fin.clientY < 120 && delta < -40) {
			metre(false);
			cachée = false;
		} else if (fin.clientY < 120 && delta > 40) {
			entete.classList.remove("entete-survol");
			metre(true);
			cachée = true;
		}
		debutY = null;
	}, { passive: true });

	document.addEventListener("click", function (e) {
		if (cachée) return;
		if (window.scrollY > 60 && entete.classList.contains("entete-survol") &&
			!entete.contains(e.target)) {
			cacherSurvol();
			cachée = true;
		}
	});
})();
`

// DonneesPage est le contexte passé à la structure commune de chaque page.
type DonneesPage struct {
	Palette kernel.Palette
	NomSite string
	Locale  string
	Entete  template.HTML
	Contenu template.HTML
	Survol  bool
}

// entete rend l'en-tête commun : logo-bouton permanent vers l'accueil,
// barre de navigation selon les pages actives, sélecteur de langue.
func entete(n *kernel.Noyau, locale, courant string, connecte bool) template.HTML {
	chemins := map[string]string{
		"actus":       "/actus",
		"planning":    "/planning",
		"equipe":      "/equipe",
		"sport":       "/sport",
		"competition": "/competition",
		"boutique":    "/boutique",
	}
	var onglets strings.Builder
	for _, nom := range []string{"actus", "planning", "equipe", "sport", "competition", "boutique"} {
		if active, ok := n.Config.Pages[nom]; !ok || !active {
			continue
		}
		courantAttr := ""
		if nom == courant {
			courantAttr = ` aria-current="page"`
		}
		fmt.Fprintf(&onglets,
			`<li><a class="nav-onglet" href="%s"%s>%s</a></li>`,
			chemins[nom], courantAttr,
			template.HTMLEscapeString(kernel.Traduit(locale, "nav."+nom)))
	}
	var connexion string
	if connecte {
		connexion = fmt.Sprintf(`<li><a class="nav-onglet" href="/deconnexion">%s</a></li>`,
			template.HTMLEscapeString(kernel.Traduit(locale, "nav.deconnexion")))
	} else {
		connexion = fmt.Sprintf(`<li><a class="nav-onglet" href="/connexion">%s</a></li>`,
			template.HTMLEscapeString(kernel.Traduit(locale, "nav.connexion")))
	}
	onglets.WriteString(connexion)

	locales := kernel.LocalesDisponibles()
	var options strings.Builder
	for _, l := range locales {
		sel := ""
		if l == locale {
			sel = ` selected`
		}
		fmt.Fprintf(&options, `<option value="%s"%s>%s</option>`, l, sel,
			template.HTMLEscapeString(kernel.Traduit(locale, "langue.nom."+l)))
	}
	sortie := strings.Join([]string{
		`<header class="entete" id="entete">`,
		`<a class="logo-bouton" href="/" aria-label="Accueil">` +
			`<img class="logo-petit" src="/assets/logo.svg" alt=""></a>`,
		`<nav aria-label="Navigation"><ul class="nav-barre">` + onglets.String() + `</ul></nav>`,
		`<div class="langue-zone"><form method="get" action="/langue" class="form-langue">` +
			`<input type="hidden" name="retour" value="` + cheminCourant(courant) + `">` +
			`<select class="langue-select" name="vers" aria-label="` +
			template.HTMLEscapeString(kernel.Traduit(locale, "nav.langue")) + `" onchange="this.form.submit()">` +
			options.String() + `</select></form></div>`,
		`</header>`,
	}, "")
	return template.HTML(sortie)
}

func cheminCourant(nom string) string {
	switch nom {
	case "actus", "equipe", "connexion":
		return "/" + nom
	default:
		return "/"
	}
}

// page rend la page complète : doctype, styles, entête, contenu, pied, script.
func page(n *kernel.Noyau, locale, nomPage, titre string, contenu template.HTML, connecte bool) string {
	dp := DonneesPage{
		Palette: n.Config.Palette,
		NomSite: n.Config.Nom,
		Locale:  locale,
		Entete:  entete(n, locale, nomPage, connecte),
		Contenu: contenu,
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
<body>
{{.Entete}}
<main>{{.Contenu}}</main>
<footer class="pied">` + "{{.NomSite}}" + `</footer>
<script>` + script + `</script>
</body>
</html>`))
	var sb strings.Builder
	if err := t.Execute(&sb, dp); err != nil {
		return "<!DOCTYPE html><html><body>erreur interne</body></html>"
	}
	return sb.String()
}
