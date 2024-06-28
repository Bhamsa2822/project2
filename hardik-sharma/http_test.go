package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestCustomerHandler_createCustomer(t *testing.T) {
	type fields struct {
		customers []Customer
	}

	tests := []struct {
		name       string
		fields     fields
		reqbody    string
		wantStatus int
		wantBody   string
	}{
		{
			name: "invalid id",
			fields: fields{
				customers: []Customer{},
			},
			reqbody: `
			{
				"id": "h",
				"customerDetails": {
					"name": "hdik",
					"address": "hsghd",
					"contactNo": 9649127584
				}
			}
			`,
			wantStatus: http.StatusBadRequest,
			wantBody:   `"invalid id"`,
		},
		{
			name: "invalid json body",
			fields: fields{
				customers: []Customer{},
			},
			reqbody: `
			{
				"id": "h,
				"customerDetails": {
					"name": "hdik",
					"address": "hsghd",
					"contactNo": 9649127584
				}
			}
			`,
			wantStatus: http.StatusBadRequest,
			wantBody:   `"invalid json body"`,
		},
		{
			name: "existing customer",
			fields: fields{
				customers: []Customer{
					{
						Id: "hs",
						CustomerDetails: CustomerDetails{
							Name:      "hardik",
							Address:   "uadipur",
							ContactNo: 9999999999,
						},
					},
				},
			},
			reqbody: `
			{
				"id": "hs",
				"customerDetails": {
					"name": "hdik",
					"address": "hsghd",
					"contactNo": 9649127584
				}
			}
			`,
			wantStatus: http.StatusConflict,
			wantBody:   `"customer exists"`,
		},
		{
			name: "invalid contact number",
			fields: fields{
				customers: []Customer{
					{
						Id: "hs",
						CustomerDetails: CustomerDetails{
							Name:      "hardik",
							Address:   "uadipur",
							ContactNo: 9999999999,
						},
					},
				},
			},
			reqbody: `
			{
				"id": "vs",
				"customerDetails": {
					"name": "varshil",
					"address": "udr",
					"contactNo": 96491275
				}
			}
			`,
			wantStatus: http.StatusBadRequest,
			wantBody:   `"invalid contact number"`,
		},
		{
			name: "non existing valid customer",
			fields: fields{
				customers: []Customer{
					{
						Id: "hs",
						CustomerDetails: CustomerDetails{
							Name:      "hardik",
							Address:   "uadipur",
							ContactNo: 9999999999,
						},
					},
				},
			},
			reqbody: `
			{
				"id": "vs",
				"customerDetails": {
					"name": "varshil",
					"address": "udr",
					"contactNo": 8888888888
				}
			}
			`,
			wantStatus: http.StatusCreated,
			wantBody:   `"customer registered"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewInMemoryRepo()
			repo.customers = tt.fields.customers
			service := NewService(repo)
			transport := NewCustomerHandler(service)

			body := strings.NewReader(tt.reqbody)
			r := httptest.NewRequest("POST", "/api/customers", body)
			w := httptest.NewRecorder()

			transport.createCustomer(w, r)

			assert.JSONEq(t, tt.wantBody, w.Body.String(), "expect body to be same")

			if w.Code != tt.wantStatus {
				t.Errorf("want status :%d got status %d", tt.wantStatus, w.Code)
			}
		})
	}
}

func TestCustomerHandler_updateCustomer(t *testing.T) {
	type fields struct {
		customers []Customer
	}

	tests := []struct {
		name     string
		fields   fields
		reqBody  string
		wantBody string
		wantCode int
	}{
		{
			name: "invalid json body",
			fields: fields{
				customers: []Customer{},
			},
			reqBody: `{
			"id":"hs ,
			"customerDetails":
			{
				"name":"hdik",
				"address":"hsghd",
				"contactNo":9649127059
			}
		}`,
			wantBody: `"invalid json body"`,
			wantCode: http.StatusBadRequest,
		},
		{
			name: "invalid id",
			fields: fields{
				customers: []Customer{},
			},
			reqBody: `{
			"id":"hsss",
			"customerDetails":
			{
				"name":"hdik",
				"address":"hsghd",
				"contactNo":9649127059
			}
		}`,
			wantBody: `"invalid id"`,
			wantCode: http.StatusBadRequest,
		},
		{
			name: "invalid contact number",
			fields: fields{
				customers: []Customer{},
			},
			reqBody: `{
			"id":"hs",
			"customerDetails":
			{
				"name":"hdik",
				"address":"hsghd",
				"contactNo":96491270
			}
		}`,
			wantBody: `"invalid contact number"`,
			wantCode: http.StatusBadRequest,
		},
		{
			name: "updating existing customer",
			fields: fields{
				customers: []Customer{
					{
						Id: "hs",
						CustomerDetails: CustomerDetails{
							Name:      "hardik",
							Address:   "udaipur",
							ContactNo: 7777777777,
						},
					},
					{
						Id: "vs",
						CustomerDetails: CustomerDetails{
							Name:      "vastghj",
							Address:   "udaipur",
							ContactNo: 9876655547,
						},
					},
				},
			},
			reqBody: `{
			"id":"hs",
			"customerDetails":
			{
				"name":"hdik",
				"address":"hsghd",
				"contactNo":9649127559
			}
		}`,
			wantBody: `"customer details updated"`,
			wantCode: http.StatusOK,
		},
		{
			name: "updating non existing customer",
			fields: fields{
				customers: []Customer{
					{
						Id: "hs",
						CustomerDetails: CustomerDetails{
							Name:      "hardik",
							Address:   "udaipur",
							ContactNo: 7777777777,
						},
					},
					{
						Id: "vs",
						CustomerDetails: CustomerDetails{
							Name:      "vastghj",
							Address:   "udaipur",
							ContactNo: 9876655547,
						},
					},
				},
			},
			reqBody: `{
			"id":"ps",
			"customerDetails":
			{
				"name":"ps",
				"address":"ps",
				"contactNo":9649127559
			}
		}`,
			wantBody: `"customer not found"`,
			wantCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &InMemoryRepo{tt.fields.customers}
			service := NewService(repo)
			transport := NewCustomerHandler(service)
			handle := registerRoutes(transport)

			body := strings.NewReader(tt.reqBody)

			r := httptest.NewRequest("PUT", "/api/customers", body)
			w := httptest.NewRecorder()

			handle.ServeHTTP(w, r)

			assert.JSONEq(t, tt.wantBody, w.Body.String(), "expect body to be same")

			assert.Equal(t, tt.wantCode, w.Code, "expect status code to be same")
		})
	}
}

func TestCustomerHandler_getAllCustomer(t *testing.T) {
	type fields struct {
		customers []Customer
	}

	tests := []struct {
		name     string
		fields   fields
		wantBody string
		wantCode int
	}{
		{
			name: "existing customers",
			fields: fields{
				customers: []Customer{
					{
						Id: "hs",
						CustomerDetails: CustomerDetails{
							Name:      "hardik",
							Address:   "udaipur",
							ContactNo: 7777777777,
						},
					},
					{
						Id: "vs",
						CustomerDetails: CustomerDetails{
							Name:      "varshil",
							Address:   "udaipur",
							ContactNo: 6666666666,
						},
					},
					{
						Id: "ps",
						CustomerDetails: CustomerDetails{
							Name:      "paramveer",
							Address:   "udaipur",
							ContactNo: 5555555555,
						},
					},
				},
			},
			wantBody: `[
				{
					"id": "hs",
					"customerDetails": {
						"name": "hardik",
						"address": "udaipur",
						"contactNo": 7777777777
					}
				},
				{
					"id": "vs",
					"customerDetails": {
						"name": "varshil",
						"address": "udaipur",
						"contactNo": 6666666666
					}
				},
				{
					"id": "ps",
					"customerDetails": {
						"name": "paramveer",
						"address": "udaipur",
						"contactNo": 5555555555
					}
				}
			]`,
			wantCode: http.StatusOK,
		},
		{
			name: "empty database",
			fields: fields{
				customers: []Customer{},
			},
			wantBody: `[]`,
			wantCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &InMemoryRepo{customers: tt.fields.customers}
			service := NewService(repo)
			transport := NewCustomerHandler(service)
			handler := registerRoutes(transport)

			w := httptest.NewRecorder()

			handler.ServeHTTP(w, httptest.NewRequest("GET", "/api/customers", nil))

			assert.JSONEq(t, tt.wantBody, w.Body.String(), "expect body to be same")

			assert.Equal(t, tt.wantCode, w.Code, "expect status code to be same")
		})
	}
}

func TestCustomerHandler_getCustomerById(t *testing.T) {
	type fields struct {
		customers []Customer
	}

	tests := []struct {
		name     string
		fields   fields
		path     string
		wantBody string
		wantCode int
	}{
		{
			name: "invalid id",
			fields: fields{
				customers: []Customer{},
			},
			path:     "/api/customers/hss",
			wantBody: `"invalid id"`,
			wantCode: http.StatusBadRequest,
		},
		{
			name: "non existing customer",
			fields: fields{
				customers: []Customer{
					{
						Id: "ps",
						CustomerDetails: CustomerDetails{
							Name:      "hshsj",
							Address:   "udr",
							ContactNo: 9999999999,
						},
					},
					{
						Id: "vs",
						CustomerDetails: CustomerDetails{
							Name:      "hshsnjmj",
							Address:   "jaiop",
							ContactNo: 9999999999,
						},
					},
				},
			},
			path:     "/api/customers/hs",
			wantBody: `"customer not found"`,
			wantCode: http.StatusNotFound,
		},
		{
			name: "non existing customer",
			fields: fields{
				customers: []Customer{
					{
						Id: "ps",
						CustomerDetails: CustomerDetails{
							Name:      "hshsj",
							Address:   "udr",
							ContactNo: 9999999999,
						},
					},
					{
						Id: "vs",
						CustomerDetails: CustomerDetails{
							Name:      "vashil",
							Address:   "udaipur",
							ContactNo: 9999999999,
						},
					},
				},
			},
			path: "/api/customers/vs",
			wantBody: `
			{
				"name": "vashil",
				"address": "udaipur",
				"contactNo": 9999999999
			}
			`,
			wantCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &InMemoryRepo{tt.fields.customers}
			service := NewService(repo)
			transport := NewCustomerHandler(service)
			handler := registerRoutes(transport)

			w := httptest.NewRecorder()

			handler.ServeHTTP(w, httptest.NewRequest("GET", tt.path, nil))

			assert.JSONEq(t, tt.wantBody, w.Body.String(), "expect body to be same")

			assert.Equal(t, tt.wantCode, w.Code, "want status code to be same")
		})
	}
}

func TestCustomerHandler_deleteCustomer(t *testing.T) {
	type fields struct {
		customers []Customer
	}

	tests := []struct {
		name     string
		fields   fields
		path     string
		wantBody string
		wantCode int
	}{
		{
			name: "invalid id",
			fields: fields{
				customers: []Customer{},
			},
			path:     "/api/customers/hsss",
			wantBody: `"invalid id"`,
			wantCode: http.StatusBadRequest,
		},
		{
			name: "non existing customer",
			fields: fields{
				customers: []Customer{
					{
						Id: "hs",
						CustomerDetails: CustomerDetails{
							Name:      "hardik",
							Address:   "uadipur",
							ContactNo: 9999999999,
						},
					},
				},
			},
			path:     "/api/customers/js",
			wantBody: `"customer not found"`,
			wantCode: http.StatusNotFound,
		},
		{
			name: "existing customer",
			fields: fields{
				customers: []Customer{
					{
						Id: "hs",
						CustomerDetails: CustomerDetails{
							Name:      "hardik",
							Address:   "uadipur",
							ContactNo: 9999999999,
						},
					},
					{
						Id: "js",
						CustomerDetails: CustomerDetails{
							Name:      "jdskik",
							Address:   "uadopksipur",
							ContactNo: 9888889999,
						},
					},
				},
			},
			path:     "/api/customers/js",
			wantBody: `"customer deleted"`,
			wantCode: http.StatusOK,
		},
		{
			name: "empty databasde",
			fields: fields{
				customers: []Customer{},
			},
			path:     "/api/customers/hs",
			wantBody: `"customer not found"`,
			wantCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &InMemoryRepo{tt.fields.customers}
			service := NewService(repo)
			transport := NewCustomerHandler(service)
			handler := registerRoutes(transport)

			w := httptest.NewRecorder()

			handler.ServeHTTP(w, httptest.NewRequest("DELETE", tt.path, nil))

			assert.JSONEq(t, tt.wantBody, w.Body.String(), "expect body to be same")

			assert.Equal(t, tt.wantCode, w.Code, "expect status code to be same")
		})
	}
}

func TestCustomerHandler_WSCreateCustomer(t *testing.T) {
	repo := &InMemoryRepo{
		customers: []Customer{
			{
				Id: "hs",
				CustomerDetails: CustomerDetails{
					Name:      "hardik",
					Address:   "udr",
					ContactNo: 8888888888,
				},
			},
		},
	}
	service := NewService(repo)
	transport := NewCustomerHandler(service)
	handler := registerRoutes(transport)

	// creating backend server
	go func() {
		if err := http.ListenAndServe(":7070", handler); err != nil {
			t.Logf("connection failed with server: %v", err)
		}
	}()

	//establishing websocket connection
	conn, wsRes, err := websocket.DefaultDialer.Dial("ws://localhost:7070/ws", nil)
	if err != nil {
		t.Fatalf("failed to establish websocket connection: %v", err)
	}

	assert.Equal(t, http.StatusSwitchingProtocols, wsRes.StatusCode, "expected status code to be same")

	//making http req
	body := strings.NewReader(`
   {
	   "id": "vs",
	   "customerDetails": {
		   "name": "varshil",
		   "address": "udr",
		   "contactNo": 8888888888
	   }
   }
   `)

	resp, err := http.Post("http://localhost:7070/api/customers", "application/json", body)
	if err != nil {
		t.Fatalf("http request failed :%v", err)
	}

	assert.Equal(t, http.StatusCreated, resp.StatusCode, "expected status code to be same")

	mt, message, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("failed to read messages :%v", err)
	}

	assert.Equal(t, websocket.TextMessage, mt, "expected message type to be same")

	wantNotification := `[
		{
			"id": "hs",
			"customerDetails": {
				"name": "hardik",
				"address": "udr",
				"contactNo": 8888888888
			}
		},
		{
			"id": "vs",
			"customerDetails": {
				"name": "varshil",
				"address": "udr",
				"contactNo": 8888888888
			}		
		}
		]`

	assert.JSONEq(t, wantNotification, string(message), "expecting customer list to be same")
}

func TestCustomerHandler_WSUpdateCustomer(t *testing.T) {
	repo := &InMemoryRepo{
		customers: []Customer{
			{
				Id: "vs",
				CustomerDetails: CustomerDetails{
					Name:      "hardik",
					Address:   "udr",
					ContactNo: 8888888888,
				},
			},
			{
				Id: "hs",
				CustomerDetails: CustomerDetails{
					Name:      "hd",
					Address:   "udr",
					ContactNo: 8888888888,
				},
			},
		},
	}
	service := NewService(repo)
	transport := NewCustomerHandler(service)
	handler := registerRoutes(transport)

	go func() {
		if err := http.ListenAndServe(":6060", handler); err != nil {
			log.Fatalf("server connection failed :%v", err)
		}
	}()

	client := &http.Client{}

	conn, wsRes, err := websocket.DefaultDialer.Dial("ws://localhost:6060/ws", nil)
	if err != nil {
		t.Fatalf("failed to establish websocket connection :%v", err)
	}

	assert.Equal(t, http.StatusSwitchingProtocols, wsRes.StatusCode, "expected status code to be same")

	body := strings.NewReader(`
   {
	   "id": "vs",
	   "customerDetails": {
		   "name": "varshil",
		   "address": "udr",
		   "contactNo": 8888888888
	   }
   }
   `)

	req, err := http.NewRequest("PUT", "http://localhost:6060/api/customers", body)
	if err != nil {
		t.Fatalf("http request failed :%v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("request failed :%v", err)
	}

	assert.Equal(t, http.StatusOK, resp.StatusCode, "expected status code to be same")

	mt, message, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("failed to read message :%v", err)
	}

	assert.Equal(t, websocket.TextMessage, mt, "expected message type to be same")

	wantCustomerList := `[
		{
			"id": "vs",
			"customerDetails": {
				"name": "varshil",
				"address": "udr",
				"contactNo": 8888888888
			}
		},
		{
			"id": "hs",
			"customerDetails": {
				"name": "hd",
				"address": "udr",
				"contactNo": 8888888888
			}		
		}
		]`

	assert.JSONEq(t, wantCustomerList, string(message), "expected customer list to be same")
}

func TestCustomerHandler_WSDeleteCustomer(t *testing.T) {
	repo := &InMemoryRepo{
		customers: []Customer{
			{
				Id: "vs",
				CustomerDetails: CustomerDetails{
					Name:      "varshil",
					Address:   "udr",
					ContactNo: 8888888888,
				},
			},
			{
				Id: "hs",
				CustomerDetails: CustomerDetails{
					Name:      "hardik",
					Address:   "udr",
					ContactNo: 8888888888,
				},
			},
		},
	}
	service := NewService(repo)
	transport := NewCustomerHandler(service)
	handler := registerRoutes(transport)

	go func() {
		if err := http.ListenAndServe(":5050", handler); err != nil {
			t.Logf("server connection failed :%v", err)
		}
	}()

	client := http.Client{}

	conn, wsRes, err := websocket.DefaultDialer.Dial("ws://localhost:5050/ws", nil)
	if err != nil {
		t.Fatalf("failed to establish websocket connection :%v", err)
	}

	assert.Equal(t, http.StatusSwitchingProtocols, wsRes.StatusCode, "expected status code to be same")

	req, err := http.NewRequest("DELETE", "http://localhost:5050/api/customers/vs", nil)
	if err != nil {
		t.Fatalf("http request failed :%v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("request failed :%v", err)
	}

	assert.Equal(t, http.StatusOK, resp.StatusCode, "expected status code to be same")

	mt, message, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("failed to read message :%v", err)
	}

	assert.Equal(t, websocket.TextMessage, mt, "expected message type to be same")

	wantCustomerList := `[
		    {
				"id": "hs",
				"customerDetails": {
					"name": "hardik",
					"address": "udr",
					"contactNo":8888888888
				}
			}
		]`

	assert.JSONEq(t, wantCustomerList, string(message), "expected customer list to be same")
}
