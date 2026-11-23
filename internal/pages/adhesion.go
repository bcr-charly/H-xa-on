package pages

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/bcr-charly/siteweb-bds-ensiie/internal/components"
	"github.com/bcr-charly/siteweb-bds-ensiie/internal/kernel"
)

// PageAdhesion est la page « L'adhésion » : pourquoi devenir membre, ce
// que l'adhésion comprend et comment adhérer. Le contenu vient de
// contenu/adhesion.txt ; le bouton d'adhésion externe (site de paiement
// de l'association) n'apparaît que si le champ lien est renseigné.
type PageAdhesion struct {
	Noyau          *kernel.Noyau
	DossierContenu string
}

func (p *PageAdhesion) Nom() string       { return "adhesion" }
func (p *PageAdhesion) Chemins() []string { return []string{"/adhesion"} }

func (p *PageAdhesion) Servir(w http.ResponseWriter, r *http.Request) {
	locale := kernel.LocaleDepuisRequete(r)
	connecte := p.Noyau.CompteCourant(r) != ""
	t := kernel.Traduit

	blocs, err := kernel.ChargerContenus(p.DossierContenu + "/adhesion.txt")
	if err != nil {
		http.Error(w, "contenu illisible", http.StatusInternalServerError)
		return
	}
	var conf *kernel.Contenu
	for _, b := range kernel.FiltreSections(blocs, "adhesion") {
		c := b
		conf = &c
		break
	}

	var pourquoi, comprend, comment, lien string
	if conf != nil {
		pourquoi = conf.Texte("pourquoi")
		comprend = conf.Texte("comprend")
		comment = conf.Texte("comment")
		lien = strings.TrimSpace(conf.Texte("lien"))
	}

	var puces strings.Builder
	if strings.TrimSpace(comprend) != "" {
		puces.WriteString(`<ul class="adhesion-puces">`)
		for _, l := range strings.Split(comprend, "\n") {
			l = strings.TrimSpace(l)
			if strings.HasPrefix(l, "- ") {
				puces.WriteString(fmt.Sprintf("<li>%s</li>", template.HTMLEscapeString(strings.TrimSpace(l[2:]))))
			}
		}
		puces.WriteString(`</ul>`)
	}

	var blocLien string
	if lien != "" {
		blocLien = fmt.Sprintf(`<div id="adhesion-payer"><a class="bouton" href="%s" rel="noopener">%s</a></div>`,
			template.HTMLEscapeString(lien),
			template.HTMLEscapeString(t(locale, "adhesion.bouton")))
	}

	corps := fmt.Sprintf(`
<div class="adhesion-section"><h2>%s</h2><p>%s</p></div>
<div class="adhesion-section"><h2>%s</h2>%s</div>
<div class="adhesion-section"><h2>%s</h2><p>%s</p>%s</div>`,
		template.HTMLEscapeString(t(locale, "adhesion.pourquoi.titre")),
		template.HTMLEscapeString(strings.TrimSpace(pourquoi)),
		template.HTMLEscapeString(t(locale, "adhesion.comprend.titre")),
		puces.String(),
		template.HTMLEscapeString(t(locale, "adhesion.comment.titre")),
		template.HTMLEscapeString(strings.TrimSpace(comment)),
		blocLien,
	)

	contenu := template.HTML(string(components.Defaut.BandeTitre(t(locale, "adhesion.titre"))) +
		`<div class="conteneur adhesion">` + corps + `</div>`)
	p.Noyau.Rendu(w, http.StatusOK, page(p.Noyau, locale, p.Nom(), t(locale, "adhesion.titre"), contenu, connecte), nil)
}
