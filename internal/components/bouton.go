package components

import (
	"fmt"
	"html/template"
)

// Bouton rend un bouton d'action. Comportement commun à tous les boutons du
// site : grossissement léger au survol sur ordinateur, animation équivalente
// au toucher.
func (b *Bibliotheque) Bouton(libelle, lien string) template.HTML {
	h := fmt.Sprintf(
		`<a class="bouton" href="%s">%s</a>`,
		template.HTMLEscapeString(lien),
		template.HTMLEscapeString(libelle),
	)
	return template.HTML(h)
}

// Lien rend un simple lien texte stylé.
func (b *Bibliotheque) Lien(libelle, lien string) template.HTML {
	h := fmt.Sprintf(
		`<a class="lien" href="%s">%s</a>`,
		template.HTMLEscapeString(lien),
		template.HTMLEscapeString(libelle),
	)
	return template.HTML(h)
}
