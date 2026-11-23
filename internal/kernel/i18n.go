package kernel

import (
	"net/http"
)

// Tables de traduction : le site est monolingue français. Le mécanisme de
// traduction est conservé pour qu'une future variante linguistique puisse
// s'ajouter sans toucher aux pages : toutes les chaînes affichées
// proviennent d'ici, jamais d'une page.
var Traductions = map[string]map[string]string{
	"fr": {
		"nav.accueil":                 "L'association",
		"nav.sports":                  "Nos sports",
		"nav.competitions":            "Compétitions",
		"nav.clubs":                   "Clubs",
		"nav.actus":                   "Actualités",
		"nav.adhesion":                "Adhésion",
		"nav.equipe":                  "Équipe",
		"nav.planning":                "Planning",
		"nav.boutique":                "Boutique",
		"nav.connexion":               "Connexion",
		"nav.deconnexion":             "Déconnexion",
		"nav.inscription":             "Inscription",
		"nav.langue":                  "Langue",
		"accueil.titre":               "L'association qui fait vivre le sport à l'ENSIIE",
		"accueil.slogan":              "Dans la vie il y a deux types de personnes, les chasseurs et les cueilleurs",
		"accueil.association.titre":   "L'association",
		"accueil.association.p1":      "Le BDS de l'ENSIIE, c'est une équipe dynamique de %d membres engagés pour offrir à tous les élèves un accès régulier, varié et de qualité au sport au sein de l'école.",
		"accueil.association.p2":      "Le sport occupe une place centrale dans la vie étudiante à l'ENSIIE. Que vous soyez un athlète aguerri ou que vous ayez simplement envie de bouger, le BDS est là pour vous accompagner toute l'année dans la pratique de vos sports favoris.",
		"accueil.association.photo":   "Photo à venir",
		"accueil.presentation.bouton": "Rencontre l'équipe",
		"accueil.stats.membres":       "Membres du BDS",
		"accueil.stats.sports":        "Sports proposés",
		"accueil.stats.seances":       "Séances hebdomadaires",
		"accueil.stats.adherents":     "Adhérents",
		"accueil.competitions.titre":  "Compétitions",
		"accueil.competitions.p1":     "Chaque année, nos sportifs participent aux tournois inter-écoles partout en France. À chaque déplacement, les délégations portent fièrement les couleurs de l'ENSIIE et de son BDS.",
		"accueil.competitions.p2":     "Compétiteurs comme supporters, ces événements sont des moments à ne pas manquer autour des valeurs du sport.",
		"accueil.competitions.bouton": "Plus d'infos",
		"accueil.evenement.site":      "Site de %s",
		"accueil.evenement.photo":     "Photo à venir",
		"accueil.clubs.intro":         "Le BDS apporte son soutien financier à plusieurs associations étudiantes spécialisées.",
		"accueil.clubs.bouton":        "Voir tous les clubs",
		"accueil.partenaires.titre":   "Nos partenaires",
		"accueil.partenaires.intro":   "Ils soutiennent le BDS toute l'année. Découvrez leurs offres pour les étudiants de l'ENSIIE.",
		"partenaire.conditions":       "Offres soumises à conditions.",
		"partenaire.site":             "Visiter le site",
		"partenaire.flyer":            "Voir le flyer",
		"sports.titre":                "Nos sports",
		"sports.intro":                "Découvre les disciplines proposées par le BDS de l'ENSIIE.",
		"sports.compte":               "%d sports",
		"sports.unite":                "sports",
		"sports.filtre.tous":          "Tous les sports",
		"sports.decouvrir":            "Découvrir",
		"sports.retour":               "Retour aux sports",
		"sports.vide":                 "Aucun sport pour le moment.",
		"actus.titre":                 "Les actualités",
		"actus.intro":                 "Tout sur les évènements et les activités récents du BDS de l'ENSIIE.",
		"actus.vide":                  "Aucun article trouvé. Modifiez votre recherche.",
		"actus.recherche":             "Rechercher une actualité…",
		"compet.titre":                "Compétitions",
		"compet.intro":                "Tout sur la compétition à l'ENSIIE, des tournois inter-écoles aux championnats FFSU.",
		"compet.vide":                 "Aucune compétition pour le moment.",
		"clubs.titre":                 "Les clubs et leur financement",
		"clubs.intro":                 "Le BDS apporte son soutien financier à des clubs et associations étudiantes dans leurs projets.",
		"clubs.vide":                  "Aucun club soutenu pour le moment.",
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
		"pied.sports":                 "Nos sports",
		"pied.competitions":           "Compétitions",
		"pied.clubs":                  "Clubs",
		"pied.siteecole":              "Site web de l'école",
		"pied.avenir":                 "à venir",
		"sport.horaires":              "Horaires",
		"sport.lieu":                  "Lieu",
		"sport.respo":                 "Respo",
		"sport.contact":               "Contact",
		"sport.competitions":          "Compétitions",
		"inscription.titre":           "Inscription",
		"inscription.identifiant":     "Identifiant",
		"inscription.motdepasse":      "Mot de passe",
		"inscription.confirmation":    "Confirme le mot de passe",
		"inscription.bouton":          "Créer mon compte",
		"inscription.erreur.existe":   "Cet identifiant est déjà pris.",
		"inscription.erreur.mdp":      "Les mots de passe ne correspondent pas.",
		"inscription.deja":            "Déjà inscrit ? Connecte-toi.",
		"accueil.actus.titre":         "Annonces",
		"accueil.actus.intro":         "Les annonces du BDS : compétitions, évènements et activités.",
		"langue.nom.fr":               "Français",
		"langue.nom.en":               "English",
	},
	"en": {
		"nav.accueil":                 "The association",
		"nav.sports":                  "Our sports",
		"nav.competitions":            "Competitions",
		"nav.clubs":                   "Clubs",
		"nav.actus":                   "News",
		"nav.equipe":                  "Team",
		"nav.connexion":               "Log in",
		"nav.deconnexion":             "Log out",
		"nav.inscription":             "Sign up",
		"nav.langue":                  "Language",
		"accueil.titre":               "The association that brings sport to life at ENSIIE",
		"accueil.association.titre":   "The association",
		"accueil.association.p1":      "The ENSIIE BDS is a dynamic team of %d members committed to providing all students with regular, varied, and high-quality access to sport within the school.",
		"accueil.association.p2":      "Sport plays a central role in student life at ENSIIE. Whether you are a seasoned athlete or just want to stay active, the BDS is here to support you all year long.",
		"accueil.association.photo":   "Photo coming soon",
		"accueil.presentation.bouton": "Meet the team",
		"accueil.stats.membres":       "BDS members",
		"accueil.stats.sports":        "Sports offered",
		"accueil.stats.seances":       "Weekly sessions",
		"accueil.stats.adherents":     "Members",
		"accueil.competitions.titre":  "Competitions",
		"accueil.competitions.p1":     "Every year our athletes take part in inter-school tournaments all over France, proudly representing the colours of ENSIIE and its BDS.",
		"accueil.competitions.p2":     "Competitors and supporters alike, these events are not to be missed.",
		"accueil.competitions.bouton": "More info",
		"accueil.evenement.site":      "%s website",
		"accueil.evenement.photo":     "Photo coming soon",
		"accueil.clubs.intro":         "The BDS financially supports several specialised student associations.",
		"accueil.clubs.bouton":        "See all clubs",
		"accueil.partenaires.titre":   "Our partners",
		"accueil.partenaires.intro":   "They support the BDS all year long. Discover their offers for ENSIIE students.",
		"accueil.actus.titre":         "Announcements",
		"accueil.actus.intro":         "BDS announcements: competitions, events and activities.",
		"partenaire.conditions":       "Offers subject to conditions.",
		"partenaire.site":             "Visit the website",
		"partenaire.flyer":            "View the flyer",
		"sports.titre":                "Our sports",
		"sports.intro":                "Discover the disciplines offered by the ENSIIE BDS.",
		"sports.compte":               "%d sports",
		"sports.unite":                "sports",
		"sports.filtre.tous":          "All sports",
		"sports.decouvrir":            "Discover",
		"sports.retour":               "Back to sports",
		"sports.vide":                 "No sport yet.",
		"actus.titre":                 "News",
		"actus.intro":                 "Everything about recent BDS ENSIIE events and activities.",
		"actus.vide":                  "No article found. Change your search.",
		"actus.recherche":             "Search a news item…",
		"compet.titre":                "Competitions",
		"compet.intro":                "Everything about competition at ENSIIE, from tournaments to FFSU championships.",
		"compet.vide":                 "No competition yet.",
		"clubs.titre":                 "Clubs and their funding",
		"clubs.intro":                 "The BDS supports student clubs and associations in their projects.",
		"clubs.vide":                  "No supported club yet.",
		"equipe.titre":                "BDS Team",
		"connexion.titre":             "Log in",
		"connexion.identifiant":       "Username",
		"connexion.motdepasse":        "Password",
		"connexion.bouton":            "Log in",
		"connexion.oublie":            "Forgot password?",
		"connexion.erreur":            "Wrong username or password.",
		"inscription.titre":           "Sign up",
		"inscription.identifiant":     "Username",
		"inscription.motdepasse":      "Password",
		"inscription.confirmation":    "Confirm the password",
		"inscription.bouton":          "Create my account",
		"inscription.erreur.existe":   "This username is already taken.",
		"inscription.erreur.mdp":      "The passwords do not match.",
		"inscription.deja":            "Already registered? Log in.",
		"pied.description":            "The Sports Bureau of ENSIIE organises and promotes sports activities for all the students of the school.",
		"pied.liensrapides":           "Quick links",
		"pied.informations":           "Information",
		"pied.contact":                "Contact",
		"pied.faq":                    "FAQ",
		"pied.apropos":                "About",
		"pied.actus":                  "News",
		"pied.equipe":                 "Our team",
		"pied.sports":                 "Our sports",
		"pied.competitions":           "Competitions",
		"pied.clubs":                  "Clubs",
		"pied.siteecole":              "School website",
		"pied.avenir":                 "coming soon",
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
	if c, err := r.Cookie("bds_ensiie_langue"); err == nil && c.Value != "" {
		return c.Value
	}
	return "fr"
}
