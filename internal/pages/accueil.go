package pages

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/bcr-charly/H-xa-on/internal/kernel"
)

// PageAccueil reproduit l'agencement du site du BDS Télécom Paris avec
// la DA braise du BDS ENSIIE : hero plein écran, « L'association »,
// bandeau de statistiques, « Compétitions », événements phares,
// « Les clubs », « Nos partenaires ».
type PageAccueil struct {
	Noyau          *kernel.Noyau
	DossierContenu string
}

func (p *PageAccueil) Nom() string       { return "accueil" }
func (p *PageAccueil) Chemins() []string { return []string{"/accueil"} }

func (p *PageAccueil) blocs(fichier, section string) []kernel.Contenu {
	blocs, err := kernel.ChargerContenus(p.DossierContenu + "/" + fichier)
	if err != nil {
		return nil
	}
	return kernel.FiltreSections(blocs, section)
}

func (p *PageAccueil) Servir(w http.ResponseWriter, r *http.Request) {
	locale := kernel.LocaleDepuisRequete(r)
	connecte := p.Noyau.CompteCourant(r) != ""
	t := kernel.Traduit

	hero := fmt.Sprintf(`
<section class="accueil-hero">
	<div class="accueil-contenu">
		<img class="logo-grand" src="/assets/logo.png?v=2026-11-b" alt="%s">
		<h2 class="accueil-titre">%s</h2>
		<div class="accueil-boutons">
			<a class="bouton" href="/adhesion">%s</a>
			<a class="bouton bouton-borde" href="/sports">%s</a>
		</div>
	</div>
</section>`,
		template.HTMLEscapeString(p.Noyau.Config.Nom),
		template.HTMLEscapeString(t(locale, "accueil.titre")),
		template.HTMLEscapeString(t(locale, "nav.adhesion")),
		template.HTMLEscapeString(t(locale, "nav.sports")),
	)

	// Section « L'association » : deux paragraphes + bouton vers l'équipe.
	association := fmt.Sprintf(`
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
</section>`,
		template.HTMLEscapeString(t(locale, "accueil.association.titre")),
		template.HTMLEscapeString(t(locale, "accueil.association.photo")),
		fmt.Sprintf(template.HTMLEscapeString(t(locale, "accueil.association.p1")), len(p.blocs("equipe.txt", "membre"))),
		template.HTMLEscapeString(t(locale, "accueil.association.p2")),
		template.HTMLEscapeString(t(locale, "accueil.presentation.bouton")),
	)

	// Bandeau de statistiques : membres, sports, séances, adhérents.
	// Valeurs réelles : contenu/statistiques.txt, sinon décomptes réels
	// des contenus (équipe, sports) ; jamais inventées.
	var membres, seances, adherents string
	for _, b := range p.blocs("statistiques.txt", "statistiques") {
		membres = b.Texte("membres")
		seances = b.Texte("seances")
		adherents = b.Texte("adherents")
	}
	nMembres := membres
	if nMembres == "" {
		nMembres = fmt.Sprintf("%d", len(p.blocs("equipe.txt", "membre")))
	}
	nSports := fmt.Sprintf("%d", len(p.blocs("sports.txt", "sport")))

	var stats strings.Builder
	stats.WriteString(`<section class="accueil-section accueil-section-stats"><div class="stats-grille">`)
	stats.WriteString(statCarte(nMembres, t(locale, "accueil.stats.membres")))
	stats.WriteString(statCarte(nSports, t(locale, "accueil.stats.sports")))
	if seances != "" {
		stats.WriteString(statCarte(seances, t(locale, "accueil.stats.seances")))
	}
	if adherents != "" {
		stats.WriteString(statCarte(adherents, t(locale, "accueil.stats.adherents")))
	}
	stats.WriteString(`</div></section>`)

	// Section « Compétitions ».
	competitions := fmt.Sprintf(`
<section class="accueil-section accueil-section-alt" id="competitions">
	<h2>%s</h2>
	<p class="accueil-texte">%s</p>
	<p class="accueil-texte">%s</p>
	<a class="bouton" href="/competitions">%s</a>
</section>`,
		template.HTMLEscapeString(t(locale, "accueil.competitions.titre")),
		template.HTMLEscapeString(t(locale, "accueil.competitions.p1")),
		template.HTMLEscapeString(t(locale, "accueil.competitions.p2")),
		template.HTMLEscapeString(t(locale, "accueil.competitions.bouton")),
	)

	// Section événements phares (répétée par événement, masquée si vide).
	var evenements strings.Builder
	for _, e := range p.blocs("evenements.txt", "evenement") {
		var bouton string
		if site := strings.TrimSpace(e.Texte("site")); site != "" {
			bouton = fmt.Sprintf(`<a class="bouton" href="%s" rel="noopener">%s</a></div>`,
				template.HTMLEscapeString(site),
				template.HTMLEscapeString(fmt.Sprintf(t(locale, "accueil.evenement.site"), e.Texte("nom"))))
		} else {
			bouton = `</div>`
		}
		evenements.WriteString(fmt.Sprintf(`
<section class="accueil-section accueil-evenement">
	<h2>%s</h2>
	<div class="evenement-grille"><div class="association-photo grande"><span>%s</span></div><div class="association-texte">
	<p>%s</p>%s</section>`,
			template.HTMLEscapeString(e.Texte("nom")),
			template.HTMLEscapeString(t(locale, "accueil.evenement.photo")),
			template.HTMLEscapeString(strings.TrimSpace(e.Texte("description"))),
			bouton,
		))
	}

	// Section « Les clubs ».
	var clubs strings.Builder
	blocsClubs := p.blocs("clubs.txt", "club")
	if len(blocsClubs) > 0 {
		var cartes strings.Builder
		for _, c := range blocsClubs {
			cartes.WriteString(fmt.Sprintf(
				`<article class="carte club-carte"><div class="carte-corps"><h3>%s</h3><p>%s</p></div></article>`,
				template.HTMLEscapeString(c.Texte("nom")),
				template.HTMLEscapeString(strings.TrimSpace(c.Texte("description"))),
			))
		}
		clubs.WriteString(fmt.Sprintf(`
<section class="accueil-section" id="clubs">
	<h2>%s</h2>
	<p class="accueil-texte">%s</p>
	<div class="clubs-grille">%s</div>
	<a class="bouton" href="/clubs">%s</a>
</section>`,
			template.HTMLEscapeString(t(locale, "clubs.titre")),
			template.HTMLEscapeString(t(locale, "accueil.clubs.intro")),
			cartes.String(),
			template.HTMLEscapeString(t(locale, "accueil.clubs.bouton")),
		))
	}

	// Section « Nos partenaires » (auto-masquée si aucun partenaire).
	var partenaires strings.Builder
	blocsPartenaires := p.blocs("partenaires.txt", "partenaire")
	if len(blocsPartenaires) > 0 {
		var cartes strings.Builder
		for _, c := range blocsPartenaires {
			var offre strings.Builder
			for _, l := range strings.Split(c.Texte("offre"), "\n") {
				l = strings.TrimSpace(l)
				if l == "" {
					continue
				}
				if strings.HasPrefix(l, "- ") {
					offre.WriteString(fmt.Sprintf("<li>%s</li>", template.HTMLEscapeString(strings.TrimSpace(l[2:]))))
				} else {
					offre.WriteString(fmt.Sprintf("<p>%s</p>", template.HTMLEscapeString(l)))
				}
			}
			var puces string
			if strings.Contains(offre.String(), "<li>") {
				puces = `<ul class="partenaire-puces">` + offre.String() + `</ul>`
			} else {
				puces = offre.String()
			}
			var conditions string
			if c.Texte("conditions") == "oui" {
				conditions = fmt.Sprintf(`<p class="partenaire-conditions">%s</p>`,
					template.HTMLEscapeString(t(locale, "partenaire.conditions")))
			}
			var boutons string
			if site := strings.TrimSpace(c.Texte("site")); site != "" {
				boutons = fmt.Sprintf(`<a class="bouton" href="%s" rel="noopener">%s</a>`,
					template.HTMLEscapeString(site),
					template.HTMLEscapeString(t(locale, "partenaire.site")))
			}
			if flyer := strings.TrimSpace(c.Texte("flyer")); flyer != "" {
				boutons += fmt.Sprintf(` <a class="bouton bouton-borde" href="%s" rel="noopener">%s</a>`,
					template.HTMLEscapeString(flyer),
					template.HTMLEscapeString(t(locale, "partenaire.flyer")))
			}
			cartes.WriteString(fmt.Sprintf(
				`<article class="carte partenaire-carte"><div class="carte-corps"><h3>%s</h3>%s%s<div class="accueil-boutons">%s</div></div></article>`,
				template.HTMLEscapeString(c.Texte("nom")),
				puces,
				conditions,
				boutons,
			))
		}
		partenaires.WriteString(fmt.Sprintf(`
<section class="accueil-section accueil-section-alt" id="partenaires">
	<h2>%s</h2>
	<p class="accueil-texte">%s</p>
	<div class="cartes-colonne">%s</div>
</section>`,
			template.HTMLEscapeString(t(locale, "accueil.partenaires.titre")),
			template.HTMLEscapeString(t(locale, "accueil.partenaires.intro")),
			cartes.String(),
		))
	}

	contenu := template.HTML(hero + association + stats.String() + competitions +
		evenements.String() + clubs.String() + partenaires.String())
	p.Noyau.Rendu(w, http.StatusOK, page(p.Noyau, locale, p.Nom(), p.Noyau.Config.Nom, contenu, connecte), nil)
}

// statCarte rend un compteur du bandeau de statistiques.
func statCarte(chiffre, libelle string) string {
	return fmt.Sprintf(`<div class="stat"><div class="stat-chiffre">%s</div><div class="stat-libelle">%s</div></div>`,
		template.HTMLEscapeString(chiffre), template.HTMLEscapeString(libelle))
}
