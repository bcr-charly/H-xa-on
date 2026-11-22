package pages

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/bcr-charly/H-xa-on/internal/kernel"
)

// PageAccueil est la page d'accueil : une grande page qui se déroule vers le
// bas comme un assemblage de sous-pages.
//
//	Sous-page 1 — hero : logo, slogan, bouton Connexion.
//	Sous-page 2 — présentation rapide de l'association.
//	Sous-page 3 — compétitions auxquelles le BDS a participé.
//
// Le fond photo reste fixe pendant le défilement ; seul le contenu monte et
// disparaît. Le contenu éditorial provient de contenu/accueil.txt et
// contenu/competitions.txt.
type PageAccueil struct {
	Noyau          *kernel.Noyau
	DossierContenu string
}

func (p *PageAccueil) Nom() string       { return "accueil" }
func (p *PageAccueil) Chemins() []string { return []string{"/accueil"} }

func (p *PageAccueil) Servir(w http.ResponseWriter, r *http.Request) {
	locale := kernel.LocaleDepuisRequete(r)
	connecte := p.Noyau.CompteCourant(r) != ""

	contenu := template.HTML(
		fondFixe() +
			p.sectionHero(locale, connecte) +
			p.sectionPresentation(locale) +
			p.sectionCompetitions(locale))
	p.Noyau.Rendu(w, http.StatusOK, page(p.Noyau, locale, p.Nom(), "", contenu, connecte), nil)
}

// fondFixe rend la couche de fond : photo assombrie par un voile de la
// couleur de fond, fixe pendant tout le défilement de la page.
func fondFixe() string {
	return `<div class="fond-fixe" aria-hidden="true">` +
		`<div class="accueil-fond" style="background-image:url('/assets/equipe.svg')"></div>` +
		`<div class="accueil-voile"></div></div>`
}

// sectionHero est la première sous-page : logo, slogan, bouton Connexion.
func (p *PageAccueil) sectionHero(locale string, connecte bool) string {
	slogan := kernel.Traduit(locale, "accueil.slogan")
	var bouton string
	if connecte {
		bouton = fmt.Sprintf(`<a class="bouton bouton-borde" href="/deconnexion">%s</a>`,
			template.HTMLEscapeString(kernel.Traduit(locale, "nav.deconnexion")))
	} else {
		bouton = fmt.Sprintf(`<a class="bouton" href="/connexion">%s</a>`,
			template.HTMLEscapeString(kernel.Traduit(locale, "nav.connexion")))
	}
	return fmt.Sprintf(`
<section class="accueil-section accueil-hero">
	<div class="accueil-contenu">
		<img class="logo-grand" src="/assets/logo.svg" alt="%s">
		<p class="accueil-slogan">%s</p>
		%s
	</div>
</section>`,
		template.HTMLEscapeString(p.Noyau.Config.Nom),
		template.HTMLEscapeString(slogan),
		bouton,
	)
}

// sectionPresentation est la deuxième sous-page : description rapide du
// BDS dans une zone à coins arrondis, espace photo à droite (logo et
// président en bulles rondes), bouton de redirection vers la page BDS.
func (p *PageAccueil) sectionPresentation(locale string) string {
	texte := kernel.Traduit(locale, "pied.description")
	president := ""
	photoPresident := "/assets/membre.svg"
	if blocs, err := kernel.ChargerContenus(p.DossierContenu + "/accueil.txt"); err == nil {
		if presents := kernel.FiltreSections(blocs, "presentation"); len(presents) > 0 {
			b := presents[0]
			if v := b.Texte("texte"); v != "" {
				texte = v
			}
			president = b.Texte("president")
			if v := b.Texte("photo"); v != "" {
				photoPresident = "/assets/" + v
			}
		}
	}
	legende := template.HTMLEscapeString(president)
	if president == "" {
		legende = ""
	}
	return fmt.Sprintf(`
<section class="accueil-section">
	<div class="carte-presentation">
		<div class="presentation-texte">
			<h2>%s</h2>
			<p>%s</p>
		</div>
		<div class="presentation-photos">
			<figure class="bulle"><img src="/assets/logo.svg" alt=""></figure>
			<figure class="bulle"><img src="%s" alt="%s"><figcaption class="bulle-legende">%s</figcaption></figure>
		</div>
	</div>
	<a class="bouton" href="/equipe">%s</a>
</section>`,
		template.HTMLEscapeString(kernel.Traduit(locale, "accueil.presentation.titre")),
		template.HTMLEscapeString(texte),
		template.HTMLEscapeString(photoPresident),
		legende,
		legende,
		template.HTMLEscapeString(kernel.Traduit(locale, "accueil.presentation.bouton")),
	)
}

// sectionCompetitions est la troisième sous-page : la barre « Compétitions »
// surgit du coin haut gauche et se centre, puis les bulles des compétitions
// apparaissent (grille 2x2 : trois compétitions + s'inscrire). Quand on
// quitte la sous-page, la barre repart d'où elle est venue et les bulles
// disparaissent. Tout est piloté par l'IntersectionObserver du script commun.
func (p *PageAccueil) sectionCompetitions(locale string) string {
	var cellules strings.Builder
	if blocs, err := kernel.ChargerContenus(p.DossierContenu + "/competitions.txt"); err == nil {
		for _, b := range kernel.FiltreSections(blocs, "compet") {
			image := b.Texte("image")
			if image == "" {
				image = "/assets/trophee.svg"
			} else if !strings.HasPrefix(image, "/") {
				image = "/assets/" + image
			}
			cellules.WriteString(fmt.Sprintf(
				`<div class="compet-bulle"><a class="compet-cellule" href="/competition/%s">`+
					`<span class="compet-ronde"><img src="%s" alt=""></span>`+
					`<span class="compet-nom">%s</span></a></div>`,
				template.HTMLEscapeString(slug(b.Texte("nom"))),
				template.HTMLEscapeString(image),
				template.HTMLEscapeString(b.Texte("nom")),
			))
		}
	}
	return fmt.Sprintf(`
<section class="accueil-section compet" id="competitions">
	<div class="compet-barre">%s</div>
	<p class="compet-note">%s</p>
	<div class="compet-grille">
		%s
		<div class="compet-bulle compet-inscrire"><a class="compet-cellule" href="/inscriptions">`+
		`<span class="compet-ronde">%s</span></a></div>
	</div>
</section>`,
		template.HTMLEscapeString(kernel.Traduit(locale, "accueil.compet.barre")),
		template.HTMLEscapeString(kernel.Traduit(locale, "accueil.compet.note")),
		cellules.String(),
		template.HTMLEscapeString(kernel.Traduit(locale, "accueil.compet.inscrire")),
	)
}

// slug transforme un nom de compétition en identifiant d'URL.
func slug(nom string) string {
	var b strings.Builder
	dernierTiret := false
	for _, r := range strings.TrimSpace(strings.ToLower(nom)) {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
			dernierTiret = false
		default:
			if !dernierTiret && b.Len() > 0 {
				b.WriteByte('-')
				dernierTiret = true
			}
		}
	}
	s := strings.Trim(b.String(), "-")
	if s == "" {
		return "compet"
	}
	return s
}
