package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"text/template"
	"time"
)

type PageData struct {
	Difficulty string
	FileName   string
	ClassName  string
}

// Pendu représente l'état du jeu du Pendu
type Pendu struct {
	MotADeviner           string
	DevineActuelle        string
	DevinettesIncorrectes []string
	TentativesRestantes   int
	LettresDevinees       []string
	MotsDevines           []string
}

// Game_Web_Data représente toutes les données à afficher dans le jeu
type Game_Web_Data struct {
	MotADeviner           string
	DevineActuelle        string
	DevinettesIncorrectes []string
	TentativesRestantes   int
	LettresDevinees       []string
	MotsDevines           []string
	HangmanDessin         string
	Difficulty            string
	DevineSaisi           string
}

const port = ":8080"

var Pendu_Web_Data *Game_Web_Data

// NouveauPendu initialise un nouveau jeu du Pendu avec un mot à deviner
func NouveauPendu(cheminMots string) *Game_Web_Data {
	motAleatoire := choisirMotAleatoireDepuisFichier(cheminMots)

	return &Game_Web_Data{
		MotADeviner:           motAleatoire,
		DevineActuelle:        revelerLettres(motAleatoire),
		DevinettesIncorrectes: []string{},
		TentativesRestantes:   10,
		LettresDevinees:       []string{},
		MotsDevines:           []string{},
	}
}

// func ChargerPendu(GameData *Game_Web_Data) *Game_Web_Data {

// 	// GameData.DevineActuelle
// 	// return &Game_Web_Data{
// 	// 	DevineActuelle:        revelerLettres(motAleatoire),
// 	// 	DevinettesIncorrectes: []string{},
// 	// 	TentativesRestantes:   10,
// 	// 	LettresDevinees:       []string{},
// 	// 	MotsDevines:           []string{},

// 	// }
// }

