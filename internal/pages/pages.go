package pages

import (
	"github.com/bcr-charly/H-xa-on/internal/kernel"
)

// Toutes construit la liste des pages embarquées. Le noyau ne connaît
// aucune page en particulier : c'est ici, et seulement ici, que le
// programme enregistre les pages disponibles.
func Toutes(n *kernel.Noyau, dossierContenu string) []kernel.Page {
	return []kernel.Page{
		&PageAccueil{Noyau: n, DossierContenu: dossierContenu},
		&PageSports{Noyau: n, DossierContenu: dossierContenu},
		&PageSport{Noyau: n, DossierContenu: dossierContenu},
		&PageCompetitions{Noyau: n, DossierContenu: dossierContenu},
		&PageClubs{Noyau: n, DossierContenu: dossierContenu},
		&PageActus{Noyau: n, DossierContenu: dossierContenu},
		&PageAdhesion{Noyau: n, DossierContenu: dossierContenu},
		&PageEquipe{Noyau: n, DossierContenu: dossierContenu},
		&PageConnexion{Noyau: n},
	}
}
