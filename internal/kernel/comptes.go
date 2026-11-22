package kernel

import (
	"crypto/pbkdf2"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

// Paramètres du hachage de mots de passe : itérations élevées pour
// rendre toute attaque par force brute coûteuse (hachage lent avec sel).
const (
	pbkdf2Iterations = 210000
	pbkdf2TailleCle  = 32
	tailleSel        = 16
)

// HacherMotDePasse produit "pbkdf2$iterations$sellenhex$empreintehex".
func HacherMotDePasse(mdp string) (string, error) {
	sel := make([]byte, tailleSel)
	if _, err := rand.Read(sel); err != nil {
		return "", err
	}
	cle, err := pbkdf2.Key(sha256.New, mdp, sel, pbkdf2Iterations, pbkdf2TailleCle)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("pbkdf2$%d$%s$%s", pbkdf2Iterations, hex.EncodeToString(sel), hex.EncodeToString(cle)), nil
}

// VerifierMotDePasse contrôle un mot de passe face à son hachage stocké.
func VerifierMotDePasse(mdp, stocke string) bool {
	parts := strings.Split(stocke, "$")
	if len(parts) != 4 || parts[0] != "pbkdf2" {
		return false
	}
	iterations, err := strconv.Atoi(parts[1])
	if err != nil {
		return false
	}
	selHex, empreinteHex := parts[2], parts[3]
	if iterations <= 0 || iterations > 10_000_000 {
		return false
	}
	sel, err := hex.DecodeString(selHex)
	if err != nil {
		return false
	}
	cle, err := pbkdf2.Key(sha256.New, mdp, sel, iterations, pbkdf2TailleCle)
	if err != nil {
		return false
	}
	empreinte, err := hex.DecodeString(empreinteHex)
	if err != nil {
		return false
	}
	if len(cle) != len(empreinte) {
		return false
	}
	var diff byte
	for i := range cle {
		diff |= cle[i] ^ empreinte[i]
	}
	return diff == 0
}
