package pages

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/bcr-charly/H-xa-on/internal/kernel"
)

// PageAccueil reproduit l'agencement du site du BDS Télécom Paris,
// avec l'identité visuelle de l'ENSIIE : hero plein écran sur le fond
// photo, section « L'association » et ses compteurs, « Les compétitions »,
// « Nos partenaires ». Le Header (barre de navigation, onglet de langue)
// et le logo permanent en haut à gauche sont rendus par la structure
// commune.
type PageAccueil struct {
	Noyau          *kernel.Noyau
	DossierContenu string
}

func (p *PageAccueil) Nom() string       { return "accueil" }
func (p *PageAccueil) Chemins() []string { return []string{"/accueil"} }

// statSection compte les blocs d'une section d'un fichier de contenu.
func (p *PageAccueil) statSection(fichier, section string) int {
	blocs, err := kernel.ChargerContenus(p.DossierContenu + "/" + fichier)
	if err != nil {
		return 0
	}
	return len(kernel.FiltreSections(blocs, section))
}

func (p *PageAccueil) Servir(w http.ResponseWriter, r *http.Request) {
	locale := kernel.LocaleDepuisRequete(r)
	connecte := p.Noyau.CompteCourant(r) != ""
	t := kernel.Traduit

	var boutonCompte string
	if connecte {
		boutonCompte = fmt.Sprintf(`<a class="bouton bouton-borde" href="/deconnexion">%s</a>`,
			template.HTMLEscapeString(t(locale, "nav.deconnexion")))
	} else {
		boutonCompte = fmt.Sprintf(`<a class="bouton bouton-borde" href="/connexion">%s</a>`,
			template.HTMLEscapeString(t(locale, "nav.connexion")))
	}

	membres := p.statSection("equipe.txt", "membre")
	actus := p.statSection("actus.txt", "actu")

	contenu := template.HTML(fmt.Sprintf(`
<section class="accueil-hero">
	<div class="accueil-contenu">
		<img class="logo-grand" src="/assets/logo.png?v=2026-11-b" alt="%s">
		<p class="accueil-slogan">%s</p>
		<div class="accueil-boutons">
			<a class="bouton" href="/actus">%s</a>
			%s
		</div>
	</div>
</section>

<section class="accueil-section" id="association">
	<h2>%s</h2>
	<div class="association-grille">
		<div class="association-photo"><span>%s</span></div>
		<div class="association-texte">
			<p>%s</p>
			<p>%s</p>
			<a class="bouton" href="/equipe">%s</a>
		</div>
	</div>
	<div class="stats-grille">
		<div class="stat"><div class="stat-chiffre">%d</div><div class="stat-libelle">%s</div></div>
		<div class="stat"><div class="stat-chiffre">%d</div><div class="stat-libelle">%s</div></div>
		<div class="stat"><div class="stat-chiffre">100%%</div><div class="stat-libelle">%s</div></div>
		<div class="stat"><div class="stat-chiffre">1</div><div class="stat-libelle">%s</div></div>
	</div>
</section>

<section class="accueil-section accueil-section-alt" id="competitions">
	<h2>%s</h2>
	<p class="accueil-texte">%s</p>
	<a class="bouton" href="/actus">%s</a>
</section>

<section class="accueil-section" id="partenaires">
	<h2>%s</h2>
	<p class="accueil-texte">%s</p>
</section>`,
		template.HTMLEscapeString(p.Noyau.Config.Nom),
		template.HTMLEscapeString(t(locale, "accueil.slogan")),
		template.HTMLEscapeString(t(locale, "accueil.bouton.actus")),
		boutonCompte,
		template.HTMLEscapeString(t(locale, "accueil.association.titre")),
		template.HTMLEscapeString(t(locale, "accueil.association.photo")),
		fmt.Sprintf(template.HTMLEscapeString(t(locale, "accueil.association.p1")), membres),
		template.HTMLEscapeString(t(locale, "accueil.association.p2")),
		template.HTMLEscapeString(t(locale, "accueil.presentation.bouton")),
		membres, template.HTMLEscapeString(t(locale, "accueil.stats.membres")),
		actus, template.HTMLEscapeString(t(locale, "accueil.stats.actus")),
		template.HTMLEscapeString(t(locale, "accueil.stats.benevoles")),
		template.HTMLEscapeString(t(locale, "accueil.stats.objectif")),
		template.HTMLEscapeString(t(locale, "accueil.competitions.titre")),
		template.HTMLEscapeString(t(locale, "accueil.competitions.texte")),
		template.HTMLEscapeString(t(locale, "accueil.competitions.bouton")),
		template.HTMLEscapeString(t(locale, "accueil.sponsors.titre")),
		template.HTMLEscapeString(t(locale, "accueil.sponsors.vide")),
	))
	p.Noyau.Rendu(w, http.StatusOK, page(p.Noyau, locale, p.Nom(), "", contenu, connecte), nil)
}
