package kernel

import (
	"fmt"
	"os"
	"strings"
)

// Contenu est un bloc de fichier de contenu : un ensemble de paires clé-valeur.
type Contenu struct {
	Section string
	Champs  map[string]string
}

func (c Contenu) Texte(cle string) string { return c.Champs[cle] }

// ParseContenus découpe un fichier de contenu en blocs.
//
// Un bloc commence par une ligne "[section]" et se poursuit jusqu'au bloc
// suivant. Chaque ligne du bloc est "cle = valeur". Les lignes vides et les
// commentaires (#) sont ignorés. Une ligne indentée (espace ou tabulation en
// début de ligne) est ajoutée à la valeur de la clé précédente, ce qui permet
// des textes multi-lignes.
func ParseContenus(texte string) ([]Contenu, error) {
	var blocs []Contenu
	derniereCle := ""
	numLigne := 0

	for _, brute := range strings.Split(texte, "\n") {
		numLigne++
		l := strings.TrimSpace(brute)
		if l == "" || strings.HasPrefix(l, "#") {
			continue
		}
		if strings.HasPrefix(l, "[") && strings.HasSuffix(l, "]") {
			blocs = append(blocs, Contenu{Section: l[1 : len(l)-1], Champs: map[string]string{}})
			derniereCle = ""
			continue
		}
		if len(blocs) == 0 {
			return nil, fmt.Errorf("ligne %d : contenu hors bloc [section]", numLigne)
		}
		if brute[0] == ' ' || brute[0] == '\t' {
			if derniereCle == "" {
				return nil, fmt.Errorf("ligne %d : valeur sans clé", numLigne)
			}
			blocs[len(blocs)-1].Champs[derniereCle] += "\n" + l
			continue
		}
		eq := strings.Index(l, "=")
		if eq < 0 {
			return nil, fmt.Errorf("ligne %d : ligne sans '='", numLigne)
		}
		derniereCle = strings.TrimSpace(l[:eq])
		blocs[len(blocs)-1].Champs[derniereCle] = strings.TrimSpace(l[eq+1:])
	}
	return blocs, nil
}

// ChargerContenus lit et analyse un fichier de contenu.
// Un fichier absent n'est pas une erreur : la page s'affiche vide.
func ChargerContenus(chemin string) ([]Contenu, error) {
	data, err := os.ReadFile(chemin)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	return ParseContenus(string(data))
}

// FiltreSections ne garde que les blocs d'une section donnée.
func FiltreSections(blocs []Contenu, section string) []Contenu {
	var res []Contenu
	for _, b := range blocs {
		if b.Section == section {
			res = append(res, b)
		}
	}
	return res
}
