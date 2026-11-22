package pages

import (
	"fmt"
	"html/template"
	"net/http"
	"sort"

	"github.com/bcr-charly/H-xa-on/internal/components"
	"github.com/bcr-charly/H-xa-on/internal/kernel"
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

	var colonne string
	if len(actus) == 0 {
		colonne = fmt.Sprintf(`<p class="actus-vide">%s</p>`,
			template.HTMLEscapeString(kernel.Traduit(locale, "actus.vide")))
	}
	for _, a := range actus {
		colonne += string(components.Defaut.Carte(a.Titre, a.Date, a.Texte, a.Image))
	}

	contenu := template.HTML(string(components.Defaut.BandeTitre(kernel.Traduit(locale, "actus.titre"))) +
		`<div class="conteneur"><div class="colonne-actus">` + colonne + `</div></div>`)
	p.Noyau.Rendu(w, http.StatusOK, page(p.Noyau, locale, p.Nom(), kernel.Traduit(locale, "actus.titre"), contenu, connecte), nil)
}
