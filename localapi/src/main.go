package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go-huma-api-server/src/utils"
	"log"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/danielgtaylor/huma/v2/humacli"
	"github.com/go-chi/chi/v5"

	_ "github.com/danielgtaylor/huma/v2/formats/cbor"
	kivik "github.com/go-kivik/kivik/v4"
	_ "github.com/go-kivik/kivik/v4/couchdb" // The CouchDB driver
)

type Options struct {
	Port int `help:"Port to listen on" short:"p" default:"8080"`
}

type CertificateInput struct {
	Body struct {
		utils.Certificate
	}
}

type ReceviedCertificateResponse struct {
	Body struct {
		CertificateID  string `json:"certificate_id" doc:"Certificate ID"`
		PilotID        string `json:"pilot_id" doc:"Pilot ID"`
		DroneID        string `json:"drone_id" doc:"Drone ID"`
		ExpirationDate string `json:"expiration_date" doc:"Expiration date"`
		Message        string `json:"message" doc:"Message"`
	}
}

type CertificateListResponse struct {
	Body struct {
		Message      string              `json:"message" doc:"Message"`
		Certificates []utils.Certificate `json:"certificates"`
	} `json:"body"`
}

type Selector struct {
	Body struct {
		Selector map[string]interface{} `json:"selector" doc:"Selector to query certificates" default:"{\"pilot_id\": \"Pilot00\"}"`
	}
}

func main() {

	// Create a new CouchDB client
	db_client, err := kivik.New("couch", "http://admin:admin123@couchdb:5984")
	// db_client, err := kivik.New("couch", "http://localhost:5984/")
	if err != nil {
		panic(err)
	}
	log.Printf("Connected to CouchDB: %s", db_client.DSN())

	// Create the database if it doesn't exist

	exists, err := db_client.DBExists(context.Background(), "certificates")
	if err != nil {
		panic(err)
	}
	if !exists {
		// Create the database
		if err := db_client.CreateDB(context.Background(), "certificates"); err != nil {
			panic(err)
		}
		log.Println("Created database: certificates")
	} else {
		log.Println("Database already exists: certificates")
	}

	db := db_client.DB("certificates")

	// Create a CLI app which takes a port option.
	cli := humacli.New(func(hooks humacli.Hooks, options *Options) {
		// Create a new router & API
		router := chi.NewMux()
		api := humachi.New(router, huma.DefaultConfig("My API", "1.0.0"))

		// Register POST /certificates/create
		huma.Register(api, huma.Operation{
			OperationID:   "post-certificate",
			Method:        http.MethodPost,
			Path:          "/certificates/create",
			Summary:       "Post a certificate",
			Tags:          []string{"Certificates"},
			Description:   "Post a certificate to the server.",
			DefaultStatus: http.StatusCreated,
		}, func(ctx context.Context, input *CertificateInput) (*ReceviedCertificateResponse, error) {

			log.Printf("Received certificate with PilotID: %v", input.Body.PilotID)
			cert_db_obj, err := utils.GetCertificateDBObject(&input.Body.Certificate)
			if err != nil {
				return nil, err
			}
			log.Printf("Received certificate: %s", cert_db_obj.CertificateID)

			// Put it in the database
			utils.PutCertificateToDB(db, cert_db_obj)

			// Generate the response
			resp := &ReceviedCertificateResponse{}
			resp.Body.CertificateID = cert_db_obj.CertificateID
			resp.Body.PilotID = cert_db_obj.PilotID
			resp.Body.DroneID = cert_db_obj.DroneID
			resp.Body.ExpirationDate = cert_db_obj.ExpirationDate
			resp.Body.Message = "Certificate received and stored successfully."
			return resp, nil
		})

		// Register POST /certificates/query
		huma.Register(api, huma.Operation{
			OperationID:   "post-certificate-query",
			Method:        http.MethodPost,
			Path:          "/certificates/query",
			Summary:       "Query certificates with selector",
			Tags:          []string{"Certificates"},
			Description:   "Query certificates with a selector.",
			DefaultStatus: http.StatusOK,
		}, func(ctx context.Context, input *Selector) (*CertificateListResponse, error) {
			log.Printf("Received selector: %v", input.Body.Selector)
			certificates, err := utils.GetCertificatesWithSelector(db, input.Body.Selector)
			if err != nil {
				return nil, err
			}
			log.Printf("Found %d certificates", len(certificates))
			var certs []utils.Certificate
			for _, cert := range certificates {
				var cert_obj utils.Certificate
				err := json.Unmarshal([]byte(cert.SerializedCertificate), &cert_obj)
				if err != nil {
					return nil, err
				}
				// fmt.Printf("Certificate: %v", cert_obj)
				certs = append(certs, cert_obj)
			}
			resp := &CertificateListResponse{}
			resp.Body.Message = fmt.Sprintf("Found %d certificates", len(certs))
			resp.Body.Certificates = certs
			return resp, nil
		})

		// Tell the CLI how to start your server.
		hooks.OnStart(func() {
			fmt.Printf("Starting server on port %d...\n", options.Port)
			http.ListenAndServe(fmt.Sprintf(":%d", options.Port), router)
		})
	})

	// Run the CLI. When passed no commands, it starts the server.
	cli.Run()
}
