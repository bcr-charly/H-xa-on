package components

import (
	"fmt"
	"html/template"
)

// Carte rend une carte d'actualité : palette d'élévation en fond, bordure
// bronze. L'image est optionnelle ; une date vide masque la ligne de date.
func (b *Bibliotheque) Carte(titre, date, texte, image string) template.HTML {
	var img string
	if image != "" {
		img = fmt.Sprintf(`<img class="carte-image" src="%s" alt="%s" loading="lazy">`,
			template.HTMLEscapeString(image), template.HTMLEscapeString(titre))
	}
	var ligneDate string
	if date != "" {
		ligneDate = fmt.Sprintf(`<div class="carte-date">%s</div>`, template.HTMLEscapeString(date))
	}
	h := fmt.Sprintf(
		`<article class="carte">%s<div class="carte-corps"><h3>%s</h3>%s<p>%s</p></div></article>`,
		img,
		template.HTMLEscapeString(titre),
		ligneDate,
		template.HTMLEscapeString(texte),
	)
	return template.HTML(h)
}

// BandeTitre rend la grande bande de titre de page.
func (b *Bibliotheque) BandeTitre(titre string) template.HTML {
	return template.HTML(fmt.Sprintf(`<div class="bande-titre"><h1>%s</h1></div>`,
		template.HTMLEscapeString(titre)))
}
