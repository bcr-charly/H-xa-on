package pages

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/bcr-charly/H-xa-on/internal/components"
	"github.com/bcr-charly/H-xa-on/internal/kernel"
)

// PageClubs est la page « Les clubs et leur financement » : le BDS
// soutient des clubs et associations étudiantes dans leurs projets.
// Un bloc [club] de contenu/clubs.txt donne une carte : titre et
// description.
type PageClubs struct {
	Noyau          *kernel.Noyau
	DossierContenu string
}

func (p *PageClubs) Nom() string       { return "clubs" }
func (p *PageClubs) Chemins() []string { return []string{"/clubs"} }

func (p *PageClubs) Servir(w http.ResponseWriter, r *http.Request) {
	locale := kernel.LocaleDepuisRequete(r)
	connecte := p.Noyau.CompteCourant(r) != ""
	t := kernel.Traduit

	blocs, err := kernel.ChargerContenus(p.DossierContenu + "/clubs.txt")
	if err != nil {
		http.Error(w, "contenu illisible", http.StatusInternalServerError)
		return
	}
	clubs := kernel.FiltreSections(blocs, "club")

	var grille strings.Builder
	grille.WriteString(`<div class="clubs-grille">`)
	for _, c := range clubs {
		grille.WriteString(fmt.Sprintf(
			`<article class="carte club-carte"><div class="carte-corps"><h3>%s</h3><p>%s</p></div></article>`,
			template.HTMLEscapeString(c.Texte("nom")),
			template.HTMLEscapeString(strings.TrimSpace(c.Texte("description"))),
		))
	}
	grille.WriteString(`</div>`)

	var vide string
	if len(clubs) == 0 {
		vide = fmt.Sprintf(`<p class="actus-vide">%s</p>`,
			template.HTMLEscapeString(t(locale, "clubs.vide")))
	}

	contenu := template.HTML(string(components.Defaut.BandeTitre(t(locale, "clubs.titre"))) +
		`<div class="conteneur">` +
		`<p class="page-intro">` + template.HTMLEscapeString(t(locale, "clubs.intro")) + `</p>` +
		grille.String() + vide + `</div>`)
	p.Noyau.Rendu(w, http.StatusOK, page(p.Noyau, locale, p.Nom(), t(locale, "clubs.titre"), contenu, connecte), nil)
}
