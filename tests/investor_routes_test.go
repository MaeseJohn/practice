package tests

import (
	"API_Rest/models"
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"net/http"
	"testing"

	"github.com/tjgq/broadcast"
	"gopkg.in/go-playground/assert.v1"
)

func TestBuyInvoice(t *testing.T) {
	CreateEnviroment()
	CreateValidUsers()
	_ = CreateInvoiceRequest(&validInvoices[0], issuerToken)

	//Nombre error t.Run(nombre funcion)
	//http constantes
	var buyParameters = []struct {
		name      string
		invoiceId string
		funds     int
		token     string
		status    int
	}{

		// Invalid tokens
		{"Zero value token", "4a9d32f2-9b26-11ee-b9d1-0242ac120002", 101, "", http.StatusUnauthorized},
		{"Invalid token", "4a9d32f2-9b26-11ee-b9d1-0242ac120002", 102, invalidToken, http.StatusUnauthorized},
		{"Issuer token", "4a9d32f2-9b26-11ee-b9d1-0242ac120002", 100, issuerToken, http.StatusForbidden},
		{"Unresgister Usertoken", "4a9d32f2-9b26-11ee-b9d1-0242ac120002", 100, unregisterInvestorToken, http.StatusNotFound},
		// Invalid uuid
		{"Not uuid", "hola", 1000, investorToken, http.StatusUnprocessableEntity},
		{"Unregister uuid", "4ddbb37b-efd5-4564-a2ba-c4ac80925b9f", 100, investorToken, http.StatusNotFound},
		// To many funds
		{"To many funds", "4a9d32f2-9b26-11ee-b9d1-0242ac120002", 1000, investorToken, http.StatusForbidden},
		// Valid purcharse
		{"Valid purcharse", "4a9d32f2-9b26-11ee-b9d1-0242ac120002", 100, investorToken, http.StatusOK},
		{"Valid purcharse", "4a9d32f2-9b26-11ee-b9d1-0242ac120002", 400, investorToken, http.StatusOK},
		// Closed invoice
		{"Closed invoice", "4a9d32f2-9b26-11ee-b9d1-0242ac120002", 100, investorToken, http.StatusNotFound},
	}

	for _, p := range buyParameters {
		t.Run(p.name, func(t *testing.T) {
			resp := BuyInvoice(p.invoiceId, p.token, p.funds)
			assert.Equal(t, resp, p.status)
		})
	}
}

func TestGetInvoices(t *testing.T) {
	CreateEnviroment()
	CreateValidUsers()
	for _, i := range validInvoices {
		_ = CreateInvoiceRequest(&i, issuerToken)
	}
	_ = BuyInvoice("8ecc1ee4-347d-4fc3-b32c-b2cd52c58f16", investorToken, 500)

	// Create a new request using http
	req, err := http.NewRequest("GET", "http://localhost:3000/invoices", nil)
	if err != nil {
		log.Fatalln("Error on request.\n[ERROR] -", err)
	}

	// Create a Bearer string by appending string access token
	var bearer = "Bearer " + investorToken

	// Add authorization header to the req
	req.Header.Add("Authorization", bearer)

	// Send req using http Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln("Error on response.\n[ERROR] -", err)
	}
	defer resp.Body.Close()

	var invoices []models.Invoice

	//Read body
	if resp.Body != nil {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln("Error while reading the response bytes:", err)
		}
		//Unmarshal body to user
		err = json.Unmarshal(body, &invoices)
		if err != nil {
			log.Fatalln(err)
		}
	}

	validInvoices[0].Status = "open"
	validInvoices[2].Status = "open"
	invoices[0].ExpireDate = validInvoices[0].ExpireDate
	invoices[1].ExpireDate = validInvoices[2].ExpireDate

	assert.Equal(t, resp.StatusCode, http.StatusOK)
	assert.Equal(t, invoices[0], validInvoices[0])
	assert.Equal(t, invoices[1], validInvoices[2])
}

