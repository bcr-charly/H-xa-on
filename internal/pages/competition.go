package pages

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/bcr-charly/H-xa-on/internal/components"
	"github.com/bcr-charly/H-xa-on/internal/kernel"
)

// PageCompetition est la page dédiée d'une compétition, servie sous
// /competition/<slug>, remplie depuis contenu/competitions.txt.
type PageCompetition struct {
	Noyau          *kernel.Noyau
	DossierContenu string
}

func (p *PageCompetition) Nom() string       { return "competition" }
func (p *PageCompetition) Chemins() []string { return []string{"/competition/"} }

func (p *PageCompetition) Servir(w http.ResponseWriter, r *http.Request) {
	locale := kernel.LocaleDepuisRequete(r)
	connecte := p.Noyau.CompteCourant(r) != ""

	slug := strings.Trim(strings.TrimPrefix(r.URL.Path, "/competition/"), "/")
	blocs, err := kernel.ChargerContenus(p.DossierContenu + "/competitions.txt")
	if err != nil {
		http.Error(w, "contenu illisible", http.StatusInternalServerError)
		return
	}
	var competition *kernel.Contenu
	for _, b := range kernel.FiltreSections(blocs, "compet") {
		if slugDe(b.Texte("nom")) == slug || (b.Texte("slug") != "" && b.Texte("slug") == slug) {
			c := b
			competition = &c
			break
		}
	}
	if competition == nil {
		http.NotFound(w, r)
		return
	}

	image := competition.Texte("image")
	if image == "" {
		image = "/assets/trophee.svg"
	} else if !strings.HasPrefix(image, "/") {
		image = "/assets/" + image
	}
	description := competition.Texte("description")
	if description == "" {
		description = kernel.Traduit(locale, "compet.vide")
	}

	pageDetail := fmt.Sprintf(`
	<div class="conteneur">
		<div class="carte" style="max-width:40rem;margin:0 auto">
			<img class="carte-image" src="%s" alt="">
			<div class="carte-corps">
				<h3>%s</h3>
				<p>%s</p>
			</div>
		</div>
	</div>`,
		template.HTMLEscapeString(image),
		template.HTMLEscapeString(competition.Texte("nom")),
		template.HTMLEscapeString(description),
	)
	contenu := template.HTML(string(components.Defaut.BandeTitre(competition.Texte("nom"))) + pageDetail)
	p.Noyau.Rendu(w, http.StatusOK, page(p.Noyau, locale, "competition", kernel.Traduit(locale, "compet.titre"), contenu, connecte), nil)
}

// PageInscriptions est la page d'inscription à une compétition, cible de la
// bulle « s'inscrire » de la sous-page compétitions. Le parcours complet
// d'inscription (formulaire généré, validation, QR code) est prévu dans la
// feuille de route : le MVP affiche l'état des inscriptions et rattache
// chaque compétition à sa page dédiée.
type PageInscriptions struct {
	Noyau          *kernel.Noyau
	DossierContenu string
}

func (p *PageInscriptions) Nom() string       { return "inscriptions" }
func (p *PageInscriptions) Chemins() []string { return []string{"/inscriptions"} }

func (p *PageInscriptions) Servir(w http.ResponseWriter, r *http.Request) {
	locale := kernel.LocaleDepuisRequete(r)
	connecte := p.Noyau.CompteCourant(r) != ""

	var liste strings.Builder
	if blocs, err := kernel.ChargerContenus(p.DossierContenu + "/competitions.txt"); err == nil {
		for _, b := range kernel.FiltreSections(blocs, "compet") {
			fmt.Fprintf(&liste,
				`<li><a href="/competition/%s">%s</a></li>`,
				template.HTMLEscapeString(slugDe(b.Texte("nom"))),
				template.HTMLEscapeString(b.Texte("nom")))
		}
	}
	corps := fmt.Sprintf(`
	<div class="conteneur">
		<div class="carte" style="max-width:40rem;margin:0 auto">
			<div class="carte-corps">
				<p>%s</p>
				<ul>%s</ul>
			</div>
		</div>
		<p style="text-align:center;margin-top:1.5rem"><a class="bouton" href="/#competitions">%s</a></p>
	</div>`,
		template.HTMLEscapeString(kernel.Traduit(locale, "inscriptions.note")),
		liste.String(),
		template.HTMLEscapeString(kernel.Traduit(locale, "compet.retour")),
	)
	contenu := template.HTML(string(components.Defaut.BandeTitre(kernel.Traduit(locale, "inscriptions.titre"))) + corps)
	p.Noyau.Rendu(w, http.StatusOK, page(p.Noyau, locale, "inscriptions", kernel.Traduit(locale, "inscriptions.titre"), contenu, connecte), nil)
}

func slugDe(nom string) string { return slug(nom) }
