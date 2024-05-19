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

type Site struct {
	Count int
}

func (s *Site) GetCount() int { return s.Count }
func (s *Site) IncCount(by int) int {
	s.Count += by
	return s.GetCount()
}
func (s *Site) DecCount(by int) int {
	s.Count -= by
	return s.Count
}

var site Site

type Data[T comparable] map[string]T

func (d *Data[T]) add(key string, val T) *Data[T] {
	(*d)[key] = val
	return d
}
func newData[T comparable]() *Data[T]                { return &Data[T]{} }
func mkdata[T comparable](key string, val T) Data[T] { return *newData[T]().add(key, val) }

var logger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
	Level:     slog.LevelDebug,
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

func fatalLog(err error, msg string, log ...any) {
	if err == nil {
		return
	}
	logger.Error(msg, log...)
	os.Exit(1)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	templat := Must(ParseFiles("./templat/index.html"))
	fatalLog(
		templat.Execute(w, site),
		"Cannot execute in the indexHandler",
		"templat", templat,
		"data", site,
	)
}

func greetHandler(w http.ResponseWriter, r *http.Request) {
	data := mkdata("Name", r.FormValue("name"))
	templat := Must(ParseFiles("./templat/response.html"))
	fatalLog(
		templat.Execute(w, data),
		"Cannot execute in the /greeting handler",
		"templat", templat,
		"data", data,
	)
}

func incHandler(w http.ResponseWriter, _ *http.Request) {
	// data := mkdata("Count", site.IncCount(1))
	// templat := Must(ParseFiles("./templat/counter.html"))
	// fatalLog(
	// 	templat.Execute(w, data),
	// 	"Cannot execute in the /inc handler",
	// 	"templat", templat,
	// 	"data", data,
	// )
	w.Write([]byte("<p id=\"counter\">1321</p>"))
}

func decHandler(w http.ResponseWriter, _ *http.Request) {
	data := mkdata("Count", site.DecCount(1))
	templat := Must(ParseFiles("./templat/counter.html"))
	fatalLog(
		templat.Execute(w, data),
		"Cannot execute in the /dec handler",
		"templat", templat,
		"data", data,
	)
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

func main() {
	logger.Debug("variable site", "site", site)

	// counterTempl := Must(ParseFiles("./templat/counter.html"))

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("POST /greeting", greetHandler)
	http.HandleFunc("POST /inc", incHandler)
	http.HandleFunc("POST /dec", decHandler)
	// http.HandleFunc("POST /inc", func() http.HandlerFunc {
	// 	site.IncCount(1)
	// 	return incHandlerWithTemplate(&site, counterTempl)
	// }())
	// http.HandleFunc("/dec", func() http.HandlerFunc {
	// 	site.DecCount(1)
	// 	return decHandlerWithTemplate(&site, counterTempl)
	// }())

	// logger.Debug("sebeblum masuk /inc")
	// http.HandleFunc("POST /inc", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	site.IncCount(1)
	// 	logger.Debug("var site after .Inc", "site", site)
	// 	fatalWrapf(counterTempl.Execute(w, map[string]int{"Count": site.IncCount(1)}), "Cannot execute in the /inc handler")
	// }))
	// http.HandleFunc("POST /dec", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	site.DecCount(1)
	// 	fatalWrapf(counterTempl.Execute(w, site), "Cannot execute in the /dec handler")
	// }))

	logger.Info("Listening to :8854")
	http.ListenAndServe(":8854", nil)
}
