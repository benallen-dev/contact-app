package handlers

import (
	"net/http"

	"github.com/a-h/templ"

	"github.com/benallen-dev/contact-app/pkg/archiver"
	"github.com/benallen-dev/contact-app/views"
)

var arkive archiver.Archiver

func init() {
	arkive = *archiver.NewArchiver()
}

func PostArchive(w http.ResponseWriter, r *http.Request) {
	arkive.Run()

	templ.Handler(views.Archive(&arkive)).ServeHTTP(w, r)
}

func GetArchive(w http.ResponseWriter, r *http.Request) {
	templ.Handler(views.Archive(&arkive)).ServeHTTP(w, r)
}

func GetArchiveFile(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "data/contacts.csv")
}
