package kernel

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/bcr-charly/siteweb-bds-ensiie/internal/components"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
)

// Page décrit une page embarquée dans le binaire. Chaque module de page
// s'enregistre auprès du noyau ; le noyau ne connaît aucune page en particulier.
type Page interface {
	// Nom est la clé de configuration de la page (ex. "accueil").
	Nom() string
	// Chemins retourne les chemins d'URL servis par la page.
	Chemins() []string
	// Servir répond à une requête sur un des chemins de la page.
	Servir(w http.ResponseWriter, r *http.Request)
}

// Noyau est le serveur : HTTP, routage, configuration, rendu, sessions.
type Noyau struct {
	Config   Config
	Sessions *GestionSessions
	Stockage Stockage
	pages    map[string]Page
	mux      *http.ServeMux
}

func NouveauNoyau(cfg Config, dossierData string) (*Noyau, error) {
	secret := cfg.Secret
	if secret == "" {
		// Sans secret configuré, un secret aléatoire est généré au démarrage :
		// les sessions ne survivent pas à un redémarrage, ce qui est acceptable
		// pour un premier lancement. Configurer [securite] secret le stabilise.
		b := make([]byte, 32)
		rand.Read(b)
		secret = hex.EncodeToString(b)
	}
	st, err := NouveauStockageFichiers(dossierData)
	if err != nil {
		return nil, err
	}
	sessions, err := NouveauGestionSessions(filepath.Join(dossierData, "sessions"), secret)
	if err != nil {
		return nil, err
	}
	n := &Noyau{
		Config:   cfg,
		Sessions: sessions,
		Stockage: st,
		pages:    map[string]Page{},
		mux:      http.NewServeMux(),
	}
	n.finaliser()
	return n, nil
}

// finaliser enregistre les chemins du noyau lui-même :
// la racine (page d'accueil si elle est embarquée, sinon 404).
func (n *Noyau) finaliser() {
	n.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		if p, ok := n.pages["accueil"]; ok {
			p.Servir(w, r)
			return
		}
		http.NotFound(w, r)
	})
}

// ServeHTTP fait du noyau un http.Handler standard.
func (n *Noyau) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	n.mux.ServeHTTP(w, r)
}

// EnregistrerPage ajoute une page à la table d'enregistrement du noyau.
// Une page absente de la configuration ou désactivée (OFF) est ignorée.
// Retourne false si la page n'est pas activée dans la configuration.
func (n *Noyau) EnregistrerPage(p Page) bool {
	active, connu := n.Config.Pages[p.Nom()]
	if connu && !active {
		return false
	}
	if !connu {
		// Page absente de conf.txt : activée par défaut, un avertissement
		// a déjà été émis par la validation de configuration côté déploiement.
		active = true
	}
	n.pages[p.Nom()] = p
	for _, c := range p.Chemins() {
		n.mux.Handle(c, http.HandlerFunc(p.Servir))
	}
	return true
}

// EnregistrerPages enregistre les pages et vérifie la cohérence avec la
// configuration : une page activée mais non embarquée produit un
// avertissement et est ignorée, jamais un plantage.
func (n *Noyau) EnregistrerPages(pages []Page) []Avertissement {
	var averts []Avertissement
	embarquees := map[string]bool{}
	for _, p := range pages {
		if n.EnregistrerPage(p) {
			embarquees[p.Nom()] = true
		}
	}
	for nom, active := range n.Config.Pages {
		if active && !embarquees[nom] {
			averts = append(averts, Avertissement{Section: "pages", Cle: nom, Message: "page activée mais non embarquée dans ce binaire, ignorée"})
		}
	}
	return averts
}

// ServirAssets expose le dossier assets/ : images, logo, photos des pages.
func (n *Noyau) ServirAssets(dossier string) {
	fs := http.FileServer(http.Dir(dossier))
	n.mux.Handle("/assets/", http.StripPrefix("/assets/", fs))
}

// SessionCourante retourne la session valide de la requête, ou nil.
func (n *Noyau) SessionCourante(r *http.Request) *Session {
	c, err := r.Cookie(NomCookieSession)
	if err != nil {
		return nil
	}
	s, err := n.Sessions.Valider(c.Value)
	if err != nil {
		return nil
	}
	return s
}

// CompteCourant retourne le login de l'adhérent connecté.
func (n *Noyau) CompteCourant(r *http.Request) string {
	if s := n.SessionCourante(r); s != nil {
		return s.Compte
	}
	return ""
}

// EstAdmin indique si la requête provient d'un administrateur connecté.
func (n *Noyau) EstAdmin(r *http.Request) bool {
	login := n.CompteCourant(r)
	if login == "" {
		return false
	}
	e, err := n.Stockage.ChercherPar("comptes", "login", login)
	if err != nil || e == nil {
		return false
	}
	return e.Data["role"] == "admin"
}

// Rendu exécute un template HTML avec les fonctions de composants et les
// données de la page.
func (n *Noyau) Rendu(w http.ResponseWriter, statut int, texteTemplate string, donnees interface{}) {
	t := template.New("page").Funcs(n.Fonctions())
	if _, err := t.Parse(texteTemplate); err != nil {
		log.Printf("rendu : %v", err)
		http.Error(w, "erreur interne", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(statut)
	if err := t.Execute(w, donnees); err != nil {
		log.Printf("rendu : %v", err)
	}
}

// Fonctions expose le vocabulaire de composants à tous les templates.
func (n *Noyau) Fonctions() template.FuncMap {
	return components.Defaut.FuncMap()
}

// Ecoute démarre le serveur HTTP sur le port configuré.
func (n *Noyau) Ecoute() error {
	adresse := ":" + strconv.Itoa(n.Config.Port)
	log.Printf("%s écoute sur %s", n.Config.Nom, adresse)
	return http.ListenAndServe(adresse, n.mux)
}
