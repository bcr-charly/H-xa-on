package pages

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/bcr-charly/siteweb-bds-ensiie/internal/components"
	"github.com/bcr-charly/siteweb-bds-ensiie/internal/kernel"
)

// PageInscription est l'écran d'inscription : un élève crée son compte
// (identifiant + mot de passe + confirmation). La gestion complète des
// inscriptions (cotisation, rôles, administration) arrive avec la partie
// back-end ; ici le compte est créé avec le rôle « adherent » et la
// session est ouverte directement.
type PageInscription struct {
	Noyau *kernel.Noyau
}

func (p *PageInscription) Nom() string       { return "inscription" }
func (p *PageInscription) Chemins() []string { return []string{"/inscription"} }

func (p *PageInscription) Servir(w http.ResponseWriter, r *http.Request) {
	locale := kernel.LocaleDepuisRequete(r)
	if p.Noyau.CompteCourant(r) != "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	erreur := ""
	if r.Method == http.MethodPost {
		login := strings.TrimSpace(r.FormValue("login"))
		mdp := r.FormValue("motdepasse")
		confirmation := r.FormValue("confirmation")
		switch {
		case login == "" || mdp == "":
			// champs requis : le formulaire est required, cas de garde.
		case mdp != confirmation:
			erreur = kernel.Traduit(locale, "inscription.erreur.mdp")
		default:
			e, err := p.Noyau.Stockage.ChercherPar("comptes", "login", login)
			if err != nil {
				http.Error(w, "stockage illisible", http.StatusInternalServerError)
				return
			}
			if e != nil {
				erreur = kernel.Traduit(locale, "inscription.erreur.existe")
			} else {
				h, err := kernel.HacherMotDePasse(mdp)
				if err != nil {
					http.Error(w, "hachage impossible", http.StatusInternalServerError)
					return
				}
				if err := p.Noyau.Stockage.Ecrire("comptes", kernel.Enregistrement{
					Data: map[string]interface{}{"login": login, "motdepasse": h, "role": "adherent"},
				}); err != nil {
					http.Error(w, "création impossible", http.StatusInternalServerError)
					return
				}
				if valeur, err := p.Noyau.Sessions.Ouvrir(login); err == nil {
					kernel.CookieSession(w, valeur)
					http.Redirect(w, r, "/", http.StatusSeeOther)
					return
				}
			}
		}
	}

	var ligneErreur string
	if erreur != "" {
		ligneErreur = fmt.Sprintf(`<p class="erreur">%s</p>`,
			template.HTMLEscapeString(erreur))
	}
	contenu := template.HTML(string(components.Defaut.BandeTitre(kernel.Traduit(locale, "inscription.titre"))) +
		fmt.Sprintf(`<div class="conteneur"><form class="connexion-form" method="post" action="/inscription">
<label>%s<input type="text" name="login" autocomplete="username" required></label>
<label>%s<input type="password" name="motdepasse" autocomplete="new-password" required></label>
<label>%s<input type="password" name="confirmation" autocomplete="new-password" required></label>
%s
<button type="submit" class="bouton">%s</button>
<a class="lien" href="/connexion">%s</a>
</form></div>`,
			template.HTMLEscapeString(kernel.Traduit(locale, "inscription.identifiant")),
			template.HTMLEscapeString(kernel.Traduit(locale, "inscription.motdepasse")),
			template.HTMLEscapeString(kernel.Traduit(locale, "inscription.confirmation")),
			ligneErreur,
			template.HTMLEscapeString(kernel.Traduit(locale, "inscription.bouton")),
			template.HTMLEscapeString(kernel.Traduit(locale, "inscription.deja")),
		))
	p.Noyau.Rendu(w, http.StatusOK, page(p.Noyau, locale, p.Nom(), kernel.Traduit(locale, "inscription.titre"), contenu, false), nil)
}
