package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type CustomerHandler struct {
	service CustomerService
}

type Subscriber interface {
	getSubscriberId() string
	update(customers []Customer)
}

func NewCustomerHandler(service CustomerService) *CustomerHandler {
	return &CustomerHandler{service: service}
}

func handleResponseErr(w http.ResponseWriter, statusCode int, errMsg string, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	log.Printf("failed due to error :%q", err)

	if err := json.NewEncoder(w).Encode(errMsg); err != nil {
		log.Printf("failed to send response :%q", err)
		return
	}
}

func (h *CustomerHandler) createCustomer(w http.ResponseWriter, r *http.Request) {
	var customer Customer
	if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
		handleResponseErr(w, http.StatusBadRequest, "invalid json body", err)
		return
	}

	if err := h.service.addCustomer(customer); err != nil {
		if errors.Is(err, ErrInvalidId) {
			handleResponseErr(w, http.StatusBadRequest, "invalid id", err)
			return
		}

		if errors.Is(err, ErrInvalidContactNo) {
			handleResponseErr(w, http.StatusBadRequest, "invalid contact number", err)
			return
		}

		if errors.Is(err, ErrConflict) {
			handleResponseErr(w, http.StatusConflict, "customer exists", err)
			return
		}
		handleResponseErr(w, http.StatusInternalServerError, "internal server error", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode("customer registered"); err != nil {
		log.Printf("failed to send response :%q", err)
		return
	}
}

func (h *CustomerHandler) updateCustomer(w http.ResponseWriter, r *http.Request) {
	var customer Customer

	if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
		handleResponseErr(w, http.StatusBadRequest, "invalid json body", err)
		return
	}

	if err := h.service.updateCustomer(customer); err != nil {
		if errors.Is(err, ErrInvalidId) {
			handleResponseErr(w, http.StatusBadRequest, "invalid id", err)
			return
		}

		if errors.Is(err, ErrInvalidContactNo) {
			handleResponseErr(w, http.StatusBadRequest, "invalid contact number", err)
			return
		}

		if errors.Is(err, ErrNotFound) {
			handleResponseErr(w, http.StatusNotFound, "customer not found", err)
			return
		}

		handleResponseErr(w, http.StatusInternalServerError, "internal server error", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode("customer details updated"); err != nil {
		log.Printf("failed to write response msg due to error :%q", err)
	}
}

func (h *CustomerHandler) getAllCustomer(w http.ResponseWriter, r *http.Request) {
	customers, err := h.service.getAllCustomer()
	if err != nil {
		handleResponseErr(w, http.StatusInternalServerError, "internal server error", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(customers); err != nil {
		log.Printf("failed to send response :%q", err)
		return
	}
}

func (h *CustomerHandler) getCustomerById(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	customer, err := h.service.getCustomerById(id)
	if err != nil {
		if errors.Is(err, ErrInvalidId) {
			handleResponseErr(w, http.StatusBadRequest, "invalid id", err)
			return
		}

		if errors.Is(err, ErrNotFound) {
			handleResponseErr(w, http.StatusNotFound, "customer not found", err)
			return
		}

		handleResponseErr(w, http.StatusInternalServerError, "internal server error", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(customer.CustomerDetails); err != nil {
		log.Printf("failed to write response due to error :%q", err)
	}
}

func (h *CustomerHandler) deleteCustomer(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if err := h.service.deleteCustomer(id); err != nil {
		if errors.Is(err, ErrInvalidId) {
			handleResponseErr(w, http.StatusBadRequest, "invalid id", err)
			return
		}

		if errors.Is(err, ErrNotFound) {
			handleResponseErr(w, http.StatusNotFound, "customer not found", err)
			return
		}

		handleResponseErr(w, http.StatusInternalServerError, "internal server error", err)
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode("customer deleted"); err != nil {
		log.Printf("failed to send response :%q", err)
		return
	}
}

// websocket implementation
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type websocketClient struct {
	clientId string
	client   *websocket.Conn
}

func NewWebsocketClient(clientId string, client *websocket.Conn) *websocketClient {
	return &websocketClient{
		clientId: clientId,
		client:   client,
	}
}

func (w *websocketClient) getSubscriberId() string {
	return w.clientId
}

func (w *websocketClient) update(customers []Customer) {
	customerList, err := json.Marshal(customers)
	if err != nil {
		log.Printf("encoding to json failed :%q\n", err)
		return
	}

	if err := w.client.WriteMessage(websocket.TextMessage, customerList); err != nil {
		log.Printf("failed to write message :%q", err)
		return
	}
}

func (h *CustomerHandler) websocketEndpoint(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("failed to upgrade request :%q", err)
		return
	}
	defer ws.Close()

	clientId := ws.RemoteAddr().String()
	client := NewWebsocketClient(clientId, ws)
	h.service.subscribe(client)
	defer h.service.unSubscribe(client)

	for {
		if _, _, err := ws.ReadMessage(); err != nil {
			log.Printf("failed to read message :%q\n", err)
			return
		}
	}
}

func registerRoutes(h *CustomerHandler) *mux.Router {
	router := mux.NewRouter()

	router.Methods("POST").Path("/api/customers").HandlerFunc(h.createCustomer)
	router.Methods("PUT").Path("/api/customers").HandlerFunc(h.updateCustomer)
	router.Methods("GET").Path("/api/customers/{id}").HandlerFunc(h.getCustomerById)
	router.Methods("GET").Path("/api/customers").HandlerFunc(h.getAllCustomer)
	router.Methods("DELETE").Path("/api/customers/{id}").HandlerFunc(h.deleteCustomer)
	router.HandleFunc("/ws", h.websocketEndpoint)

	return router
}
