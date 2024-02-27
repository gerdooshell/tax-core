package restApi

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gerdooshell/tax-core/controller/rest_api/handlers"
	optimalRrsp "github.com/gerdooshell/tax-core/controller/rest_api/handlers/optimal_rrsp"
	taxCalculator "github.com/gerdooshell/tax-core/controller/rest_api/handlers/tax_calculator"
	taxMargin "github.com/gerdooshell/tax-core/controller/rest_api/handlers/tax_margin"

	logger "github.com/gerdooshell/tax-logger-client-go"
	"github.com/gorilla/mux"
)

func ServeHTTP() {
	muxRouter := mux.NewRouter()
	apiHandlers := make([]handlers.Handler, 0, 5)
	apiHandlers = append(apiHandlers, taxCalculator.NewTaxCalculatorController())
	apiHandlers = append(apiHandlers, taxMargin.NewTaxMarginController())
	apiHandlers = append(apiHandlers, optimalRrsp.NewOptimalRRSPController())
	for _, handler := range apiHandlers {
		RegisterMuxHTTP(muxRouter, handler)
	}
	launchHTTPServerMux(muxRouter)
	//launchHTTPServer(muxRouter)
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

func RegisterMuxHTTP(muxRouter *mux.Router, handler handlers.Handler) {
	h := createHTTPHandler(handler)
	muxRouter.HandleFunc(handler.URL(), func(w http.ResponseWriter, r *http.Request) {
		notAllowed := true
		for _, method := range append(handler.Methods(), http.MethodOptions) {
			if method == r.Method {
				w.Header().Add("Access-Control-Allow-Origin", "*")
				w.Header().Add("Access-Control-Allow-Methods", "GET, OPTIONS")
				w.Header().Add("Cache-Control", "public, max-age=600")
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

func launchHTTPServer(handler http.Handler) {
	addrFlag := flag.String("addr", ":8185", "address to run unified server on")
	fmt.Println("http server ready")
	_ = http.ListenAndServe(*addrFlag, handler)
}

func launchHTTPServerMux(handler http.Handler) {
	server := &http.Server{
		Addr:              ":8185",
		Handler:           handler,
		ReadTimeout:       time.Second * 2,
		ReadHeaderTimeout: time.Second * 2,
		WriteTimeout:      time.Second * 5,
		IdleTimeout:       time.Second * 5,
	}
	fmt.Println("http server ready")
	_ = server.ListenAndServe()
}

func createHTTPHandler(handler handlers.Handler) http.Handler {
	hf := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				logger.FatalFormat("Recovered fatal error: \"%v\"", rec)
			}
		}()
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		logger.InfoFormat("request:\"%v\" from \"%v\"", r.URL, r.RemoteAddr)
		err := handler.ParseArgs(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			message := "failed submitting the message"
			if _, err = w.Write([]byte(fmt.Sprintf("%v. error: %v", message, err))); err != nil {
				logger.FatalFormat("Error writing response body to writer. err: %s", err.Error())
			}
			return
		}
		if err = handler.Authorize(); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		resp := handler.Process(r)
		w.WriteHeader(resp.StatusCode)
		if resp != nil && resp.Body != nil {
			defer resp.Body.Close()
			body, _ := io.ReadAll(resp.Body)
			if _, err = w.Write(body); err != nil {
				logger.FatalFormat("Error writing response body to writer. err: %s", err.Error())
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
