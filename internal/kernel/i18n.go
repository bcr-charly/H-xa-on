package kernel

import (
	"net/http"
)

// Tables de traduction : toutes les chaînes affichées proviennent d'ici,
// jamais d'une page. Une locale sans traduction retombe sur le français.
var Traductions = map[string]map[string]string{
	"fr": {
		"nav.actus":             "Actus",
		"nav.planning":          "Planning",
		"nav.equipe":            "BDS",
		"nav.sport":             "Sport",
		"nav.competition":       "Compétition",
		"nav.boutique":          "Boutique",
		"nav.connexion":         "Connexion",
		"nav.deconnexion":       "Déconnexion",
		"nav.langue":            "Langue",
		"accueil.slogan":        "Dans la vie il y a deux types de personnes, les chasseurs et les cueilleurs",
		"actus.titre":           "Actualités",
		"actus.vide":            "Aucune actualité pour le moment.",
		"equipe.titre":          "Équipe du BDS",
		"connexion.titre":       "Connexion",
		"connexion.identifiant": "Identifiant",
		"connexion.motdepasse":  "Mot de passe",
		"connexion.bouton":      "Se connecter",
		"connexion.oublie":      "Mot de passe oublié ?",
		"connexion.erreur":      "Identifiant ou mot de passe incorrect.",
		"pied.description":      "Le Bureau Des Sports de l'ENSIIE organise et promeut les activités sportives pour tous les étudiants de l'école.",
		"pied.liensrapides":     "Liens rapides",
		"pied.informations":     "Informations",
		"pied.contact":          "Contact",
		"pied.faq":              "FAQ",
		"pied.apropos":          "À propos",
		"pied.actus":            "Actualités",
		"pied.equipe":           "Notre équipe",
		"pied.siteecole":        "Site web de l'école",
		"pied.suiveznous":       "Suivez-nous",
		"pied.avenir":           "à venir",
		"langue.nom.fr":         "Français",
		"langue.nom.en":         "English",
	},
	"en": {
		"nav.actus":             "News",
		"nav.planning":          "Schedule",
		"nav.equipe":            "Team",
		"nav.sport":             "Sports",
		"nav.competition":       "Competition",
		"nav.boutique":          "Shop",
		"nav.connexion":         "Log in",
		"nav.deconnexion":       "Log out",
		"nav.langue":            "Language",
		"accueil.slogan":        "In life there are two kinds of people, hunters and gatherers",
		"actus.titre":           "News",
		"actus.vide":            "No news yet.",
		"equipe.titre":          "BDS Team",
		"connexion.titre":       "Log in",
		"connexion.identifiant": "Username",
		"connexion.motdepasse":  "Password",
		"connexion.bouton":      "Log in",
		"connexion.oublie":      "Forgot password?",
		"connexion.erreur":      "Wrong username or password.",
		"pied.description":      "The Sports Bureau of ENSIIE organises and promotes sports activities for all the students of the school.",
		"pied.liensrapides":     "Quick links",
		"pied.informations":     "Information",
		"pied.contact":          "Contact",
		"pied.faq":              "FAQ",
		"pied.apropos":          "About",
		"pied.actus":            "News",
		"pied.equipe":           "Our team",
		"pied.siteecole":        "School website",
		"pied.suiveznous":       "Follow us",
		"pied.avenir":           "coming soon",
		"langue.nom.fr":         "French",
		"langue.nom.en":         "English",
	},
}

// LocalesDisponibles liste les locales livrées.
func LocalesDisponibles() []string {
	res := []string{}
	for l := range Traductions {
		res = append(res, l)
	}
	return res
}

// Traduit retourne la chaîne traduite pour la locale donnée.
func Traduit(locale, cle string) string {
	if t, ok := Traductions[locale]; ok {
		if v, ok := t[cle]; ok && v != "" {
			return v
		}
	}
	if v, ok := Traductions["fr"][cle]; ok {
		return v
	}
	return cle
}

// LocaleDepuisRequete détermine la locale : paramètre ?lang=, puis cookie,
// puis français par défaut.
func LocaleDepuisRequete(r *http.Request) string {
	if l := r.URL.Query().Get("lang"); l != "" {
		return l
	}
	if c, err := r.Cookie("hexaion_langue"); err == nil && c.Value != "" {
		return c.Value
	}
	return "fr"
}
