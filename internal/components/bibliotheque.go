package components

import "html/template"

// Bibliotheque est la bibliothèque de composants réutilisables : toute page
// s'exprime avec ce vocabulaire. Les composants prennent des chaînes en
// entrée et retournent du HTML sûr : les valeurs sont échappées
// systématiquement via html/template.
type Bibliotheque struct {
	funcs template.FuncMap
}

// Defaut est la bibliothèque enregistrée pour tous les rendus du noyau.
var Defaut = nouvelle()

func nouvelle() *Bibliotheque {
	b := &Bibliotheque{funcs: template.FuncMap{}}
	b.funcs["bouton"] = b.Bouton
	b.funcs["carte"] = b.Carte
	b.funcs["iconeRonde"] = b.IconeRonde
	b.funcs["bandeTitre"] = b.BandeTitre
	b.funcs["lien"] = b.Lien
	return b
}

// FuncMap expose les composants aux templates des pages.
func (b *Bibliotheque) FuncMap() template.FuncMap { return b.funcs }
