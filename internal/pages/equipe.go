package pages

import (
	"html/template"
	"math"
	"net/http"
	"strings"

	"github.com/bcr-charly/siteweb-bds-ensiie/internal/components"
	"github.com/bcr-charly/siteweb-bds-ensiie/internal/kernel"
)

// PageEquipe affiche la grille des membres du bureau, remplie depuis
// contenu/equipe.txt (nom, rôle, photo, lien éventuel).
//
// La disposition est symétrique : d'abord la rangée du bureau
// (président, vice-président, trésorier), puis des rangées régulières
// centrées — le nombre de colonnes est calculé pour que la grille reste
// équilibrée quelle que soit l'effectif (19 membres : 3 puis 4×4).
type PageEquipe struct {
	Noyau          *kernel.Noyau
	DossierContenu string
}

func (p *PageEquipe) Nom() string       { return "equipe" }
func (p *PageEquipe) Chemins() []string { return []string{"/equipe"} }

func (p *PageEquipe) Servir(w http.ResponseWriter, r *http.Request) {
	locale := kernel.LocaleDepuisRequete(r)
	connecte := p.Noyau.CompteCourant(r) != ""

	blocs, err := kernel.ChargerContenus(p.DossierContenu + "/equipe.txt")
	if err != nil {
		http.Error(w, "contenu illisible", http.StatusInternalServerError)
		return
	}

	var bureau, autres []kernel.Contenu
	for _, b := range kernel.FiltreSections(blocs, "membre") {
		if b.Texte("rang") == "bureau" {
			bureau = append(bureau, b)
		} else {
			autres = append(autres, b)
		}
	}

	var grille strings.Builder
	if len(bureau) > 0 {
		grille.WriteString(`<div class="rangee-equipe">`)
		for _, b := range bureau {
			grille.WriteString(string(components.Defaut.IconeRonde(
				cheminImage(b.Texte("photo")),
				b.Texte("nom"),
				b.Texte("role"),
				b.Texte("lien"),
			)))
		}
		grille.WriteString(`</div>`)
	}
	if len(autres) > 0 {
		colonnes := colonnesSymetriques(len(autres))
		for i, b := range autres {
			if i%colonnes == 0 {
				grille.WriteString(`<div class="rangee-equipe">`)
			}
			grille.WriteString(string(components.Defaut.IconeRonde(
				cheminImage(b.Texte("photo")),
				b.Texte("nom"),
				b.Texte("role"),
				b.Texte("lien"),
			)))
			if (i+1)%colonnes == 0 || i == len(autres)-1 {
				grille.WriteString(`</div>`)
			}
		}
	}

	contenu := template.HTML(string(components.Defaut.BandeTitre(kernel.Traduit(locale, "equipe.titre"))) +
		`<div class="conteneur"><div class="grille-equipe">` + grille.String() + `</div></div>`)
	p.Noyau.Rendu(w, http.StatusOK, page(p.Noyau, locale, p.Nom(), kernel.Traduit(locale, "equipe.titre"), contenu, connecte), nil)
}

// colonnesSymetriques choisit un nombre de colonnes qui rend la grille
// régulière : racine carrée arrondie au supérieur (16 membres → 4 colonnes).
func colonnesSymetriques(n int) int {
	if n <= 0 {
		return 1
	}
	c := int(math.Ceil(math.Sqrt(float64(n))))
	if c < 2 {
		c = 2
	}
	return c
}

func cheminImage(photo string) string {
	if photo == "" {
		return "/assets/membre.svg"
	}
	if strings.HasPrefix(photo, "/") {
		return photo
	}
	return "/assets/" + photo
}
