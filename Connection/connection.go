package Connection

import (
	"html/template"
	"log"
	"net/http"
)

// we need to call Handlefunc to every link to render these pages
func Connection() {
	http.HandleFunc("/home.html", Render_Game)
}

func Render_Game(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("logged-in")
	if err != nil {
		http.SetCookie(w, &http.Cookie{
			Name:  "logged-in",
			Value: "0",
			Path:  "/",
		})
	}
	parsedTemplate, _ := template.ParseFiles("./template/home.html")
	err_tmpl := parsedTemplate.Execute(w, c)
	if err_tmpl != nil {
		log.Println("Error executing template :", err_tmpl)
		return
	}
}
