package pages

import (
	"html/template"
	"net/http"
	"strings"

	"github.com/bcr-charly/H-xa-on/internal/components"
	"github.com/bcr-charly/H-xa-on/internal/kernel"
)

// PageEquipe affiche la grille des membres du bureau, remplie depuis
// contenu/equipe.txt (nom, rôle, photo, lien éventuel).
type PageEquipe struct {
	Noyau          *kernel.Noyau
	DossierContenu string
}

func (p *PageEquipe) Nom() string       { return "equipe" }
func (p *PageEquipe) Chemins() []string { return []string{"/equipe"} }

func (p *PageEquipe) Servir(w http.ResponseWriter, r *http.Request) {
	locale := kernel.LocaleDepuisRequete(r)
	connecte := p.Noyau.CompteCourant(r) != ""

	blocs, err := kernel.ChargerContenus(p.DossierContenu + "/equipe.txt")
	if err != nil {
		http.Error(w, "contenu illisible", http.StatusInternalServerError)
		return
	}

	var grille strings.Builder
	for _, b := range kernel.FiltreSections(blocs, "membre") {
		grille.WriteString(string(components.Defaut.IconeRonde(
			cheminImage(b.Texte("photo")),
			b.Texte("nom"),
			b.Texte("role"),
			b.Texte("lien"),
		)))
	}

	contenu := template.HTML(string(components.Defaut.BandeTitre(kernel.Traduit(locale, "equipe.titre"))) +
		`<div class="conteneur"><div class="grille-equipe">` + grille.String() + `</div></div>`)
	p.Noyau.Rendu(w, http.StatusOK, page(p.Noyau, locale, p.Nom(), kernel.Traduit(locale, "equipe.titre"), contenu, connecte), nil)
}

func cheminImage(photo string) string {
	if photo == "" {
		return "/assets/membre.svg"
	}
	if strings.HasPrefix(photo, "/") {
		return photo
	}
	return "/assets/" + photo
}
