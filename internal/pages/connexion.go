package pages

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/bcr-charly/H-xa-on/internal/components"
	"github.com/bcr-charly/H-xa-on/internal/kernel"
)

// PageConnexion est l'écran de connexion servi par le même binaire.
// La création des comptes reste réservée aux administrateurs.
type PageConnexion struct {
	Noyau *kernel.Noyau
}

func (p *PageConnexion) Nom() string { return "connexion" }
func (p *PageConnexion) Chemins() []string {
	return []string{"/connexion", "/deconnexion", "/admin/export"}
}

func (p *PageConnexion) Servir(w http.ResponseWriter, r *http.Request) {
	locale := kernel.LocaleDepuisRequete(r)
	switch r.URL.Path {
	case "/connexion":
		p.servirConnexion(w, r, locale)
	case "/deconnexion":
		p.servirDeconnexion(w, r, locale)
	case "/admin/export":
		p.servirExport(w, r, locale)
	default:
		http.NotFound(w, r)
	}
}

func (p *PageConnexion) servirDeconnexion(w http.ResponseWriter, r *http.Request, locale string) {
	if c, err := r.Cookie(kernel.NomCookieSession); err == nil {
		p.Noyau.Sessions.Fermer(c.Value)
	}
	kernel.EffaceCookieSession(w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// servirExport livre le CSV d'une collection à un administrateur connecté,
// composant du noyau dès la phase stockage.
func (p *PageConnexion) servirExport(w http.ResponseWriter, r *http.Request, locale string) {
	if !p.Noyau.EstAdmin(r) {
		http.Redirect(w, r, "/connexion", http.StatusSeeOther)
		return
	}
	collection := r.URL.Query().Get("collection")
	if collection == "" {
		http.Error(w, "collection manquante", http.StatusBadRequest)
		return
	}
	data, err := p.Noyau.Stockage.ExporterCSV(collection)
	if err != nil {
		http.Error(w, "export impossible", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q.csv", collection))
	w.Write(data)
}

func (p *PageConnexion) servirConnexion(w http.ResponseWriter, r *http.Request, locale string) {
	if p.Noyau.CompteCourant(r) != "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	erreur := ""
	if r.Method == http.MethodPost {
		login := r.FormValue("login")
		mdp := r.FormValue("motdepasse")
		e, err := p.Noyau.Stockage.ChercherPar("comptes", "login", login)
		if err == nil && e != nil {
			if stocke, _ := e.Data["motdepasse"].(string); kernel.VerifierMotDePasse(mdp, stocke) {
				valeur, err := p.Noyau.Sessions.Ouvrir(login)
				if err == nil {
					kernel.CookieSession(w, valeur)
					http.Redirect(w, r, "/", http.StatusSeeOther)
					return
				}
			}
		}
		erreur = kernel.Traduit(locale, "connexion.erreur")
	}

	var ligneErreur string
	if erreur != "" {
		ligneErreur = fmt.Sprintf(`<p class="erreur">%s</p>`, template.HTMLEscapeString(erreur))
	}
	contenu := template.HTML(string(components.Defaut.BandeTitre(kernel.Traduit(locale, "connexion.titre"))) +
		fmt.Sprintf(`<div class="conteneur"><form class="connexion-form" method="post" action="/connexion">
<label>%s<input type="text" name="login" autocomplete="username" required></label>
<label>%s<input type="password" name="motdepasse" autocomplete="current-password" required></label>
%s
<button type="submit" class="bouton">%s</button>
<a class="lien" href="/connexion">%s</a>
</form></div>`,
			template.HTMLEscapeString(kernel.Traduit(locale, "connexion.identifiant")),
			template.HTMLEscapeString(kernel.Traduit(locale, "connexion.motdepasse")),
			ligneErreur,
			template.HTMLEscapeString(kernel.Traduit(locale, "connexion.bouton")),
			template.HTMLEscapeString(kernel.Traduit(locale, "connexion.oublie")),
		))
	p.Noyau.Rendu(w, http.StatusOK, page(p.Noyau, locale, p.Nom(), kernel.Traduit(locale, "connexion.titre"), contenu, false), nil)
}
