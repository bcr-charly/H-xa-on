package kernel

import (
	"net/http"
)

// Tables de traduction : toutes les chaînes affichées proviennent d'ici,
// jamais d'une page. Une locale sans traduction retombe sur le français.
var Traductions = map[string]map[string]string{
	"fr": {
		"nav.actus":                   "Actus",
		"nav.planning":                "Planning",
		"nav.equipe":                  "BDS",
		"nav.sport":                   "Sport",
		"nav.competition":             "Compétition",
		"nav.boutique":                "Boutique",
		"nav.connexion":               "Connexion",
		"nav.deconnexion":             "Déconnexion",
		"nav.langue":                  "Langue",
		"accueil.slogan":              "Dans la vie il y a deux types de personnes, les chasseurs et les cueilleurs",
		"actus.titre":                 "Actualités",
		"actus.vide":                  "Aucune actualité pour le moment.",
		"equipe.titre":                "Équipe du BDS",
		"connexion.titre":             "Connexion",
		"connexion.identifiant":       "Identifiant",
		"connexion.motdepasse":        "Mot de passe",
		"connexion.bouton":            "Se connecter",
		"connexion.oublie":            "Mot de passe oublié ?",
		"connexion.erreur":            "Identifiant ou mot de passe incorrect.",
		"pied.description":            "Le Bureau Des Sports de l'ENSIIE organise et promeut les activités sportives pour tous les étudiants de l'école.",
		"pied.liensrapides":           "Liens rapides",
		"pied.informations":           "Informations",
		"pied.contact":                "Contact",
		"pied.faq":                    "FAQ",
		"pied.apropos":                "À propos",
		"pied.actus":                  "Actualités",
		"pied.equipe":                 "Notre équipe",
		"pied.siteecole":              "Site web de l'école",
		"pied.suiveznous":             "Suivez-nous",
		"pied.avenir":                 "à venir",
		"accueil.presentation.titre":  "Qui sommes-nous ?",
		"accueil.sponsors.titre":      "Sponsors",
		"accueil.sponsors.vide":       "Nos sponsors bientôt ici. Un partenariat avec le BDS de l ENSIIE ? Écrivez-nous.",
		"equipe.bureau":               "Le bureau",
		"equipe.membres":              "Les membres",
		"accueil.presentation.bouton": "Découvrir l'équipe du BDS",
		"accueil.compet.barre":        "Compétitions",
		"accueil.compet.inscrire":     "S'inscrire",
		"accueil.bouton.actus":        "Voir les actualités",
		"accueil.association.titre":   "L'association",
		"accueil.association.p1":      "Le BDS de l'ENSIIE, c'est une équipe dynamique de %d membres engagés pour offrir à tous les élèves un accès régulier, varié et de qualité au sport au sein de l'école.",
		"accueil.association.p2":      "Le sport occupe une place centrale dans la vie étudiante à l'ENSIIE. Que vous soyez un athlète aguerri ou que vous ayez simplement envie de bouger, le BDS est là pour vous accompagner toute l'année dans la pratique de vos sports favoris.",
		"accueil.association.photo":   "Photo à venir",
		"accueil.stats.membres":       "Membres du BDS",
		"accueil.stats.actus":         "Actualités publiées",
		"accueil.stats.benevoles":     "Bénévoles",
		"accueil.stats.objectif":      "Objectif : le sport",
		"accueil.competitions.titre":  "Les compétitions",
		"accueil.competitions.texte":  "Chaque année, nos sportifs participent aux tournois inter-écoles partout en France. À chaque déplacement, les délégations portent fièrement les couleurs de l'ENSIIE et de son BDS.",
		"accueil.competitions.bouton": "S'inscrire",
		"accueil.sports.titre":        "Sports",
		"compet.retour":               "Retour à l'accueil",
		"compet.titre":                "Compétitions",
		"compet.vide":                 "Aucune compétition pour le moment.",
		"inscriptions.titre":          "Inscriptions aux compétitions",
		"inscriptions.note":           "Les inscriptions ouvriront prochainement. En attendant, suivez les actus du BDS pour ne rien manquer.",
		"sport.inscrire":              "S'inscrire au %s",
		"sport.horaires":              "Horaires",
		"sport.lieu":                  "Lieu",
		"sport.respo":                 "Respo",
		"sport.contact":               "Contact",
		"sport.competitions":          "Compétitions",
		"langue.nom.fr":               "Français",
		"langue.nom.en":               "English",
	},
	"en": {
		"nav.actus":                   "News",
		"nav.planning":                "Schedule",
		"nav.equipe":                  "Team",
		"nav.sport":                   "Sports",
		"nav.competition":             "Competition",
		"nav.boutique":                "Shop",
		"nav.connexion":               "Log in",
		"nav.deconnexion":             "Log out",
		"nav.langue":                  "Language",
		"accueil.slogan":              "In life there are two kinds of people, hunters and gatherers",
		"actus.titre":                 "News",
		"actus.vide":                  "No news yet.",
		"equipe.titre":                "BDS Team",
		"connexion.titre":             "Log in",
		"connexion.identifiant":       "Username",
		"connexion.motdepasse":        "Password",
		"connexion.bouton":            "Log in",
		"connexion.oublie":            "Forgot password?",
		"connexion.erreur":            "Wrong username or password.",
		"pied.description":            "The Sports Bureau of ENSIIE organises and promotes sports activities for all the students of the school.",
		"pied.liensrapides":           "Quick links",
		"pied.informations":           "Information",
		"pied.contact":                "Contact",
		"pied.faq":                    "FAQ",
		"pied.apropos":                "About",
		"pied.actus":                  "News",
		"pied.equipe":                 "Our team",
		"pied.siteecole":              "School website",
		"pied.suiveznous":             "Follow us",
		"pied.avenir":                 "coming soon",
		"accueil.presentation.titre":  "Who are we?",
		"accueil.sponsors.titre":      "Sponsors",
		"accueil.sponsors.vide":       "Our sponsors soon in this space. Partnership with the ENSIIE BDS? Get in touch.",
		"equipe.bureau":               "The board",
		"equipe.membres":              "The members",
		"accueil.presentation.bouton": "Meet the BDS team",
		"accueil.compet.barre":        "Competitions",
		"accueil.compet.inscrire":     "Sign up",
		"accueil.bouton.actus":        "Read the news",
		"accueil.association.titre":   "The association",
		"accueil.association.p1":      "The ENSIIE BDS is a dynamic team of %d members committed to providing all students with regular, varied, and high-quality access to sport within the school.",
		"accueil.association.p2":      "Sport plays a central role in student life at ENSIIE. Whether you are a seasoned athlete or just want to stay active, the BDS is here to support you all year long in practicing your favorite sports.",
		"accueil.association.photo":   "Photo coming soon",
		"accueil.stats.membres":       "BDS members",
		"accueil.stats.actus":         "News published",
		"accueil.stats.benevoles":     "Volunteers",
		"accueil.stats.objectif":      "One goal: sport",
		"accueil.competitions.titre":  "Competitions",
		"accueil.competitions.texte":  "Every year our athletes take part in inter-school tournaments all over France, proudly representing the colours of ENSIIE and its BDS.",
		"accueil.competitions.bouton": "Sign up",
		"accueil.sports.titre":        "Sports",
		"compet.retour":               "Back to home",
		"compet.titre":                "Competitions",
		"compet.vide":                 "No competition yet.",
		"inscriptions.titre":          "Competition sign-up",
		"inscriptions.note":           "Sign-ups will open soon. Follow the BDS news in the meantime.",
		"sport.inscrire":              "Sign up for %s",
		"sport.horaires":              "Schedule",
		"sport.lieu":                  "Venue",
		"sport.respo":                 "Manager",
		"sport.contact":               "Contact",
		"sport.competitions":          "Competitions",
		"langue.nom.fr":               "French",
		"langue.nom.en":               "English",
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
