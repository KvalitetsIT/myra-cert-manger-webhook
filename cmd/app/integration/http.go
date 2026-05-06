package main_test

import (
	"bytes"
	"encoding/json"

	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/KvalitetsIT/cert-manager-webhook-myra/internal/testutil"
	"github.com/Myra-Security-GmbH/myrasec-go/v2"
	"github.com/gorilla/mux"
)

// The purpose of this handler is to mimic the real Myra API
type MyraHttpHandler struct {
	store  *testutil.Storage
	logger *slog.Logger
}

func newMyraHttpHandler(store *testutil.Storage, logger *slog.Logger) *MyraHttpHandler {
	return &MyraHttpHandler{
		store:  store,
		logger: logger,
	}
}

func (handler *MyraHttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// --- Read request body as before ---
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		handler.logger.Error("Failed to read body", slog.Any("error", err))
		http.Error(w, "failed to read body", http.StatusInternalServerError)
		return
	}
	r.Body.Close()

	var bodyMap map[string]any
	if len(bodyBytes) > 0 {
		if err := json.Unmarshal(bodyBytes, &bodyMap); err != nil {
			handler.logger.Warn("Body is not valid JSON", slog.String("body", string(bodyBytes)), slog.Any("error", err))
		}
	}
	handler.logger.Info("Incoming request", slog.String("path", r.URL.Path), slog.String("method", r.Method), slog.Any("body", bodyMap))
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	// --- Serve through mux ---
	router := mux.NewRouter()
	router.HandleFunc("/domain/{domainID}/dns-records", handler.create).Methods("POST")
	router.HandleFunc("/domain/{domainID}/dns-records/{recordID}", handler.delete).Methods("DELETE")
	router.HandleFunc("/domain/{domainID}/dns-records", handler.get_records).Methods("GET")
	router.HandleFunc("/domains", handler.get_domains).Methods("GET")

	router.ServeHTTP(w, r)

}

func (m *MyraHttpHandler) create(w http.ResponseWriter, r *http.Request) {
	var record myrasec.DNSRecord
	if err := json.NewDecoder(r.Body).Decode(&record); err != nil {
		m.logger.Error("Failed to decode request body", slog.Any("error", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		panic(err)
	}

	vars := mux.Vars(r)
	domainId := m.getId("domainID", vars)

	record, err := m.store.AddRecord(domainId, record)
	if err != nil {
		m.logger.Error("Failed to add record", slog.Any("record", record), slog.Any("error", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		panic(err)

	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201

	err = json.NewEncoder(w).Encode(record)

	if err != nil {
		m.logger.Error("Failed to write response", slog.Any("error", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (m *MyraHttpHandler) getId(key string, vars map[string]string) int {
	id, err := strconv.Atoi(vars[key])
	if err != nil {
		m.logger.Error("Failed to parse key", slog.String("key", key))
		panic(err)
	}
	return id
}

func (m *MyraHttpHandler) delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	domainID := m.getId("domainID", vars)
	recordID := m.getId("recordID", vars)

	record, err := m.store.DeleteRecord(domainID, recordID)
	if err != nil {
		m.logger.Error("Failed to delete record", slog.Any("domain id", domainID), slog.Any("record id", recordID))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 204
	err = json.NewEncoder(w).Encode(record)

	if err != nil {
		m.logger.Error("Failed to encode response", slog.Any("error", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (m *MyraHttpHandler) get_records(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	domainID := m.getId("domainID", vars)
	records, err := m.store.GetRecords(domainID)

	resp := newResponse(records)

	if err != nil {
		m.logger.Error("Could not provide records", slog.Any("domain id", domainID))
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(resp)

	if err != nil {
		m.logger.Error("Failed to encode response", slog.Any("error", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (m *MyraHttpHandler) get_domains(w http.ResponseWriter, req *http.Request) {
	type DomainVO struct {
		ObjectType   string              `json:"objectType"`
		Name         string              `json:"name"`
		Organization int                 `json:"organizationId"`
		AutoUpdate   bool                `json:"autoUpdate"`
		Maintenance  bool                `json:"maintenance"`
		DnsRecords   []myrasec.DNSRecord `json:"dnsRecords"`
		Reversed     bool                `json:"reversed"`
		Environment  string              `json:"environment"`
		Locked       bool                `json:"locked"`
		Id           int                 `json:"id"`
		Modified     string              `json:"modified"`
		Created      string              `json:"created"`
		Label        string              `json:"label"`
	}

	// Build domains from the store
	domains := []DomainVO{}
	for _, domain := range m.store.GetDomains() { // helper we'll add
		records, err := m.store.GetRecords(domain.ID)
		if err != nil {
			m.logger.Error("Could not retrieve recods for domain", slog.Any("Domain id", domain.ID), slog.Any("error", err))
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		domain := DomainVO{
			ObjectType:   "DomainVO",
			Name:         domain.Name,
			Organization: 1000430, // hardcode or track if you want
			AutoUpdate:   true,
			Maintenance:  false,
			DnsRecords:   records,
			Reversed:     false,
			Environment:  "live",
			Locked:       false,
			Id:           domain.ID,
			Modified:     "2025-10-15T15:57:38+0200",
			Created:      "2025-10-15T15:57:38+0200",
			Label:        domain.Name,
		}
		domains = append(domains, domain)
	}

	resp := newResponse(domains)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(resp)

	if err != nil {
		m.logger.Error("Failed to write response", slog.Any("error", err))
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
}

func newResponse[T any](entries []T) myrasec.Response {
	data := make([]any, len(entries))
	for i, e := range entries {
		data[i] = e
	}
	return myrasec.Response{
		Error:         false,
		ViolationList: make([]*myrasec.Violation, 0),
		WarningList:   make([]*myrasec.Warning, 0),
		Data:          data,
		Page:          1,
		Count:         len(data),
		PageSize:      50,
	}
}