// RendrePagePendu rend la page du jeu du Pendu
func RendrePagePendu(jeu *Pendu) {

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
func PenduHandler(jeu *Game_Web_Data, saisie string) {
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
func (jeu *Game_Web_Data) DejaDevine(saisie string) bool {
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
()

// func (jeu *Pendu) DejaDevine(saisie string) bool {
// 	for _, lettre := range jeu.LettresDevinees {
// 		if lettre == saisie {
// 			return true
// 		}
// 	}

// 	for _, mot := range jeu.MotsDevines {
// 		if mot == saisie {
// 			return true
// 		}
// 	}

// 	return false
// }

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
	// Printdessin((10 - tentativesRestantes) * 7)
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

// // Printdessin affiche les lignes i à i+6 du fichier hangman.txt
// func Printdessin(i int) {
// 	file, err := os.Open("hangman.txt")
// 	if err != nil {
// 		fmt.Println("Erreur lors de l'ouverture du fichier:", err)
// 		return
// 	}
// 	defer file.Close()

// 	// Création d'un scanner pour lire le fichier ligne par ligne
// 	scanner := bufio.NewScanner(file)

//		// Lire et afficher les lignes
//		var lineCount int
//		for scanner.Scan() {
//			lineCount++
//			if lineCount >= i && lineCount <= i+6 {
//				fmt.Println(scanner.Text())
//			}
//			if lineCount >= i+6 {
//				break
//			}
//		}
//	}
//
// Printdessin affiche les lignes i à i+6 du fichier hangman.txt
func Printdessin_WEB(i int) string {
	displayed_hangmanText := ""
	file, err := os.Open("hangman.txt")
	if err != nil {
		fmt.Println("Erreur lors de l'ouverture du fichier:", err)
		return ""
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
			displayed_hangmanText += scanner.Text()
		}
		if lineCount >= i+6 {
			break
		}
	}
	return displayed_hangmanText
}

// ChargerJeu charge l'état du jeu depuis un fichier
func ChargerJeu(cheminFichier string) (*Game_Web_Data, error) {
	fichier, err := os.Open(cheminFichier)
	if err != nil {
		return nil, err
	}
	defer fichier.Close()

	decodeur := json.NewDecoder(fichier)
	var jeu Game_Web_Data
	err = decodeur.Decode(&jeu)
	if err != nil {
		return nil, err
	}

	return &jeu, nil
}

func Home(w http.ResponseWriter, r *http.Request) {

	Difficulties := []PageData{
		{"Easy", "words.txt", "easy"},
		{"Medium", "words2.txt", "medium"},
		{"Hard", "words3.txt", "hard"},
	}

	// Parse the HTML template file
	tmpl, err := template.ParseFiles("html/home.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Execute the template with the data and write the result to the response writer
	err = tmpl.Execute(w, Difficulties)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

}
func bienvenue(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "html/bienvenue.html")
}
func Game(w http.ResponseWriter, r *http.Request) {
	difficulty := r.FormValue("difficulty_value")
	fmt.Println("diff : ", difficulty)
	fileName := ""
	switch difficulty {
	case "Easy":
		fileName = "words.txt"
	case "Medium":
		fileName = "words2.txt"
	case "Hard":
		fileName = "words3.txt"
	}
	Pendu_Web_Data := &Game_Web_Data{
		Difficulty: fileName,
	}

	// ********************* A NE PAS DECOMMENTER **************************** // DEJA DECLAREE EN HAUT
	// type Game_Web_Data struct {
	// 	MotADeviner           string
	// 	DevineActuelle        string
	// 	DevinettesIncorrectes []string
	// 	TentativesRestantes   int
	// 	LettresDevinees       []string
	// 	MotsDevines           []string
	// 	HangmanDessin         string
	// 	Difficulty            string
	// 	DevineSaisi           string
	// }
	// ********************* A NE PAS DECOMMENTER **************************** //
	//all of the Form Values
	Mot_saisi := r.PostFormValue("user_guess_value")  // ce que l'utilisateur a saisi
	bouton_envoyer := r.PostFormValue("submit_guess") // le bouton de submit (Envoyer)
	// un booléen pour dire si on a cliqué sur le bouton submit ou pas
	bouton_cliquer := bouton_envoyer != ""
	// vérifier si on a cliqué sur le bouton
	if bouton_cliquer {
		//Pendu* -> Game_Web_Data*
		//Jeu -> Pendu_Web_Data
		Pendu_Web_Data.DevineSaisi = Mot_saisi
		// 1. prendre ce que l'utilisateur a saisi (la variable Mot_saisi) et faire le traitement necéssaire avec cette variable
		if len(Mot_saisi) == 1 {
			//si il a déjà été devinée
			if Pendu_Web_Data.DejaDevine(Mot_saisi) {
				fmt.Println("le carrecteur est correct")
			} else if strings.Contains(Pendu_Web_Data.MotADeviner, Mot_saisi) {
				// si la lettre est correscte
				Pendu_Web_Data.DevineActuelle = revelerLettresCorrectes(Pendu_Web_Data.MotADeviner, Pendu_Web_Data.DevineActuelle, Mot_saisi)
				Pendu_Web_Data.LettresDevinees = append(Pendu_Web_Data.LettresDevinees, Mot_saisi)
			} else {
				//faux
				Pendu_Web_Data.TentativesRestantes--
				Pendu_Web_Data.DevinettesIncorrectes = append(Pendu_Web_Data.DevinettesIncorrectes, Mot_saisi)
				Pendu_Web_Data.LettresDevinees = append(Pendu_Web_Data.LettresDevinees, Mot_saisi)
			}
		} else if len(Mot_saisi) > 1 {
			// si

		}
	}

	// 2. mettre les valeurs après le traitement de Mot_saisi dans Pendu_Web_Data (exemple Pendu_Web_Data.DevineActuelle = "hello") pour bien sauvegarder les données
	// 3. dans Pendu_Web_Data y'a HangmanDessin dans lequel vous metterezle dessin du hangman en accord avec les tentatives (utilisez  Printdessin_WEB(i int)// EXEMPLE : Pendu_Web_Data.HangmanDessin 0= Printdessin_WEB(3))
	// 4. ATTENTION !!! quand vous mettez lettresIncorrectes faut leur faire un join (faut pas avoir un tableau mais un string) (exemple : ["a","b"] -> NON ,"a,b"-> OUI)

	// } else if  {
	// 	// Pendu_Web_Data = NouveauPendu("words3.txt")
	// }
	if len(Mot_saisi) > 1 {
		fmt.Print("le mot ecri sur le web : ", Mot_saisi)
	} else {
		fmt.Print("la lettre ecri sur le web : ", Mot_saisi)
	}
	// Parse the HTML template file
	tmpl, err := template.ParseFiles("html/game.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Execute the template with the data and write the result to the response writer
	err = tmpl.Execute(w, Pendu_Web_Data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func main() {
	http.HandleFunc("/home", Home)
	http.HandleFunc("/game", Game)
	fmt.Println("hangman c'est tout")
	// cheminMots := "words.txt"
	// var jeu *Pendu

	// // Vérifier si l'indicateur --startWith est fourni
	// if len(os.Args) == 3 && os.Args[1] == "--startWith" {
	// 	var err error
	// 	jeu, err = ChargerJeu(os.Args[2])
	// 	if err != nil {
	// 		fmt.Println("Erreur lors du chargement du jeu :", err)
	// 		os.Exit(1)
	// 	}
	// } else {
	// 	jeu = NouveauPendu(cheminMots)
	// }

	// for jeu.TentativesRestantes > 0 && strings.Contains(jeu.DevineActuelle, "_") {
	// 	RendrePagePendu(jeu)

	// 	var saisie string
	// 	fmt.Scanln(&saisie)

	// 	PenduHandler(jeu, saisie)
	// }

	// RendrePagePendu(jeu)
	fmt.Println("\n(http://localhost:8080/home) - Server started on port", port)
	http.ListenAndServe(port, nil)
}
