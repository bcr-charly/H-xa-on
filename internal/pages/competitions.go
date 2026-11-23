package pages

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/bcr-charly/H-xa-on/internal/components"
	"github.com/bcr-charly/H-xa-on/internal/kernel"
)

// PageCompetitions est la page « Compétitions » : les grands rendez-vous
// de l'année sportive, tournois inter-écoles et championnats FFSU. Chaque
// bloc [compet] de contenu/competitions.txt donne une carte : nom, brève
// description, période de l'année si renseignée.
type PageCompetitions struct {
	Noyau          *kernel.Noyau
	DossierContenu string
}

func (p *PageCompetitions) Nom() string       { return "competitions" }
func (p *PageCompetitions) Chemins() []string { return []string{"/competitions"} }

func (p *PageCompetitions) Servir(w http.ResponseWriter, r *http.Request) {
	locale := kernel.LocaleDepuisRequete(r)
	connecte := p.Noyau.CompteCourant(r) != ""
	t := kernel.Traduit

	blocs, err := kernel.ChargerContenus(p.DossierContenu + "/competitions.txt")
	if err != nil {
		http.Error(w, "contenu illisible", http.StatusInternalServerError)
		return
	}
	compets := kernel.FiltreSections(blocs, "compet")

	var grille strings.Builder
	grille.WriteString(`<div class="cartes-colonne">`)
	for _, c := range compets {
		periode := ""
		if v := c.Texte("periode"); v != "" {
			periode = fmt.Sprintf(`<div class="carte-date">%s</div>`,
				template.HTMLEscapeString(v))
		}
		grille.WriteString(fmt.Sprintf(
			`<article class="carte"><img class="carte-image" src="%s" alt="%s" loading="lazy"><div class="carte-corps"><h3>%s</h3>%s<p>%s</p></div></article>`,
			template.HTMLEscapeString(cheminImage(c.Texte("image"))),
			template.HTMLEscapeString(c.Texte("nom")),
			template.HTMLEscapeString(c.Texte("nom")),
			periode,
			template.HTMLEscapeString(strings.TrimSpace(c.Texte("description"))),
		))
	}
	grille.WriteString(`</div>`)

	var vide string
	if len(compets) == 0 {
		vide = fmt.Sprintf(`<p class="actus-vide">%s</p>`,
			template.HTMLEscapeString(t(locale, "compet.vide")))
	}

	contenu := template.HTML(string(components.Defaut.BandeTitre(t(locale, "compet.titre"))) +
		`<div class="conteneur">` +
		`<p class="page-intro">` + template.HTMLEscapeString(t(locale, "compet.intro")) + `</p>` +
		grille.String() + vide + `</div>`)
	p.Noyau.Rendu(w, http.StatusOK, page(p.Noyau, locale, p.Nom(), t(locale, "compet.titre"), contenu, connecte), nil)
}
