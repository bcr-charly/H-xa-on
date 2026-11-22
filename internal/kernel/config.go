package kernel

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config est la configuration complète d'une instance Héxaïon,
// lue depuis conf.txt au démarrage.
type Config struct {
	Port    int
	Nom     string
	Palette Palette
	Pages   map[string]bool
	Secret  string
}

// Palette contient les couleurs de l'instance. La palette par défaut
// est celle de l'identité visuelle du BDS ENSIIE.
type Palette struct {
	Fond      string
	Elevation string
	Clair     string
	Accent    string
	Bordure   string
}

func PaletteDefaut() Palette {
	return Palette{
		Fond:      "#0B0A07",
		Elevation: "#1F1A0C",
		Clair:     "#FFC61A",
		Accent:    "#E8A400",
		Bordure:   "#8C6400",
	}
}

// Avertissement est un problème de configuration non fatal.
type Avertissement struct {
	Section string
	Cle     string
	Message string
}

func (a Avertissement) String() string {
	if a.Section == "" {
		return fmt.Sprintf("configuration : %s", a.Message)
	}
	return fmt.Sprintf("configuration [%s] %s : %s", a.Section, a.Cle, a.Message)
}

// ParseConfig analyse un fichier conf.txt au format texte à indentation.
func ParseConfig(contenu string) (Config, []Avertissement) {
	cfg := Config{
		Port:    8080,
		Nom:     "Héxaïon",
		Palette: PaletteDefaut(),
		Pages:   map[string]bool{},
		Secret:  "",
	}
	averts := []Avertissement{}

	section := ""
	for numero, ligne := range strings.Split(contenu, "\n") {
		l := strings.TrimSpace(ligne)
		if l == "" || strings.HasPrefix(l, "#") {
			continue
		}
		if !strings.HasPrefix(ligne, " ") && !strings.HasPrefix(ligne, "\t") {
			section = l
			switch section {
			case "serveur", "palette", "pages", "securite":
			default:
				averts = append(averts, Avertissement{Message: fmt.Sprintf("ligne %d : section inconnue %q, ignorée", numero+1, section)})
			}
			continue
		}
		if section == "" {
			averts = append(averts, Avertissement{Message: fmt.Sprintf("ligne %d : clé hors section, ignorée", numero+1)})
			continue
		}
		eq := strings.Index(l, "=")
		if eq < 0 {
			averts = append(averts, Avertissement{Section: section, Message: fmt.Sprintf("ligne %d : ligne sans '=', ignorée", numero+1)})
			continue
		}
		cle := strings.TrimSpace(l[:eq])
		val := strings.TrimSpace(l[eq+1:])

		switch section {
		case "serveur":
			switch cle {
			case "port":
				if p, err := strconv.Atoi(val); err == nil && p > 0 && p < 65536 {
					cfg.Port = p
				} else {
					averts = append(averts, Avertissement{Section: section, Cle: cle, Message: fmt.Sprintf("valeur %q invalide, port par défaut conservé", val)})
				}
			case "nom":
				cfg.Nom = val
			default:
				averts = append(averts, Avertissement{Section: section, Cle: cle, Message: "clé inconnue, ignorée"})
			}
		case "palette":
			p := cfg.Palette
			cible := &p
			switch cle {
			case "fond":
				cible.Fond = val
			case "elevation":
				cible.Elevation = val
			case "clair":
				cible.Clair = val
			case "accent":
				cible.Accent = val
			case "bordure":
				cible.Bordure = val
			default:
				averts = append(averts, Avertissement{Section: section, Cle: cle, Message: "clé inconnue, ignorée"})
				continue
			}
			cfg.Palette = p
			if !couleurValide(val) {
				averts = append(averts, Avertissement{Section: section, Cle: cle, Message: fmt.Sprintf("valeur %q n'est pas une couleur #RRGGBB", val)})
			}
		case "pages":
			on, err := parseONOFF(val)
			if err != nil {
				averts = append(averts, Avertissement{Section: section, Cle: cle, Message: fmt.Sprintf("valeur %q : attendre ON ou OFF", val)})
				continue
			}
			cfg.Pages[cle] = on
		case "securite":
			switch cle {
			case "secret":
				cfg.Secret = val
			default:
				averts = append(averts, Avertissement{Section: section, Cle: cle, Message: "clé inconnue, ignorée"})
			}
		}
	}
	return cfg, averts
}

func parseONOFF(v string) (bool, error) {
	switch strings.ToUpper(v) {
	case "ON":
		return true, nil
	case "OFF":
		return false, nil
	}
	return false, fmt.Errorf("valeur invalide")
}

func couleurValide(v string) bool {
	if len(v) != 7 || v[0] != '#' {
		return false
	}
	for _, c := range v[1:] {
		if !strings.ContainsRune("0123456789abcdefABCDEF", c) {
			return false
		}
	}
	return true
}

// ChargerConfig lit conf.txt depuis le disque.
func ChargerConfig(chemin string) (Config, []Avertissement, error) {
	data, err := os.ReadFile(chemin)
	if err != nil {
		if os.IsNotExist(err) {
			// Une instance sans conf.txt démarre avec les valeurs par défaut.
			cfg, a := ParseConfig("")
			return cfg, a, nil
		}
		return Config{}, nil, err
	}
	cfg, a := ParseConfig(string(data))
	return cfg, a, nil
}
