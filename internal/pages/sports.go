package pages

import (
	"fmt"
	"html/template"
	"net/http"
	"sort"
	"strings"

	"github.com/bcr-charly/H-xa-on/internal/components"
	"github.com/bcr-charly/H-xa-on/internal/kernel"
)

// PageSports est la page « Nos sports » : barre de filtres par famille
// (champ categorie du contenu) au-dessus d'une grille de cartes, une
// carte par sport (photo, nom, bouton Découvrir), et compteur du nombre
// de sports affichés.
type PageSports struct {
	Noyau          *kernel.Noyau
	DossierContenu string
}

func (p *PageSports) Nom() string       { return "sports" }
func (p *PageSports) Chemins() []string { return []string{"/sports"} }

func (p *PageSports) Servir(w http.ResponseWriter, r *http.Request) {
	locale := kernel.LocaleDepuisRequete(r)
	connecte := p.Noyau.CompteCourant(r) != ""
	t := kernel.Traduit

	blocs, err := kernel.ChargerContenus(p.DossierContenu + "/sports.txt")
	if err != nil {
		http.Error(w, "contenu illisible", http.StatusInternalServerError)
		return
	}
	sports := kernel.FiltreSections(blocs, "sport")

	categories := map[string]bool{}
	for _, s := range sports {
		if c := s.Texte("categorie"); c != "" {
			categories[c] = true
		}
	}
	cles := make([]string, 0, len(categories))
	for c := range categories {
		cles = append(cles, c)
	}
	sort.Strings(cles)

	var compteurs strings.Builder
	for _, s := range sports {
		compteurs.WriteString(fmt.Sprintf("<span hidden data-categorie=\"%s\" data-total=\"%d\"></span>",
			template.HTMLEscapeString(s.Texte("categorie")), len(sports)))
	}

	var filtres strings.Builder
	filtres.WriteString(fmt.Sprintf(
		`<button class="filtre actif" data-filtre="">%s</button>`,
		template.HTMLEscapeString(t(locale, "sports.filtre.tous"))))
	for _, c := range cles {
		filtres.WriteString(fmt.Sprintf(
			`<button class="filtre" data-filtre="%s">%s</button>`,
			template.HTMLEscapeString(c), template.HTMLEscapeString(c)))
	}

	var grille strings.Builder
	grille.WriteString(`<div class="sports-grille">`)
	for _, s := range sports {
		nom := s.Texte("nom")
		grille.WriteString(fmt.Sprintf(
			`<a class="carte sport-carte" data-categorie="%s" href="/sport/%s">
<img class="carte-image" src="%s" alt="%s" loading="lazy">
<div class="carte-corps"><h3>%s</h3><span class="bouton bouton-borde">%s</span></div>
</a>`,
			template.HTMLEscapeString(s.Texte("categorie")),
			slugDe(nom),
			template.HTMLEscapeString(cheminImage(s.Texte("image"))),
			template.HTMLEscapeString(nom),
			template.HTMLEscapeString(nom),
			template.HTMLEscapeString(t(locale, "sports.decouvrir")),
		))
	}
	grille.WriteString(`</div>`)

	var vide string
	if len(sports) == 0 {
		vide = fmt.Sprintf(`<p class="actus-vide">%s</p>`,
			template.HTMLEscapeString(t(locale, "sports.vide")))
	}

	compte := fmt.Sprintf(t(locale, "sports.compte"), len(sports))
	contenu := template.HTML(string(components.Defaut.BandeTitre(t(locale, "sports.titre"))) +
		`<div class="conteneur">
<p class="page-intro">` + template.HTMLEscapeString(t(locale, "sports.intro")) + `</p>
<div class="sports-compte" id="sports-compte" data-unite="sports">` + compte + `</div>` + compteurs.String() +
		`<div class="sports-filtres">` + filtres.String() + `</div>` +
		grille.String() + vide +
		`<script>
(function () {
	var boutons = document.querySelectorAll(".filtre");
	var cartes = document.querySelectorAll(".sport-carte");
	var compte = document.getElementById("sports-compte");
	function appliquer(f) {
		var n = 0;
		cartes.forEach(function (c) {
			var ok = !f || c.getAttribute("data-categorie") === f;
			c.style.display = ok ? "" : "none";
			if (ok) n++;
		});
		if (compte) { compte.textContent = n + " " + compte.getAttribute("data-unite"); }
	}
	boutons.forEach(function (b) {
		b.addEventListener("click", function () {
			boutons.forEach(function (x) { x.classList.remove("actif"); });
			b.classList.add("actif");
			appliquer(b.getAttribute("data-filtre"));
		});
	});
})();
</script>
</div>`)
	p.Noyau.Rendu(w, http.StatusOK, page(p.Noyau, locale, p.Nom(), t(locale, "sports.titre"), contenu, connecte), nil)
}

