package main3

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

// Pendu représente l'état du jeu du Pendu
type Pendu struct {
	MotADeviner           string
	DevineActuelle        string
	DevinettesIncorrectes []string
	TentativesRestantes   int
	LettresDevinees       []string
	MotsDevines           []string
}

// NouveauPendu initialise un nouveau jeu du Pendu avec un mot à deviner
func NouveauPendu(cheminMots string) *Pendu {
	motAleatoire := choisirMotAleatoireDepuisFichier(cheminMots)

	return &Pendu{
		MotADeviner:           motAleatoire,
		DevineActuelle:        revelerLettres(motAleatoire),
		DevinettesIncorrectes: []string{},
		TentativesRestantes:   10,
		LettresDevinees:       []string{},
		MotsDevines:           []string{},
	}
}

// RendrePagePendu rend la page du jeu du Pendu
func RendrePagePendu(jeu *Pendu) {
	clearTerminal()

	fmt.Println("Mot à deviner :", jeu.DevineActuelle)
	fmt.Println("Devinettes incorrectes :", strings.Join(jeu.DevinettesIncorrectes, ", "))
	fmt.Println("Tentatives restantes :", jeu.TentativesRestantes)

	afficherPositionJose(jeu.TentativesRestantes)

	if jeu.TentativesRestantes == 0 {
		fmt.Println("Désolé, vous avez épuisé toutes vos tentatives. Le mot était :", jeu.MotADeviner)
	} else if !strings.Contains(jeu.DevineActuelle, "_") {
		fmt.Println("Félicitations ! Vous avez trouvé le mot :", jeu.MotADeviner)
	} else {
		fmt.Print("Entrez une lettre ou un mot : ")
		fmt.Println("Entrer STOP pour sauvegarder et arreter")
	}
}

// PenduHandler gère les essais pour le jeu du Pendu
func PenduHandler(jeu *Pendu, saisie string) {
	if jeu.DejaDevine(saisie) {
		fmt.Println("Vous avez déjà deviné cela. Veuillez essayer autre chose.")
		return
	}

	if len(saisie) == 1 {
		// Vérifier si la lettre devinée est dans le mot
		if strings.Contains(jeu.MotADeviner, saisie) {
			jeu.DevineActuelle = revelerLettresCorrectes(jeu.MotADeviner, jeu.DevineActuelle, saisie)

			// Vérifier si le mot a été entièrement découvert
			if !strings.Contains(jeu.DevineActuelle, "_") {
				return
			}
		} else {
			// Mauvaise devinette, ajouter à la liste des devinettes incorrectes
			jeu.TentativesRestantes--
			jeu.DevinettesIncorrectes = append(jeu.DevinettesIncorrectes, saisie)
		}

		jeu.LettresDevinees = append(jeu.LettresDevinees, saisie)
	} else if len(saisie) == len(jeu.MotADeviner) && saisie == jeu.MotADeviner {
		// Le joueur tente de deviner le mot complet
		jeu.DevineActuelle = jeu.MotADeviner
	} else {
		// Mauvaise tentative de deviner le mot complet, déduire 2 essais
		jeu.TentativesRestantes -= 2
	}

	jeu.MotsDevines = append(jeu.MotsDevines, saisie)
}

// DejaDevine vérifie si la lettre ou le mot a déjà été deviné
func (jeu *Pendu) DejaDevine(saisie string) bool {
	for _, lettre := range jeu.LettresDevinees {
		if lettre == saisie {
			return true
		}
	}

	for _, mot := range jeu.MotsDevines {
		if mot == saisie {
			return true
		}
	}

	return false
}

// choisirMotAleatoireDepuisFichier choisit un mot aléatoire depuis un fichier
func choisirMotAleatoireDepuisFichier(chemin string) string {
	fichier, err := os.Open(chemin)
	if err != nil {
		fmt.Println("Erreur lors de l'ouverture du fichier :", err)
		os.Exit(1)
	}
	defer fichier.Close()

	scanner := bufio.NewScanner(fichier)
	mots := []string{}
	for scanner.Scan() {
		mots = append(mots, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Erreur lors de la lecture du fichier :", err)
		os.Exit(1)
	}

	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(mots))
	return mots[index]
}

