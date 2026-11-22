package kernel

import (
	"crypto/rand"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

// Enregistrement est une ligne d'une collection : identifiant unique
// et données arbitraires sérialisables.
type Enregistrement struct {
	ID   string                 `json:"id"`
	Data map[string]interface{} `json:"data"`
}

// Stockage est l'interface unique d'accès aux données persistantes.
// L'implémentation de référence est un fichier plat par collection ;
// elle peut être remplacée (SQLite par exemple) sans toucher au reste.
type Stockage interface {
	Charger(collection string) ([]Enregistrement, error)
	Chercher(collection, id string) (*Enregistrement, error)
	ChercherPar(collection, cle, valeur string) (*Enregistrement, error)
	Ecrire(collection string, e Enregistrement) error
	Supprimer(collection, id string) error
	ExporterCSV(collection string) ([]byte, error)
}

// StockageFichiers est l'implémentation de référence : un fichier JSON par
// collection dans le dossier data/, écritures sérialisées par verrou mémoire
// puis réécriture atomique (fichier temporaire + remplacement).
type StockageFichiers struct {
	dossier string
	mu      sync.Mutex
}

func NouveauStockageFichiers(dossier string) (*StockageFichiers, error) {
	if err := os.MkdirAll(dossier, 0o700); err != nil {
		return nil, err
	}
	return &StockageFichiers{dossier: dossier}, nil
}

func (s *StockageFichiers) chemin(collection string) string {
	return filepath.Join(s.dossier, collection+".json")
}

func (s *StockageFichiers) Charger(collection string) ([]Enregistrement, error) {
	return s.charger(collection)
}

// charger lit la collection sans prendre le verrou : les méthodes qui
// détiennent déjà le verrou (Ecrire, Supprimer) l'utilisent.
func (s *StockageFichiers) charger(collection string) ([]Enregistrement, error) {
	data, err := os.ReadFile(s.chemin(collection))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var res []Enregistrement
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, fmt.Errorf("collection %s corrompue : %w", collection, err)
	}
	return res, nil
}

func (s *StockageFichiers) Chercher(collection, id string) (*Enregistrement, error) {
	items, err := s.Charger(collection)
	if err != nil {
		return nil, err
	}
	for i := range items {
		if items[i].ID == id {
			return &items[i], nil
		}
	}
	return nil, nil
}

func (s *StockageFichiers) ChercherPar(collection, cle, valeur string) (*Enregistrement, error) {
	items, err := s.Charger(collection)
	if err != nil {
		return nil, err
	}
	for i := range items {
		if v, ok := items[i].Data[cle].(string); ok && v == valeur {
			return &items[i], nil
		}
	}
	return nil, nil
}

func (s *StockageFichiers) Ecrire(collection string, e Enregistrement) error {
	if e.ID == "" {
		id, err := nouvelID()
		if err != nil {
			return err
		}
		e.ID = id
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	items, err := s.charger(collection)
	if err != nil {
		return err
	}
	remplace := false
	for i := range items {
		if items[i].ID == e.ID {
			items[i] = e
			remplace = true
			break
		}
	}
	if !remplace {
		items = append(items, e)
	}
	data, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return err
	}
	return ecritureAtomique(s.chemin(collection), data)
}

func (s *StockageFichiers) Supprimer(collection, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	items, err := s.charger(collection)
	if err != nil {
		return err
	}
	var res []Enregistrement
	for _, it := range items {
		if it.ID != id {
			res = append(res, it)
		}
	}
	data, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		return err
	}
	return ecritureAtomique(s.chemin(collection), data)
}

// ExporterCSV produit l'export de la collection : une colonne par clé
// rencontrée, colonnes triées pour un résultat stable, en-têtes inclus.
func (s *StockageFichiers) ExporterCSV(collection string) ([]byte, error) {
	items, err := s.Charger(collection)
	if err != nil {
		return nil, err
	}
	cles := map[string]bool{}
	for _, it := range items {
		for k := range it.Data {
			cles[k] = true
		}
	}
	colonnes := []string{"id"}
	for k := range cles {
		colonnes = append(colonnes, k)
	}
	sort.Strings(colonnes[1:])

	var sortie strings.Builder
	cw := csv.NewWriter(&sortie)
	if err := cw.Write(colonnes); err != nil {
		return nil, err
	}
	for _, it := range items {
		ligne := make([]string, len(colonnes))
		ligne[0] = it.ID
		for i, k := range colonnes[1:] {
			switch v := it.Data[k].(type) {
			case string:
				ligne[i+1] = v
			case nil:
				ligne[i+1] = ""
			default:
				ligne[i+1] = fmt.Sprintf("%v", v)
			}
		}
		if err := cw.Write(ligne); err != nil {
			return nil, err
		}
	}
	cw.Flush()
	return []byte(sortie.String()), nil
}

// ecritureAtomique écrit via un fichier temporaire puis renomme,
// pour qu'un crash ne laisse jamais de fichier à moitié écrit.
func ecritureAtomique(chemin string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(chemin), 0o700); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(chemin), "."+filepath.Base(chemin)+".tmp*")
	if err != nil {
		return err
	}
	nomTmp := tmp.Name()
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		os.Remove(nomTmp)
		return err
	}
	if err := tmp.Close(); err != nil {
		os.Remove(nomTmp)
		return err
	}
	if err := os.Rename(nomTmp, chemin); err != nil {
		os.Remove(nomTmp)
		return err
	}
	return nil
}

func nouvelID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", b), nil
}

// PurgeAnnuelle supprime toutes les données des collections données,
// dans le dossier data/ comme pour les sessions.
// C'est la fonction de fin d'année scolaire.
func (s *StockageFichiers) PurgeAnnuelle(collections ...string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, c := range collections {
		if err := os.Remove(s.chemin(c)); err != nil && !os.IsNotExist(err) {
			return err
		}
		if err := os.RemoveAll(filepath.Join(s.dossier, "sessions")); err != nil {
			return err
		}
	}
	return nil
}

// DateNow est l'heure courante, isolée pour les tests.
var DateNow = time.Now
