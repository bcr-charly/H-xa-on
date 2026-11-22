package kernel

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestParseConfigSections(t *testing.T) {
	cfg, averts := ParseConfig(`
serveur
  port = 9090
  nom  = BDS ENSIIE

palette
  fond = #202020

pages
  accueil = ON
  actus = ON
  inconnue = ON

securite
  secret = s3cret
`)
	if cfg.Port != 9090 {
		t.Errorf("port : %d", cfg.Port)
	}
	if cfg.Nom != "BDS ENSIIE" {
		t.Errorf("nom : %q", cfg.Nom)
	}
	if cfg.Palette.Fond != "#202020" {
		t.Errorf("palette : %+v", cfg.Palette)
	}
	if cfg.Palette.Accent != "#E8A400" {
		t.Errorf("palette par défaut non conservée : %+v", cfg.Palette)
	}
	if !cfg.Pages["accueil"] || cfg.Pages["planning"] {
		t.Errorf("pages : %+v", cfg.Pages)
	}
	if cfg.Secret != "s3cret" {
		t.Errorf("secret non lu")
	}
	if len(averts) != 0 {
		t.Fatalf("avertissements : %v (une page inconnue activée est avertie à l'enregistrement, pas au parse)", averts)
	}
}

func TestParseConfigValeursInvalidesNonFatales(t *testing.T) {
	cfg, averts := ParseConfig(`
serveur
  port = pas-un-port

palette
  fond = rouge-bizarre

pages
  accueil = PEUT-ETRE
`)
	if cfg.Port != 8080 {
		t.Errorf("port par défaut : %d", cfg.Port)
	}
	if cfg.Palette.Fond != "rouge-bizarre" {
		t.Errorf("couleur appliquée quand même attendu warning seulement")
	}
	if len(averts) != 3 {
		t.Errorf("avertissements : %+v", averts)
	}
}

func TestParseContenus(t *testing.T) {
	blocs, err := ParseContenus(`
[actu]
titre = Coupe desENSIÉ
date = 2026-01-12
texte = Victoire en finale.
  Score : 2-1.
image = assets/match.jpg

[actu]
titre = Semaine du sport
texte = Venez nombreux
`)
	if err != nil {
		t.Fatal(err)
	}
	actus := FiltreSections(blocs, "actu")
	if len(actus) != 2 {
		t.Fatalf("blocs : %d", len(actus))
	}
	if actus[0].Texte("texte") != "Victoire en finale.\nScore : 2-1." {
		t.Errorf("multi-ligne : %q", actus[0].Texte("texte"))
	}
	if actus[0].Texte("image") != "assets/match.jpg" {
		t.Errorf("image : %q", actus[0].Texte("image"))
	}
}

func TestContenuHorsBlocErreur(t *testing.T) {
	if _, err := ParseContenus("titre = sans bloc"); err == nil {
		t.Error("attendu une erreur pour du contenu hors bloc")
	}
}

func TestChargerContenusAbsent(t *testing.T) {
	blocs, err := ChargerContenus("/chemin/qui/nexiste/pas.txt")
	if err != nil || blocs != nil {
		t.Errorf("fichier absent : %v %v", blocs, err)
	}
}

func TestMotDePasse(t *testing.T) {
	h, err := HacherMotDePasse("bds2026!")
	if err != nil {
		t.Fatal(err)
	}
	if !VerifierMotDePasse("bds2026!", h) {
		t.Error("mot de passe correct refusé")
	}
	if VerifierMotDePasse("autre", h) {
		t.Error("mot de passe incorrect accepté")
	}
	if VerifierMotDePasse("x", "pas-un-hachage") {
		t.Error("hachage corrompu accepté")
	}
	h2, _ := HacherMotDePasse("bds2026!")
	if h == h2 {
		t.Error("deux hachages du même mot de passe identiques : sel non utilisé")
	}
}

