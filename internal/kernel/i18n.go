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
		"adhesion.titre":              "L'adhésion",
		"adhesion.pourquoi.titre":     "Pourquoi devenir membre ?",
		"adhesion.comprend.titre":     "Ton adhésion comprend",
		"adhesion.comment.titre":      "Comment adhérer",
		"adhesion.bouton":             "Adhérer en ligne",
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
		"pied.adhesion":               "Adhésion",
		"pied.siteecole":              "Site web de l'école",
		"pied.avenir":                 "à venir",
		"sport.horaires":              "Horaires",
		"sport.lieu":                  "Lieu",
		"sport.respo":                 "Respo",
		"sport.contact":               "Contact",
		"sport.competitions":          "Compétitions",
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

// LocaleDepuisRequete ramène toujours au français : le site est
// monolingue, la mécanique reste pour une future variante linguistique.
func LocaleDepuisRequete(r *http.Request) string {
	return "fr"
}