// PageSport est la page dédiée d'un sport, servie sous /sport/<slug>,
// remplie depuis contenu/sports.txt : présentation, compétitions, respo,
// contact, horaires et lieu. Un bouton retour ramène à la liste.
type PageSport struct {
	Noyau          *kernel.Noyau
	DossierContenu string
}

func (p *PageSport) Nom() string       { return "sport" }
func (p *PageSport) Chemins() []string { return []string{"/sport/"} }

func (p *PageSport) Servir(w http.ResponseWriter, r *http.Request) {
	locale := kernel.LocaleDepuisRequete(r)
	connecte := p.Noyau.CompteCourant(r) != ""
	t := kernel.Traduit

	cle := strings.Trim(strings.TrimPrefix(r.URL.Path, "/sport/"), "/")
	blocs, err := kernel.ChargerContenus(p.DossierContenu + "/sports.txt")
	if err != nil {
		http.Error(w, "contenu illisible", http.StatusInternalServerError)
		return
	}
	var sport *kernel.Contenu
	for _, b := range kernel.FiltreSections(blocs, "sport") {
		if slugDe(b.Texte("nom")) == cle || (b.Texte("slug") != "" && b.Texte("slug") == cle) {
			c := b
			sport = &c
			break
		}
	}
	if sport == nil {
		http.NotFound(w, r)
		return
	}
	nom := sport.Texte("nom")

	var competitions strings.Builder
	if lignes := strings.TrimSpace(sport.Texte("competitions")); lignes != "" {
		for _, l := range strings.Split(lignes, "\n") {
			competitions.WriteString(fmt.Sprintf("<li>%s</li>", template.HTMLEscapeString(strings.TrimSpace(l))))
		}
	}
	var blocCompet string
	if competitions.Len() > 0 {
		blocCompet = fmt.Sprintf(`<div class="sport-detail-bloc"><h2>%s</h2><ul>%s</ul></div>`,
			template.HTMLEscapeString(t(locale, "sport.competitions")), competitions.String())
	}

	var blocRespo string
	if sport.Texte("respo") != "" || sport.Texte("contact") != "" {
		blocRespo = fmt.Sprintf(`<div class="sport-detail-bloc"><h2>%s</h2><p>%s</p><p>%s</p></div>`,
			template.HTMLEscapeString(t(locale, "sport.respo")),
			template.HTMLEscapeString(sport.Texte("respo")),
			template.HTMLEscapeString(strings.TrimSpace(sport.Texte("contact"))),
		)
	}

	horairesLieu := fmt.Sprintf(`<div class="sport-horaires-lieu">
<div><h2>%s</h2><p>%s</p></div>
<div><h2>%s</h2><p>%s</p></div>
</div>`,
		template.HTMLEscapeString(t(locale, "sport.horaires")),
		template.HTMLEscapeString(strings.TrimSpace(sport.Texte("horaires"))),
		template.HTMLEscapeString(t(locale, "sport.lieu")),
		template.HTMLEscapeString(strings.TrimSpace(sport.Texte("lieu"))),
	)

	detail := fmt.Sprintf(`<div class="sport-detail">
<h1>%s</h1>
<div class="sport-detail-intro">%s</div>
%s
%s
%s
</div>`,
		template.HTMLEscapeString(nom),
		template.HTMLEscapeString(strings.TrimSpace(sport.Texte("presentation"))),
		blocCompet,
		blocRespo,
		horairesLieu,
	)

	retour := fmt.Sprintf(`<div class="retour-haut"><a class="bouton bouton-borde" href="/sports">%s</a></div>`,
		template.HTMLEscapeString(t(locale, "sports.retour")))

	contenu := template.HTML(string(components.Defaut.BandeTitre(nom)) +
		retour + `<div class="conteneur">` + detail + `</div>`)
	p.Noyau.Rendu(w, http.StatusOK, page(p.Noyau, locale, p.Nom(), nom, contenu, connecte), nil)
}
