package serve_mux

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// ServeMux returns a server multiplexer
// reference:
// https://go.dev/blog/routing-enhancements
// https://pkg.go.dev/net/http#ServeMux
// https://github.com/golang/go/issues/61410
func ServeMux() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /users/{id}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("QUERYING USER... %s\n", r.PathValue("id"))
		fmt.Println("DONE")
	})

	mux.HandleFunc("POST /users", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("CREATING USER...")
		fmt.Println("DONE")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`User created: 5`))
	})

	mux.HandleFunc("DELETE /users/{id}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("DELETING USER... %s\n", r.PathValue("id"))
		fmt.Println("DONE")
	})

	mux.HandleFunc("PUT /users/{id...}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("UPDATING USER... %s\n", r.PathValue("id"))
		fmt.Println("DONE")
	})

	mux.HandleFunc("PUT /users/{id}/{attribute}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("UPDATING USER... %s\n", r.PathValue("id"))
		fmt.Printf("UPDATING USER... %s\n", r.PathValue("attribute"))
		fmt.Println("DONE")
	})

	// unacceptable path {id...} must at tail
	// mux.HandleFunc("GET /org/{id...}/commit", func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Printf("QUERYING ORG...%s\n", r.PathValue("id"))
	// })

	// unacceptable path {id...} must at tail
	// mux.HandleFunc("GET /v2/org/{id...}/commit/{version}", func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Printf("QUERYING ORG...%s\n", r.PathValue("id"))
	// })

	return mux
}

func TestMux(t *testing.T) {
	srv := httptest.NewServer(ServeMux())
	defer srv.Close()

	c := srv.Client()

	req, _ := http.NewRequest(http.MethodGet, srv.URL+"/users/1", nil)
	_, _ = c.Do(req)

	req, _ = http.NewRequest(http.MethodPost, srv.URL+"/users", nil)
	_, _ = c.Do(req)

	req, _ = http.NewRequest(http.MethodDelete, srv.URL+"/users/1", nil)
	_, _ = c.Do(req)

	req, _ = http.NewRequest(http.MethodPut, srv.URL+"/users/1/2/3", nil)
	_, _ = c.Do(req)

	// panic: parsing "GET /org/{id...}/commit": at offset 9: {...} wildcard not at end [recovered, repanicked]
	// req, _ = http.NewRequest(http.MethodGet, srv.URL+"/org/1/2/3/commit", nil)
	// _, _ = c.Do(req)

	// panic: parsing "GET /v2/org/{id...}/commit/{version}": at offset 12: {...} wildcard not at end [recovered, repanicked]
	req, _ = http.NewRequest(http.MethodGet, srv.URL+"/v2/org/1/2/3/commit/v2", nil)
	_, _ = c.Do(req)
}
