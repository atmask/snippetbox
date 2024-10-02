package main

import (
    "errors"
    "fmt"
    "net/http"
    "strconv"
    "strings"
    "unicode/utf8"

    "github.com/atmask/snippetbox/internal/models"
)

// Change the signature of the home handler so it is defined as a method against
// *application.
func (app *application) home(w http.ResponseWriter, r *http.Request) {
    snippets, err := app.snippets.Latest()
    if err != nil {
        app.serverError(w,r, err)
        return 
    }

    data := app.newTemplateData(r)
    data.Snippets = snippets

    app.render(w, r, http.StatusOK, "home.tmpl.html", data)
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

    data := app.newTemplateData(r)
    data.Snippet = snippet

    app.render(w, r, http.StatusOK, "view.tmpl.html", data)
}

// Change the signature of the snippetCreate handler so it is defined as a method
// against *application.
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
    data := app.newTemplateData(r)

    app.render(w, r, http.StatusOK, "create.tmpl.html", data)
}

// Change the signature of the snippetCreatePost handler so it is defined as a method
// against *application.
func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
    err := r.ParseForm()
    if err != nil {
        app.clientError(w, http.StatusBadRequest)
        return
    }

    title := r.PostForm.Get("title")
    content := r.PostForm.Get("content")

    // PostForm.Get() always returns a string so we need to convert
    // to an int
    expires, err := strconv.Atoi(r.PostForm.Get("expires"))
    if err != nil {
        app.clientError(w, http.StatusBadRequest)
        return
    }

    // Init a map to hold any validation errors
    fieldErrors := make(map[string]string)
    
    // Validate the title field
    if strings.TrimSpace(title) == "" {
        fieldErrors["title"] = "This field cannot be blank"
    } else if utf8.RuneCountInString(title) > 100 {
        fieldErrors["title"] = "This field cannot be more than 100 characters long"
    }

    // Check that the Content value is not blank
    if strings.TrimSpace(content) == "" {
        fieldErrors["content"] = "This field cannot be blank"
    }

    // Check that expires value matches one of the accepted values (1 ,7, or 365)
    if expires != 1 && expires != 7 && expires != 365 {
        fieldErrors["expires"] = "This field must be set to 1, 7, or 365"
    }

    // If there are any errors, dump them in a plain text http response
    // and return from the handlet
    if len(fieldErrors) > 0 {
        fmt.Fprint(w, fieldErrors)
        return
    }

    id, err := app.snippets.Insert(title, content, expires)
    if err != nil {
        app.serverError(w, r, err)
        return
    }

    http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}