package components

import (
	"fmt"
	"html/template"
)

// IconeRonde rend une icône ronde contenant la photo du membre : grossissement
// au survol sur ordinateur, léger agrandissement animé au toucher sur mobile.
// Le lien est optionnel.
func (b *Bibliotheque) IconeRonde(image, nom, role, lien string) template.HTML {
	contenu := fmt.Sprintf(
		`<img class="icone-image" src="%s" alt="%s"><div class="icone-nom">%s</div><div class="icone-role">%s</div>`,
		template.HTMLEscapeString(image),
		template.HTMLEscapeString(nom),
		template.HTMLEscapeString(nom),
		template.HTMLEscapeString(role),
	)
	var h string
	if lien != "" {
		h = fmt.Sprintf(`<a class="icone-ronde" href="%s">%s</a>`,
			template.HTMLEscapeString(lien), contenu)
	} else {
		h = fmt.Sprintf(`<div class="icone-ronde">%s</div>`, contenu)
	}
	return template.HTML(h)
}