func TestStockageFichiers(t *testing.T) {
	d := t.TempDir()
	s, err := NouveauStockageFichiers(d)
	if err != nil {
		t.Fatal(err)
	}
	if err := s.Ecrire("comptes", Enregistrement{Data: map[string]interface{}{"login": "ada", "role": "admin"}}); err != nil {
		t.Fatal(err)
	}
	e, err := s.ChercherPar("comptes", "login", "ada")
	if err != nil || e == nil {
		t.Fatalf("cherche : %v %v", e, err)
	}
	id := e.ID
	e.Data["role"] = "adherent"
	if err := s.Ecrire("comptes", *e); err != nil {
		t.Fatal(err)
	}
	items, _ := s.Charger("comptes")
	if len(items) != 1 || items[0].Data["role"] != "adherent" {
		t.Fatalf("réécriture : %+v", items)
	}
	if err := s.Ecrire("comptes", Enregistrement{Data: map[string]interface{}{"login": "bob"}}); err != nil {
		t.Fatal(err)
	}

	csvData, err := s.ExporterCSV("comptes")
	if err != nil {
		t.Fatal(err)
	}
	enTete := "id,login,role\n"
	if sortie := string(csvData); len(sortie) < len(enTete) || sortie[:len(enTete)] != enTete {
		t.Errorf("CSV en-tete : %q", sortie)
		t.Errorf("CSV en-tete : %q", string(csvData))
	}

	if err := s.PurgeAnnuelle("comptes"); err != nil {
		t.Fatal(err)
	}
	items, _ = s.Charger("comptes")
	if len(items) != 0 {
		t.Errorf("purge : %+v", items)
	}
	_ = id
}

func TestSessions(t *testing.T) {
	d := t.TempDir()
	g, err := NouveauGestionSessions(filepath.Join(d, "sessions"), "secret-partage")
	if err != nil {
		t.Fatal(err)
	}
	valeur, err := g.Ouvrir("admin")
	if err != nil {
		t.Fatal(err)
	}
	s, err := g.Valider(valeur)
	if err != nil || s == nil || s.Compte != "admin" {
		t.Fatalf("valider : %v %v", s, err)
	}
	if v, _ := g.Valider(valeur + "0"); v != nil {
		t.Error("signature falsifiée acceptée")
	}
	if v, _ := g.Valider("nimportequoi"); v != nil {
		t.Error("cookie mal formé accepté")
	}
	g.Fermer(valeur)
	if v, _ := g.Valider(valeur); v != nil {
		t.Error("session fermée encore valide")
	}

	valeurExpiree, _ := g.Ouvrir("admin")
	id := valeurExpiree[:36]
	sF := Session{ID: id, Compte: "admin", Expire: time.Now().Add(-time.Minute)}
	fichier := filepath.Join(d, "sessions", id+".json")
	donnees, _ := json.Marshal(sF)
	os.WriteFile(fichier, donnees, 0o600)
	if v, _ := g.Valider(valeurExpiree); v != nil {
		t.Error("session expirée acceptée")
	}
	if _, err := os.Stat(fichier); !os.IsNotExist(err) {
		t.Error("session expirée non purgée")
	}
}

func TestTraductions(t *testing.T) {
	if Traduit("fr", "nav.connexion") != "Connexion" {
		t.Error("traduction fr")
	}
	if Traduit("en", "nav.connexion") != "Log in" {
		t.Error("traduction en")
	}
	if Traduit("xx", "nav.actus") != "Actus" {
		t.Error("locale inconnue doit retomber sur fr")
	}
	if Traduit("fr", "cle.inexistante") != "cle.inexistante" {
		t.Error("clé inconnue doit retourner la clé")
	}
}

func TestEnregistrerPages(t *testing.T) {
	n, err := NouveauNoyau(Config{Pages: map[string]bool{"accueil": true, "boutique": true}}, filepath.Join(t.TempDir(), "data"))
	if err != nil {
		t.Fatal(err)
	}
	page := &pageTest{nom: "accueil"}
	averts := n.EnregistrerPages([]Page{page})
	if len(averts) != 1 || averts[0].Cle != "boutique" {
		t.Errorf("page activée non embarquée : %+v", averts)
	}
}

type pageTest struct{ nom string }

func (p *pageTest) Nom() string                                   { return p.nom }
func (p *pageTest) Chemins() []string                             { return []string{"/" + p.nom} }
func (p *pageTest) Servir(w http.ResponseWriter, r *http.Request) {}
