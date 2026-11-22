package pages

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bcr-charly/H-xa-on/internal/kernel"
)

func noyauTest(t *testing.T, racine string) (*kernel.Noyau, *httptest.Server) {
	t.Helper()
	cfg := kernel.Config{
		Port:    8080,
		Nom:     "BDS ENSIIE",
		Palette: kernel.PaletteDefaut(),
		Pages:   map[string]bool{"accueil": true, "actus": true, "equipe": true, "connexion": true},
		Secret:  "secret-test",
	}
	n, err := kernel.NouveauNoyau(cfg, filepath.Join(racine, "data"))
	if err != nil {
		t.Fatal(err)
	}
	n.ServirAssets(filepath.Join(racine, "assets"))
	for _, a := range n.EnregistrerPages(Toutes(n, filepath.Join(racine, "contenu"))) {
		t.Logf("avertissement : %s", a)
	}
	ts := httptest.NewServer(n)
	t.Cleanup(ts.Close)
	return n, ts
}

func get(t *testing.T, url string) (int, string) {
	t.Helper()
	resp, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, string(b)
}

func instanceTest(t *testing.T) string {
	racine := t.TempDir()
	for d := range map[string]bool{"contenu": true, "assets": true, "data": true} {
		os.MkdirAll(filepath.Join(racine, d), 0o755)
	}
	os.WriteFile(filepath.Join(racine, "contenu", "actus.txt"), []byte(`
[actu]
titre = Championnat d'ultimate
date = 2026-03-08
texte = Victoire 13-9 en finale départementale.
image = assets/match.jpg
`), 0o644)
	os.WriteFile(filepath.Join(racine, "contenu", "equipe.txt"), []byte(`
[membre]
nom = Camille Dupont
role = Présidente
photo = camille.jpg
`), 0o644)
	return racine
}

func TestAccueilRendu(t *testing.T) {
	racine := instanceTest(t)
	_, ts := noyauTest(t, racine)
	code, corps := get(t, ts.URL+"/")
	if code != 200 {
		t.Fatalf("code : %d", code)
	}
	if !strings.Contains(corps, "les chasseurs et les cueilleurs") {
		t.Error("slogan absent de l'accueil")
	}
	if !strings.Contains(corps, "/connexion") {
		t.Error("bouton connexion absent")
	}
	if !strings.Contains(corps, "/assets/logo.svg") {
		t.Error("logo permanent absent")
	}
	if !strings.Contains(corps, `class="logo-flottant"`) {
		t.Error("logo permanent (bouton menu home) absent")
	}
	if !strings.Contains(corps, `body class="page-accueil"`) {
		t.Error("la page d'accueil doit porter la classe page-accueil")
	}
	if !strings.Contains(corps, `class="langue-bloc"`) || !strings.Contains(corps, `class="langue-menu"`) {
		t.Error("onglet de langue attendu sous forme de menu ouvrable")
	}
	if !strings.Contains(corps, "bds_ensiie") || !strings.Contains(corps, "ensiie.fr") {
		t.Error("liens du pied de page (Instagram, école) absents")
	}
	if !strings.Contains(corps, "Évry-Courcouronnes") {
		t.Error("coordonnées du pied de page (mentions légales) absentes")
	}
	main := corps[strings.Index(corps, "<main>"):strings.Index(corps, "</main>")]
	if strings.Contains(main, "http://") || strings.Contains(main, "https://") {
		t.Error("ressource distante dans le contenu de la page (le pied de page doit rester la seule exception)")
	}
}

func TestPageInterneHeaderBlanc(t *testing.T) {
	racine := instanceTest(t)
	_, ts := noyauTest(t, racine)
	_, corps := get(t, ts.URL+"/equipe")
	if !strings.Contains(corps, `body class="page-interne"`) {
		t.Error("une page interne doit porter la classe page-interne (header blanc)")
	}
	if !strings.Contains(corps, "rangee-equipe") {
		t.Error("rangées symétriques attendues sur la page équipe")
	}
}

func TestActusRendu(t *testing.T) {
	racine := instanceTest(t)
	_, ts := noyauTest(t, racine)
	code, corps := get(t, ts.URL+"/actus")
	if code != 200 {
		t.Fatalf("code : %d", code)
	}
	if !strings.Contains(corps, "Championnat d&#39;ultimate") {
		t.Error("carte d'actualité absente")
	}
	if !strings.Contains(corps, `class="carte"`) {
		t.Error("composant carte non utilisé")
	}
}

func TestEquipeRendu(t *testing.T) {
	racine := instanceTest(t)
	_, ts := noyauTest(t, racine)
	code, corps := get(t, ts.URL+"/equipe")
	if code != 200 {
		t.Fatalf("code : %d", code)
	}
	if !strings.Contains(corps, "Camille Dupont") {
		t.Error("membre absent de la grille")
	}
	if !strings.Contains(corps, "Présidente") {
		t.Error("role absent de la grille")
	}
}

func TestConnexionRenduEtRefus(t *testing.T) {
	racine := instanceTest(t)
	_, ts := noyauTest(t, racine)
	code, corps := get(t, ts.URL+"/connexion")
	if code != 200 || !strings.Contains(corps, `name="motdepasse"`) {
		t.Fatalf("formulaire absent : %d", code)
	}
	resp, err := http.PostForm(ts.URL+"/connexion", map[string][]string{
		"login": {"personne"}, "motdepasse": {"mauvais"},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(b), "incorrect") && !strings.Contains(string(b), "Identifiant") {
		t.Error("échec de connexion : message d'erreur attendu")
	}
}

func TestPageNonEmbarqueeAvertie(t *testing.T) {
	racine := instanceTest(t)
	cfg := kernel.Config{
		Pages:  map[string]bool{"accueil": true, "boutique": true},
		Secret: "s",
	}
	n, err := kernel.NouveauNoyau(cfg, filepath.Join(racine, "data"))
	if err != nil {
		t.Fatal(err)
	}
	averts := n.EnregistrerPages(Toutes(n, filepath.Join(racine, "contenu")))
	trouve := false
	for _, a := range averts {
		if a.Cle == "boutique" {
			trouve = true
		}
	}
	if !trouve {
		t.Errorf("boutique activée non embarquée doit produire un avertissement : %v", averts)
	}
}
