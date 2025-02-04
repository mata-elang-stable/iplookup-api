package app

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/fadhilyori/iplookup-go/internal/cache"
	"github.com/fadhilyori/iplookup-go/internal/iplookup"
	"github.com/fadhilyori/iplookup-go/internal/logger"
	"github.com/fadhilyori/iplookup-go/internal/schema"
	"github.com/gorilla/mux"
	"net"
	"net/http"
	"strings"
	"time"
)

var log = logger.GetLogger()

type App struct {
	regionLookup  *iplookup.IPLookup
	asnLookup     *iplookup.IPLookup
	listenAddr    string
	listenPort    int
	cacheInstance *cache.ValKeyInstance
}

func NewApp(listenAddr string, listenPort int) *App {
	return &App{
		listenAddr: listenAddr,
		listenPort: listenPort,
	}
}

func (a *App) EnableCache(addresses []string, ttl time.Duration) {
	a.cacheInstance = cache.MustNewValkey(addresses, ttl)
	log.Infof("cache enabled with redis backend: %v", addresses)
}

func (a *App) LoadRegionMaxmindDB(dbPath string) {
	a.regionLookup = iplookup.NewIPLookup(dbPath)
}

func (a *App) LoadASNMaxmindDB(dbPath string) {
	a.asnLookup = iplookup.NewIPLookup(dbPath)
}

func sendResponse(w http.ResponseWriter, r *http.Request, code int, data []byte) error {
	log.WithFields(logger.Fields{
		"client": strings.Split(r.RemoteAddr, ":")[0],
		"method": r.Method,
		"path":   r.URL.RawPath,
		"code":   code,
	}).Infoln()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if _, err := w.Write(data); err != nil {
		return err
	}

	return nil
}

func sendSuccessResponse(w http.ResponseWriter, r *http.Request, data []byte) {
	if err := sendResponse(w, r, http.StatusOK, data); err != nil {
		log.Errorf("Error sending response: %v", err)
	}
}

func sendErrorResponse(w http.ResponseWriter, r *http.Request, code int, message string) {
	responseText := fmt.Sprintf(`{"error": "%s"}`, message)

	if err := sendResponse(w, r, code, []byte(responseText)); err != nil {
		log.Errorf("Error sending response: %v", err)
	}
}

func (a *App) handleHealthCheckRequest(w http.ResponseWriter, r *http.Request) {
	log.WithField("client", r.RemoteAddr).Debugf("Received request: %s %s", r.Method, r.URL.Path)

	sendSuccessResponse(w, r, []byte(`{"status": "ok"}`))
}

func (a *App) handleLookupRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ip := vars["ip"]

	// Validate the IP address
	parsedIP := net.ParseIP(ip)

	if parsedIP == nil {
		sendErrorResponse(w, r, http.StatusBadRequest, "invalid IP address")
		return
	}

	if parsedIP.IsPrivate() {
		sendErrorResponse(w, r, http.StatusBadRequest, "private IP address")
		return
	}

	if a.cacheInstance != nil {
		// Check if the request already cached
		if val, err := a.cacheInstance.Get(r.Context(), ip); err == nil {
			log.Debugf("Serving request from cache: %s", ip)

			sendSuccessResponse(w, r, []byte(val))
			return
		}
	}

	var regionData schema.MaxmindDBRegion
	var asnData schema.MaxmindDBAS

	if a.regionLookup != nil {
		if err := a.regionLookup.Lookup(parsedIP, &regionData); err != nil {
			// if error is "invalid IP address" or "private IP address", return 400

			sendErrorResponse(w, r, http.StatusInternalServerError, err.Error())
			return
		}
	}

	if a.asnLookup != nil {
		if err := a.asnLookup.Lookup(parsedIP, &asnData); err != nil {
			sendErrorResponse(w, r, http.StatusInternalServerError, err.Error())
			return
		}
	}

	responseDataByte, _ := json.Marshal(map[string]interface{}{
		"region": regionData,
		"asn":    asnData,
	})

	if a.cacheInstance != nil {
		log.Debugf("Caching request: %s", ip)
		if err := a.cacheInstance.Cache(r.Context(), ip, string(responseDataByte)); err != nil {
			log.Errorf("Error while caching request: %v", err)
		}
	}

	sendSuccessResponse(w, r, responseDataByte)
}

func (a *App) Run(ctx context.Context) error {
	listenAddr := fmt.Sprintf("%s:%d", a.listenAddr, a.listenPort)

	log.Printf("Starting server on %s", listenAddr)

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/lookup", a.handleLookupRequest).Methods("GET").Queries("ip", "{ip}")
	router.HandleFunc("/health", a.handleHealthCheckRequest).Methods("GET")

	server := &http.Server{
		Addr:         listenAddr,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		BaseContext: func(net.Listener) context.Context {
			return ctx
		},
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Errorf("%v", err)
		}
	}()

	<-ctx.Done()
	log.Info("Shutting down server")

	return server.Close()
}