func TestConcurrentInvoicePurchase(t *testing.T) {
	CreateEnviroment()
	var randomUsers = []*models.User{
		models.NewUser("89a08faa-df68-4df9-8499-4a9497dd43f3", "miguel@gmail.com", "miguel", "miguel", "investor", 1000),
		models.NewUser("7783d257-4d0b-470e-80b1-4d54f022ad2d", "jose@gmail.com", "jose", "jose", "investor", 100),
		models.NewUser("e9473b04-459b-4b7f-871c-d71f84d0b163", "pepe@gmail.com", "pepe", "pepe", "investor", 1000),
		models.NewUser("bbdf30ca-2161-4dc1-9ece-6fca26de7287", "manuel@gmail.com", "manuel", "manuel", "investor", 500),
		models.NewUser("0fb482c7-420d-49ec-afa9-b68dd322b664", "max@gmail.com", "max", "max", "investor", 100),
		models.NewUser("89e60bfc-9b26-4e5e-9c6e-d0775c2a6cc1", "josian@gmail.com", "josian", "josian", "investor", 1000),
		models.NewUser("5ea8f8b9-f38e-45ed-a736-42b960c76534", "rosa@gmail.com", "rosa", "rosa", "investor", 1000),
		models.NewUser("808fa41b-1460-4250-adc0-12b23b8633b0", "darius@gmail.com", "darius", "darius", "investor", 300),
		models.NewUser("ce646384-4c78-44ff-a085-8984e8ab54aa", "gareen@gmail.com", "gareen", "gareen", "investor", 2000),
		models.NewUser("4199efa2-a5e3-4f2e-8470-e2c66a45484b", "jarvan@gmail.com", "jarvan", "jarvan", "investor", 400),
		models.NewUser("93d685f5-b485-4156-bf46-9209fd3ad92f", "dimitry@gmail.com", "dimitry", "dimitry", "investor", 300),
		models.NewUser("aa67c643-7fda-4d07-b203-c960ef4d8a18", "vlad@gmail.com", "vlad", "vlad", "investor", 250),
		models.NewUser("cafe0a1d-c079-48d5-af0d-a840a135c099", "morgan@gmail.com", "morgan", "morgan", "investor", 6000),
		models.NewUser("5526a2f3-17f2-4379-985f-698ae6abc350", "olivia@gmail.com", "olivia", "olivia", "investor", 7000),
		models.NewUser("e065ff2d-e488-4a4a-bb77-2ef93733cb1b", "alicia@gmail.com", "alicia", "alicia", "investor", 0),
		//randomUsers[15] == issuers
		models.NewUser("89a08faa-df68-4df9-8499-4a9497dd43f3", "jorge@gmail.com", "jorge", "jorge", "issuer", 0),
		models.NewUser("7783d257-4d0b-470e-80b1-4d54f022ad2d", "vicente@gmail.com", "vicente", "vicente", "issuer", 0),
		models.NewUser("e9473b04-459b-4b7f-871c-d71f84d0b163", "marcos@gmail.com", "marcos", "marcos", "issuer", 0),
		models.NewUser("bbdf30ca-2161-4dc1-9ece-6fca26de7287", "sara@gmail.com", "sara", "sara", "issuer", 0),
		models.NewUser("0fb482c7-420d-49ec-afa9-b68dd322b664", "emilia@gmail.com", "emilia", "emilia", "issuer", 0),
		models.NewUser("89e60bfc-9b26-4e5e-9c6e-d0775c2a6cc1", "pocoyo@gmail.com", "pocoyo", "pocoyo", "issuer", 0),
	}

	var randomInvoices = []*models.Invoice{
		models.NewInvoice("3ffcbee5-60db-4ef8-89ad-15e0c8b4db21", "invoicename", "01-01-2050", rand.Intn(5000)),
		models.NewInvoice("cfcc9434-09f5-422e-b4cf-2e437e612149", "invoicename", "01-01-2050", rand.Intn(5000)),
		models.NewInvoice("1a7e3e22-c07e-4e38-b3a4-3af72ccb430d", "invoicename", "01-01-2050", rand.Intn(5000)),
		models.NewInvoice("70763a21-a0ec-4aab-910a-a5e81f6f8d41", "invoicename", "01-01-2050", rand.Intn(5000)),
		models.NewInvoice("d933091f-272c-42de-9dfb-77ec41a7d986", "invoicename", "01-01-2050", rand.Intn(5000)),
		models.NewInvoice("e4f1aa17-e4b8-451a-83a1-f96a568081bf", "invoicename", "01-01-2050", rand.Intn(5000)),
		models.NewInvoice("541a3623-00ee-4a10-af81-8e6fa0bd7232", "invoicename", "01-01-2050", rand.Intn(5000)),
		models.NewInvoice("eba4926b-d40f-4c16-bd2a-ec668c1a0165", "invoicename", "01-01-2050", rand.Intn(5000)),
		models.NewInvoice("57fbe9d8-2493-4c88-8c69-d1b8c91de792", "invoicename", "01-01-2050", rand.Intn(5000)),
		models.NewInvoice("2633e440-d8f9-4299-a9c1-35abe28c1813", "invoicename", "01-01-2050", rand.Intn(5000)),
		models.NewInvoice("978ca4c7-2e6c-4c11-8ec1-d36e854ff74a", "invoicename", "01-01-2050", rand.Intn(5000)),
		models.NewInvoice("bd1a8de9-109b-4e78-8875-eb66be200166", "invoicename", "01-01-2050", rand.Intn(5000)),
		models.NewInvoice("9e12d48e-3ebe-41b4-ae46-2014b7a38e12", "invoicename", "01-01-2050", rand.Intn(5000)),
		models.NewInvoice("fd1672dd-883c-4181-a6e7-ec61c39481c0", "invoicename", "01-01-2050", rand.Intn(5000)),
		models.NewInvoice("63526b29-84ab-4f54-87eb-8779b33d8c4a", "invoicename", "01-01-2050", rand.Intn(5000)),
		models.NewInvoice("93d5555d-dad8-4c2d-ad02-f9f48550838a", "invoicename", "01-01-2050", rand.Intn(5000)),
	}

	// CREATE USERS AND COUNT TOTAL FUNDS
	var totalInvestorsFunds int
	var totalIssuersFunds int
	for _, u := range randomUsers {
		_ = CreateUserRequest(u)
		if u.Role == "investor" {
			totalInvestorsFunds += u.Funds
		}
		if u.Role == "issuer" {
			totalIssuersFunds += u.Funds
		}
	}

	// CREATING INVOICES WITH RANDOM ISSUERS
	for _, i := range randomInvoices {
		randomIssuer := randomUsers[rand.Intn(6)+15]
		_, tokenString := LoginUserRequest(randomIssuer.Email, randomIssuer.Password)
		CreateInvoiceRequest(i, tokenString)
	}

	var b broadcast.Broadcaster
	for i := 0; i < 100; i++ {
		l := b.Listen()
		invoice := randomInvoices[rand.Intn(16)]
		investor := randomUsers[rand.Intn(15)]
		_, tokenString := LoginUserRequest(investor.Email, investor.Password)
		invested := (rand.Intn(5) + 1) * 100

		go func() {
			for _ = range l.Ch {
				BuyInvoice(invoice.InvoiceId, tokenString, invested)
			}
		}()
	}
	b.Send("OPEN")
	b.Close()

}