// afficherPositionJose affiche la position actuelle de José
func afficherPositionJose(tentativesRestantes int) {
	Printdessin((10 - tentativesRestantes) * 7)
}

// revelerLettres révèle n lettres aléatoires dans le mot
func revelerLettres(mot string) string {
	n := len(mot)/2 - 1
	if n < 0 {
		n = 0
	}

	lettresAReveler := make([]int, n)
	for i := 0; i < n; i++ {
		lettresAReveler[i] = rand.Intn(len(mot))
	}

	resultat := strings.Builder{}
	for i, char := range mot {
		if contains(lettresAReveler, i) {
			resultat.WriteRune(char)
		} else {
			resultat.WriteString("_")
		}
	}

	return resultat.String()
}

// revelerLettresCorrectes révèle toutes les occurrences de la lettre correcte dans le mot
func revelerLettresCorrectes(mot, devineActuelle, lettre string) string {
	resultat := strings.Builder{}
	for i, char := range mot {
		if string(char) == lettre {
			resultat.WriteString(lettre)
		} else {
			resultat.WriteString(string(devineActuelle[i]))
		}
	}
	return resultat.String()
}

// contains vérifie si la liste d'entiers contient une valeur donnée
func contains(liste []int, valeur int) bool {
	for _, elem := range liste {
		if elem == valeur {
			return true
		}
	}
	return false
}

// clearTerminal efface le terminal
func clearTerminal() {
	switch runtime.GOOS {
	case "linux", "darwin":
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

// Printdessin affiche les lignes i à i+6 du fichier hangman.txt
func Printdessin(i int) {
	file, err := os.Open("hangman.txt")
	if err != nil {
		fmt.Println("Erreur lors de l'ouverture du fichier:", err)
		return
	}
	defer file.Close()

	// Création d'un scanner pour lire le fichier ligne par ligne
	scanner := bufio.NewScanner(file)

	// Lire et afficher les lignes
	var lineCount int
	for scanner.Scan() {
		lineCount++
		if lineCount >= i && lineCount <= i+6 {
			fmt.Println(scanner.Text())
		}
		if lineCount >= i+6 {
			break
		}
	}
}

// SauvegarderJeu enregistre l'état actuel du jeu dans un fichier
func SauvegarderJeu(jeu *Pendu, cheminFichier string) error {
	fichier, err := os.Create(cheminFichier)
	if err != nil {
		return err
	}
	defer fichier.Close()

	encodeur := json.NewEncoder(fichier)
	return encodeur.Encode(jeu)
}

// ChargerJeu charge l'état du jeu depuis un fichier
func ChargerJeu(cheminFichier string) (*Pendu, error) {
	fichier, err := os.Open(cheminFichier)
	if err != nil {
		return nil, err
	}
	defer fichier.Close()

	decodeur := json.NewDecoder(fichier)
	var jeu Pendu
	err = decodeur.Decode(&jeu)
	if err != nil {
		return nil, err
	}

	return &jeu, nil
}

func main3() {
	fmt.Println("hangman c'est tout")
	cheminMots := "words.txt"
	var jeu *Pendu

	// Vérifier si l'indicateur --startWith est fourni
	if len(os.Args) == 3 && os.Args[1] == "--startWith" {
		var err error
		jeu, err = ChargerJeu(os.Args[2])
		if err != nil {
			fmt.Println("Erreur lors du chargement du jeu :", err)
			os.Exit(1)
		}
	} else {
		jeu = NouveauPendu(cheminMots)
	}

	for jeu.TentativesRestantes > 0 && strings.Contains(jeu.DevineActuelle, "_") {
		RendrePagePendu(jeu)

		var saisie string
		fmt.Scanln(&saisie)

		if saisie == "STOP" {
			err := SauvegarderJeu(jeu, "save.txt")
			if err != nil {
				fmt.Println("Erreur lors de la sauvegarde du jeu :", err)
				os.Exit(1)
			}
			fmt.Println("Jeu sauvegardé. Pour reprendre, utilisez '--startWith save.txt' comme indicateur.")
			os.Exit(0)
		}

		PenduHandler(jeu, saisie)
	}

	RendrePagePendu(jeu)
}
