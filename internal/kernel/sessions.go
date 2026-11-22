package kernel

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	NomCookieSession = "hexaion_session"
	DureeSession     = 12 * time.Hour
)

// Session est l'état serveur d'une session ouverte.
type Session struct {
	ID     string    `json:"id"`
	Compte string    `json:"compte"`
	Expire time.Time `json:"expire"`
}

// GestionSessions gère les sessions : identifiants opaques en cookie signé,
// état complet côté serveur dans data/sessions, expiration automatique.
type GestionSessions struct {
	dossier string
	secret  []byte
	mu      sync.Mutex
}

func NouveauGestionSessions(dossier, secret string) (*GestionSessions, error) {
	if err := os.MkdirAll(dossier, 0o700); err != nil {
		return nil, err
	}
	return &GestionSessions{dossier: dossier, secret: []byte(secret)}, nil
}

func (g *GestionSessions) chemin(id string) string {
	return filepath.Join(g.dossier, id+".json")
}

func (g *GestionSessions) signer(id string) string {
	mac := hmac.New(sha256.New, g.secret)
	mac.Write([]byte(id))
	return hex.EncodeToString(mac.Sum(nil))
}

// Ouvrir crée une session pour un compte, l'persiste et retourne la valeur
// opaque du cookie : "id.signature".
func (g *GestionSessions) Ouvrir(compte string) (string, error) {
	b := make([]byte, 18)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	id := hex.EncodeToString(b)
	s := Session{ID: id, Compte: compte, Expire: time.Now().Add(DureeSession)}
	data, err := json.Marshal(s)
	if err != nil {
		return "", err
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	if err := ecritureAtomique(g.chemin(id), data); err != nil {
		return "", err
	}
	return id + "." + g.signer(id), nil
}

// Valider contrôle le cookie et retourne la session si elle est valide.
func (g *GestionSessions) Valider(valeur string) (*Session, error) {
	parts := strings.SplitN(valeur, ".", 2)
	if len(parts) != 2 {
		return nil, nil
	}
	id, sig := parts[0], parts[1]
	if !hmac.Equal([]byte(g.signer(id)), []byte(sig)) {
		return nil, nil
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	data, err := os.ReadFile(g.chemin(id))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var s Session
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, nil
	}
	if time.Now().After(s.Expire) {
		os.Remove(g.chemin(id))
		return nil, nil
	}
	return &s, nil
}

// Fermer détruit la session côté serveur.
func (g *GestionSessions) Fermer(valeur string) {
	parts := strings.SplitN(valeur, ".", 2)
	if len(parts) != 2 {
		return
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	os.Remove(g.chemin(parts[0]))
}

// CookieSession écrit le cookie de session sur la réponse.
func CookieSession(w http.ResponseWriter, valeur string) {
	http.SetCookie(w, &http.Cookie{
		Name:     NomCookieSession,
		Value:    valeur,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(DureeSession.Seconds()),
	})
}

// EffaceCookieSession retire le cookie de session.
func EffaceCookieSession(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     NomCookieSession,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})
}

func (s Session) String() string {
	return fmt.Sprintf("session %s (compte %s)", s.ID, s.Compte)
}
