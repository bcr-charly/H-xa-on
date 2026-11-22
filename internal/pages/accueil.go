package pages

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/bcr-charly/H-xa-on/internal/kernel"
)

// PageAccueil est la page d'accueil : photo d'équipe assombrie, logo au
// centre, slogan, bouton Connexion.
type PageAccueil struct {
	Noyau          *kernel.Noyau
	DossierContenu string
}

func (p *PageAccueil) Nom() string       { return "accueil" }
func (p *PageAccueil) Chemins() []string { return []string{"/accueil"} }

func (p *PageAccueil) Servir(w http.ResponseWriter, r *http.Request) {
	locale := kernel.LocaleDepuisRequete(r)
	connecte := p.Noyau.CompteCourant(r) != ""

	slogan := kernel.Traduit(locale, "accueil.slogan")
	libConnexion := kernel.Traduit(locale, "nav.connexion")
	libDeconnexion := kernel.Traduit(locale, "nav.deconnexion")
	var bouton string
	if connecte {
		bouton = fmt.Sprintf(`<a class="bouton" href="/deconnexion">%s</a>`,
			template.HTMLEscapeString(libDeconnexion))
	} else {
		bouton = fmt.Sprintf(`<a class="bouton" href="/connexion">%s</a>`,
			template.HTMLEscapeString(libConnexion))
	}

	contenu := template.HTML(fmt.Sprintf(`
<section class="accueil-hero">
	<div class="accueil-fond" style="background-image:url('/assets/equipe.svg')"></div>
	<div class="accueil-voile"></div>
	<div class="accueil-contenu">
		<img class="logo-grand" src="/assets/logo.svg" alt="%s">
		<p class="accueil-slogan">%s</p>
		%s
	</div>
</section>`,
		template.HTMLEscapeString(p.Noyau.Config.Nom),
		template.HTMLEscapeString(slogan),
		bouton,
	))
	p.Noyau.Rendu(w, http.StatusOK, page(p.Noyau, locale, p.Nom(), "", contenu, connecte), nil)
}
