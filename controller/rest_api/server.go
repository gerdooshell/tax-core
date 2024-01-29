package restApi

import (
	"flag"
	"fmt"
	"github.com/gerdooshell/tax-core/controller/rest_api/handlers"
	taxCalculator "github.com/gerdooshell/tax-core/controller/rest_api/handlers/tax_calculator"
	"io"
	"log"
	"net/http"
)

func ServeHTTP() {
	apiHandlers := make([]handlers.Handler, 0, 3)
	apiHandlers = append(apiHandlers, taxCalculator.NewTaxCalculatorController())
	for _, handler := range apiHandlers {
		RegisterHTTP(handler)
	}
	launchHTTPServer()
}

func RegisterHTTP(handler handlers.Handler) {
	h := createHTTPHandler(handler)
	http.HandleFunc(handler.URL(), func(w http.ResponseWriter, r *http.Request) {
		notAllowed := true
		for _, method := range append(handler.Methods(), http.MethodOptions) {
			if method == r.Method {
				w.Header().Add("Access-Control-Allow-Origin", "*")
				w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, x-api-key")
				h.ServeHTTP(w, r)
				notAllowed = false
			}
		}
		if notAllowed {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
		return
	})
}

func launchHTTPServer() {
	addrFlag := flag.String("addr", ":8185", "address to run unified server on")
	fmt.Println("http server ready")
	_ = http.ListenAndServe(*addrFlag, nil)
}

func createHTTPHandler(handler handlers.Handler) http.Handler {
	hf := func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		request, err := handler.ParseArgs(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			message := "failed submitting the message"
			if _, err = w.Write([]byte(fmt.Sprintf("%v. error: %v", message, err))); err != nil {
				log.Fatalf("Error writing response body to writer. err: %s", err.Error())
			}
			return
		}
		if err = handler.Authorize(); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		resp := handler.Process(request)
		w.WriteHeader(resp.StatusCode)
		if resp != nil && resp.Body != nil {
			defer resp.Body.Close()
			body, _ := io.ReadAll(resp.Body)
			if _, err = w.Write(body); err != nil {
				log.Fatalf("Error writing response body to writer. err: %s", err.Error())
			}
		}
		return
	}
	return httpHandler{
		serveHTTP: hf,
	}
}

// httpHandler fulfills the http.Handler interface, allowing us to use logging http middleware
type httpHandler struct {
	serveHTTP http.HandlerFunc
}

// ServeHTTP calls through to the constructed function
func (h httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.serveHTTP(w, r)
}
