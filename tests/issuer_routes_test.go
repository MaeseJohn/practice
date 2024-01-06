package tests

import (
	"API_Rest/models"
	"net/http"
	"testing"

	"gopkg.in/go-playground/assert.v1"
)

func TestCreateInvoice(t *testing.T) {
	CreateEnviroment()
	CreateValidUsers()
	var invoices = []*models.Invoice{
		models.NewInvoice("32143", "Invoice1", "2030-01-01", 500),
		models.NewInvoice("4a9d32f2-9b26-11ee-b9d1-0242ac120002", "Invoice1", "32141", 500),
		models.NewInvoice("4a9d32f2-9b26-11ee-b9d1-0242ac120002", "Invoice1", "2012-01-01", 500),
		models.NewInvoice("4a9d32f2-9b26-11ee-b9d1-0242ac120002", "Invoice1", "2150-01-01", 0),
		//Valid invoices
		models.NewInvoice("4a9d32f2-9b26-11ee-b9d1-0242ac120002", "Invoice1", "2150-01-01", 500),
	}
	var testData = []struct {
		name    string
		invoice *models.Invoice
		token   string
		status  int
	}{
		// Invalid Tokens
		{"Invalid token", invoices[0], invalidToken, http.StatusUnauthorized},
		{"Investor token", invoices[0], investorToken, http.StatusForbidden},
		{"Emty token", invoices[0], "", http.StatusUnauthorized},
		{"Unregister token", invoices[0], unregisterIssuerToken, http.StatusNotFound},
		// Invalid invoice uuid
		{"Invalid invoice uuid", invoices[0], issuerToken, http.StatusUnprocessableEntity},
		// Invalid Date
		{"Invalid date", invoices[1], issuerToken, http.StatusUnprocessableEntity},
		{"Invalid date", invoices[2], issuerToken, http.StatusUnprocessableEntity},
		// Invalid price
		{"Invalid price", invoices[3], issuerToken, http.StatusBadRequest},
		// Valid invoice
		{"Valid invoice", invoices[4], issuerToken, http.StatusOK},
		// Try to register same uuid invoice
		{"Duplicate key value uuid", invoices[4], issuerToken, http.StatusInternalServerError},
	}
	for _, i := range testData {
		t.Run(i.name, func(*testing.T) {
			resp := CreateInvoiceRequest(i.invoice, i.token)
			assert.Equal(t, resp, i.status)
		})
	}
}

func TestApproveInvoice(t *testing.T) {
	CreateEnviroment()
	CreateValidUsers()
	for _, i := range validInvoices {
		_ = CreateInvoiceRequest(&i, issuerToken)
	}
	_ = BuyInvoice(validInvoices[1].InvoiceId, investorToken, 500)

	var invoicesId = []struct {
		name      string
		invoiceId string
		status    int
	}{
		{"Unregister uuid", "e05b5242-a991-4841-9200-27d3157a506d", http.StatusNotFound},
		{"Emty uuid", "", http.StatusNotFound},
		{"Open invoice", validInvoices[0].InvoiceId, http.StatusBadRequest},
		{"Valid invoice", validInvoices[1].InvoiceId, http.StatusOK},
		{"Open invoice", validInvoices[2].InvoiceId, http.StatusBadRequest},
	}

	for _, i := range invoicesId {
		t.Run(i.name, func(t *testing.T) {
			resp := ApproveInvoceRequest(i.invoiceId, issuerToken)
			assert.Equal(t, resp, i.status)
		})
	}

}
