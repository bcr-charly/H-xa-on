package pages

import (
	"fmt"
	"html/template"
	"net/http"
	"sort"
	"strings"

	"github.com/bcr-charly/siteweb-bds-ensiie/internal/components"
	"github.com/bcr-charly/siteweb-bds-ensiie/internal/kernel"
)

// PageActus affiche les actualités lues dans contenu/actus.txt :
// une carte par entrée, les plus récentes en tête.
type PageActus struct {
	Noyau          *kernel.Noyau
	DossierContenu string
}

type actu struct {
	Titre, Date, Texte, Image string
}

func (p *PageActus) Nom() string       { return "actus" }
func (p *PageActus) Chemins() []string { return []string{"/actus"} }

func (p *PageActus) Servir(w http.ResponseWriter, r *http.Request) {
	locale := kernel.LocaleDepuisRequete(r)
	connecte := p.Noyau.CompteCourant(r) != ""

	blocs, err := kernel.ChargerContenus(p.DossierContenu + "/actus.txt")
	if err != nil {
		http.Error(w, "contenu illisible", http.StatusInternalServerError)
		return
	}
	var actus []actu
	for _, b := range kernel.FiltreSections(blocs, "actu") {
		actus = append(actus, actu{
			Titre: b.Texte("titre"),
			Date:  b.Texte("date"),
			Texte: b.Texte("texte"),
			Image: b.Texte("image"),
		})
	}
	sort.SliceStable(actus, func(i, j int) bool { return actus[i].Date > actus[j].Date })

	recherche := strings.TrimSpace(strings.ToLower(r.URL.Query().Get("q")))
	visibles := actus
	if recherche != "" {
		visibles = nil
		for _, a := range actus {
			if strings.Contains(strings.ToLower(a.Titre+" "+a.Texte), recherche) {
				visibles = append(visibles, a)
			}
		}
	}

	var colonne string
	if len(visibles) == 0 {
		colonne = fmt.Sprintf(`<p class="actus-vide">%s</p>`,
			template.HTMLEscapeString(kernel.Traduit(locale, "actus.vide")))
	}
	for _, a := range visibles {
		colonne += string(components.Defaut.Carte(a.Titre, a.Date, a.Texte, a.Image))
	}

	formRecherche := fmt.Sprintf(`<form class="actus-recherche" method="get" action="/actus">
<input type="search" name="q" value="%s" placeholder="%s" aria-label="%s">
</form>`,
		template.HTMLEscapeString(r.URL.Query().Get("q")),
		template.HTMLEscapeString(kernel.Traduit(locale, "actus.recherche")),
		template.HTMLEscapeString(kernel.Traduit(locale, "actus.recherche")),
	)

	contenu := template.HTML(string(components.Defaut.BandeTitre(kernel.Traduit(locale, "actus.titre"))) +
		`<div class="conteneur">` +
		`<p class="page-intro">` + template.HTMLEscapeString(kernel.Traduit(locale, "actus.intro")) + `</p>` +
		formRecherche + `<div class="colonne-actus">` + colonne + `</div></div>`)
	p.Noyau.Rendu(w, http.StatusOK, page(p.Noyau, locale, p.Nom(), kernel.Traduit(locale, "actus.titre"), contenu, connecte), nil)
}
