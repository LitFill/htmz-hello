// htmz-hello, LitFill <author at email dot com>
// program for...
package main

import (
	"fmt"
	. "html/template"
	"log/slog"
	"net/http"
	"os"
)

var logger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
	// Level: slog.LevelDebug,
	AddSource: true,
}))

func mayFatal[T any](val T, err error) T {
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	return val
}

func wrapErr(err error, msg string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s, error: '%w'", msg, err)
}

func fatalWrap(err error, msg string) { mayFatal(0, wrapErr(err, msg)) }
func fatalWrapf(err error, format string, a ...any) {
	mayFatal(0, wrapErr(err, fmt.Sprintf(format, a...)))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	indexFile := mayFatal(os.ReadFile("./templat/index.html"))
	w.Write(indexFile)
}

func greetingHandlerWithTemplate(templat *Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nama := r.FormValue("name")
		// resp := struct{ Name string }{Name: nama}
		resp := map[string]string{"Name": nama}
		fatalWrapf(templat.Execute(w, resp), "Cannot execute in the /greeting handler")
	}
}

func incHandlerWithTemplate(s *Site, templat *Template) http.HandlerFunc {
	s.IncCount(1)
	return func(w http.ResponseWriter, r *http.Request) {
		fatalWrapf(templat.Execute(w, s), "Cannot execute in the /inc handler")
	}
}

func decHandlerWithTemplate(s *Site, templat *Template) http.HandlerFunc {
	s.DecCount(1)
	return func(w http.ResponseWriter, r *http.Request) {
		fatalWrapf(templat.Execute(w, s), "Cannot execute in the /dec handler")
	}
}

type Site struct {
	Count int
}

func (s *Site) GetCount() int { return s.Count }
func (s *Site) IncCount(by int) int {
	s.Count += by
	return s.Count
}
func (s *Site) DecCount(by int) int {
	s.Count -= by
	return s.Count
}

func main() {
	var site Site

	respTempl := mayFatal(ParseFiles("./templat/response.html"))
	counterTempl := Must(ParseFiles("./templat/counter.html"))

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("POST /greeting", greetingHandlerWithTemplate(respTempl))
	// http.HandleFunc("POST /greeting", greetingHandlerWithTemplate(mayFatal(ParseFiles("./templat/response.html"))))
	http.HandleFunc("POST /inc", func() http.HandlerFunc {
		site.IncCount(1)
		return incHandlerWithTemplate(&site, counterTempl)
	}())
	http.HandleFunc("POST /dec", func() http.HandlerFunc {
		site.DecCount(1)
		return decHandlerWithTemplate(&site, counterTempl)
	}())

	logger.Info("Listening to :8854")
	http.ListenAndServe(":8854", nil)
}
