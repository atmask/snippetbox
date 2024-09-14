package main

import (
    "errors"
    "fmt"
    // "html/template"
    "net/http"
    "strconv"

    "github.com/atmask/snippetbox/internal/models"
)

// Change the signature of the home handler so it is defined as a method against
// *application.
func (app *application) home(w http.ResponseWriter, r *http.Request) {
    w.Header().Add("Server", "Go")

    snippets, err := app.snippets.Latest()
    if err != nil {
        app.serverError(w,r, err)
        return 
    }

    for _, snippet := range snippets {
        fmt.Fprintf(w, "%v\n", snippet)
    }

    // files := []string{
    //     "./ui/html/base.tmpl.html",
    //     "./ui/html/partials/nav.tmpl.html",
    //     "./ui/html/pages/home.tmpl.html",
    // }

    // ts, err := template.ParseFiles(files...)
    // if err != nil {
    //     app.serverError(w, r, err) // Use the serverError() helper.
    //     return
    // }

    // err = ts.ExecuteTemplate(w, "base", nil)
    // if err != nil {
	// 	app.serverError(w, r, err)       
    // }
}

// Change the signature of the snippetView handler so it is defined as a method
// against *application.
func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(r.PathValue("id"))
    if err != nil || id < 1 {
        http.NotFound(w, r)
        return
    }

    snippet, err := app.snippets.Get(id)
    if err != nil {
        if errors.Is(err, models.ErrNoRecord) {
            http.NotFound(w,r)
        } else {
            app.serverError(w,r, err)
        }
        return
    }

    fmt.Fprintf(w, "%+v", snippet)
}

// Change the signature of the snippetCreate handler so it is defined as a method
// against *application.
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Display a form for creating a new snippet..."))
}

// Change the signature of the snippetCreatePost handler so it is defined as a method
// against *application.
func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
    title  := "O nsail"
    content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
    expires := 7    

    id, err := app.snippets.Insert(title, content, expires)
    if err != nil {
        app.serverError(w, r, err)
        return
    }

    http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
    
    w.WriteHeader(http.StatusCreated)
    w.Write([]byte("Save a new snippet..."))
}