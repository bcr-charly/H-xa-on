package pages

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/bcr-charly/H-xa-on/internal/kernel"
)

// PageAccueil est la page d'accueil, ramenée à sa base : logo central,
// slogan, bouton Connexion. Le Header (barre de navigation, onglet de
// langue) et le logo permanent en haut à gauche sont rendus par la
// structure commune.
type PageAccueil struct {
	Noyau *kernel.Noyau
}

func (p *PageAccueil) Nom() string       { return "accueil" }
func (p *PageAccueil) Chemins() []string { return []string{"/accueil"} }

func (p *PageAccueil) Servir(w http.ResponseWriter, r *http.Request) {
	locale := kernel.LocaleDepuisRequete(r)
	connecte := p.Noyau.CompteCourant(r) != ""

	var bouton string
	if connecte {
		bouton = fmt.Sprintf(`<a class="bouton" href="/deconnexion">%s</a>`,
			template.HTMLEscapeString(kernel.Traduit(locale, "nav.deconnexion")))
	} else {
		bouton = fmt.Sprintf(`<a class="bouton" href="/connexion">%s</a>`,
			template.HTMLEscapeString(kernel.Traduit(locale, "nav.connexion")))
	}

	contenu := template.HTML(fmt.Sprintf(`
<section class="accueil-hero">
	<div class="accueil-contenu">
		<img class="logo-grand" src="/assets/logo.png" alt="%s">
		<p class="accueil-slogan">%s</p>
		<div>%s</div>
	</div>
</section>`,
		template.HTMLEscapeString(p.Noyau.Config.Nom),
		template.HTMLEscapeString(kernel.Traduit(locale, "accueil.slogan")),
		bouton,
	))
	p.Noyau.Rendu(w, http.StatusOK, page(p.Noyau, locale, p.Nom(), "", contenu, connecte), nil)
}
