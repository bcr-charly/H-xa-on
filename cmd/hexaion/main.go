package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"

	"github.com/bcr-charly/H-xa-on/internal/kernel"
	"github.com/bcr-charly/H-xa-on/internal/pages"
)

func main() {
	dossier := flag.String("dossier", ".", "dossier de déploiement (contient conf.txt, contenu/, assets/, data/)")
	creerAdmin := flag.String("creer-admin", "", "créer un compte administrateur (login:motdepasse) puis quitter")
	flag.Parse()

	cfg, averts, err := kernel.ChargerConfig(filepath.Join(*dossier, "conf.txt"))
	if err != nil {
		log.Fatalf("lecture de conf.txt : %v", err)
	}
	for _, a := range averts {
		log.Printf("avertissement : %s", a)
	}

	dossierData := filepath.Join(*dossier, "data")
	noyau, err := kernel.NouveauNoyau(cfg, dossierData)
	if err != nil {
		log.Fatalf("initialisation du noyau : %v", err)
	}
	noyau.ServirAssets(filepath.Join(*dossier, "assets"))

	if *creerAdmin != "" {
		login, mdp, ok := couper(*creerAdmin)
		if !ok {
			log.Fatal("-creer-admin attend \"login:motdepasse\"")
		}
		if err := creerCompteAdmin(noyau, login, mdp); err != nil {
			log.Fatalf("création du compte : %v", err)
		}
		fmt.Printf("compte administrateur %q créé dans data/comptes.json\n", login)
		return
	}

	embarquees := pages.Toutes(noyau, filepath.Join(*dossier, "contenu"))
	for _, a := range noyau.EnregistrerPages(embarquees) {
		log.Printf("avertissement : %s", a)
	}

	log.Printf("%s démarre — pages actives : %s", cfg.Nom, nomsPages(embarquees, cfg))
	if err := noyau.Ecoute(); err != nil {
		log.Fatalf("serveur : %v", err)
	}
}

func couper(s string) (string, string, bool) {
	for i := 0; i < len(s); i++ {
		if s[i] == ':' {
			login, mdp := s[:i], s[i+1:]
			if login != "" && mdp != "" {
				return login, mdp, true
			}
			return "", "", false
		}
	}
	return "", "", false
}

func creerCompteAdmin(n *kernel.Noyau, login, mdp string) error {
	if e, err := n.Stockage.ChercherPar("comptes", "login", login); err != nil {
		return err
	} else if e != nil {
		return fmt.Errorf("le login %q existe déjà", login)
	}
	h, err := kernel.HacherMotDePasse(mdp)
	if err != nil {
		return err
	}
	return n.Stockage.Ecrire("comptes", kernel.Enregistrement{
		Data: map[string]interface{}{"login": login, "motdepasse": h, "role": "admin"},
	})
}

func nomsPages(ps []kernel.Page, cfg kernel.Config) string {
	res := ""
	for _, p := range ps {
		if active, ok := cfg.Pages[p.Nom()]; ok && !active {
			continue
		}
		if res != "" {
			res += ", "
		}
		res += p.Nom()
	}
	return res
}
