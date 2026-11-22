package pages

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/bcr-charly/H-xa-on/internal/components"
	"github.com/bcr-charly/H-xa-on/internal/kernel"
)

// PageSport est la page dédiée d'un sport, servie sous /sport/<slug>,
// remplie depuis contenu/sports.txt. Chaque bloc [sport] supporte :
//
//	nom           — nom affiché
//	image         — photo de la case (sous-page sports)
//	presentation  — présentation du sport au BDS (multi-lignes)
//	competitions  — une compétition par ligne indentée
//	respo         — responsable du sport
//	contact       — contacts du respo (multi-lignes)
//	photo, photo2 — photos de la page
//	horaires      — créneaux
//	lieu          — où se déroulent les séances
//
// Le bouton « S'inscrire au ... » redirige vers l'onglet inscription.
// Le parcours complet (création de compte, cotisation HelloAsso) est prévu
// dans la feuille de route : il sera branché ici quand le compte BDS existera.
type PageSport struct {
	Noyau          *kernel.Noyau
	DossierContenu string
}

func (p *PageSport) Nom() string       { return "sport" }
func (p *PageSport) Chemins() []string { return []string{"/sport/"} }

func (p *PageSport) Servir(w http.ResponseWriter, r *http.Request) {
	locale := kernel.LocaleDepuisRequete(r)
	connecte := p.Noyau.CompteCourant(r) != ""
	imageFond := "/assets/sport.svg"

	cle := strings.Trim(strings.TrimPrefix(r.URL.Path, "/sport/"), "/")
	blocs, err := kernel.ChargerContenus(p.DossierContenu + "/sports.txt")
	if err != nil {
		http.Error(w, "contenu illisible", http.StatusInternalServerError)
		return
	}
	var sport *kernel.Contenu
	for _, b := range kernel.FiltreSections(blocs, "sport") {
		if slug(b.Texte("nom")) == cle || (b.Texte("slug") != "" && b.Texte("slug") == cle) {
			c := b
			sport = &c
			break
		}
	}
	if sport == nil {
		http.NotFound(w, r)
		return
	}
	if v := sport.Texte("image"); v != "" {
		if !strings.HasPrefix(v, "/") {
			v = "/assets/" + v
		}
		imageFond = v
	}
	p.servirDetail(w, r, locale, connecte, sport, imageFond)
}

func (p *PageSport) servirDetail(w http.ResponseWriter, r *http.Request, locale string, connecte bool, sport *kernel.Contenu, imageFond string) {
	nom := sport.Texte("nom")

	var competitions strings.Builder
	if lignes := strings.TrimSpace(sport.Texte("competitions")); lignes != "" {
		for _, l := range strings.Split(lignes, "\n") {
			competitions.WriteString(fmt.Sprintf("<li>%s</li>", template.HTMLEscapeString(strings.TrimSpace(l))))
		}
	}
	var blocCompet string
	if competitions.Len() > 0 {
		blocCompet = fmt.Sprintf(`<div class="sport-detail-bloc"><h2>%s</h2><ul>%s</ul></div>`,
			template.HTMLEscapeString(kernel.Traduit(locale, "sport.competitions")),
			competitions.String())
	}

	var photos strings.Builder
	for _, champ := range []string{"photo", "photo2", "photo3"} {
		if v := sport.Texte(champ); v != "" {
			if !strings.HasPrefix(v, "/") {
				v = "/assets/" + v
			}
			photos.WriteString(fmt.Sprintf(`<img src="%s" alt="%s" loading="lazy">`,
				template.HTMLEscapeString(v), template.HTMLEscapeString(nom)))
		}
	}
	var blocPhotos string
	if photos.Len() > 0 {
		blocPhotos = `<div class="sport-photos">` + photos.String() + `</div>`
	}

	var blocRespo string
	if sport.Texte("respo") != "" || sport.Texte("contact") != "" {
		blocRespo = fmt.Sprintf(`<div class="sport-detail-bloc"><h2>%s</h2><p>%s</p><p>%s</p></div>`,
			template.HTMLEscapeString(kernel.Traduit(locale, "sport.respo")),
			template.HTMLEscapeString(sport.Texte("respo")),
			template.HTMLEscapeString(strings.TrimSpace(sport.Texte("contact"))),
		)
	}

	horairesLieu := fmt.Sprintf(`<div class="sport-horaires-lieu">
<div><h2>%s</h2><p>%s</p></div>
<div><h2>%s</h2><p>%s</p></div>
</div>`,
		template.HTMLEscapeString(kernel.Traduit(locale, "sport.horaires")),
		template.HTMLEscapeString(strings.TrimSpace(sport.Texte("horaires"))),
		template.HTMLEscapeString(kernel.Traduit(locale, "sport.lieu")),
		template.HTMLEscapeString(strings.TrimSpace(sport.Texte("lieu"))),
	)

	bouton := fmt.Sprintf(`<div class="sport-inscrire"><a class="bouton" href="/inscriptions">%s</a></div>`,
		template.HTMLEscapeString(fmt.Sprintf(kernel.Traduit(locale, "sport.inscrire"), nom)))

	detail := fmt.Sprintf(`<div class="sport-detail">
<h1>%s</h1>
<div class="sport-detail-intro">%s</div>
%s
%s
%s
%s
%s
</div>`,
		template.HTMLEscapeString(nom),
		template.HTMLEscapeString(strings.TrimSpace(sport.Texte("presentation"))),
		blocCompet,
		blocRespo,
		blocPhotos,
		horairesLieu,
		bouton,
	)
	contenu := template.HTML(fondFixe(imageFond) + string(components.Defaut.BandeTitre(nom)) + `<div class="conteneur">` + detail + `</div>`)
	p.Noyau.Rendu(w, http.StatusOK, page(p.Noyau, locale, p.Nom(), nom, contenu, connecte), nil)
}
