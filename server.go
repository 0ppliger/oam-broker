package main

import (
	"encoding/json"
	"net/http"
	"reflect"
	"context"
	"github.com/owasp-amass/asset-db/repository/neo4j"
	oam "github.com/owasp-amass/open-asset-model"
	oam_dns "github.com/owasp-amass/open-asset-model/dns"
	oam_net "github.com/owasp-amass/open-asset-model/network"
	oam_org "github.com/owasp-amass/open-asset-model/org"
	oam_url "github.com/owasp-amass/open-asset-model/url"
	oam_cert "github.com/owasp-amass/open-asset-model/certificate"
	oam_pf "github.com/owasp-amass/open-asset-model/platform"
	oam_contact "github.com/owasp-amass/open-asset-model/contact"
	oam_file "github.com/owasp-amass/open-asset-model/file"
	oam_financial "github.com/owasp-amass/open-asset-model/financial"
	oam_general "github.com/owasp-amass/open-asset-model/general"
	oam_people "github.com/owasp-amass/open-asset-model/people"
	oam_reg "github.com/owasp-amass/open-asset-model/registration"
	oam_account "github.com/owasp-amass/open-asset-model/account"
)


type AssetInput struct {
	Type oam.AssetType `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

var assetTypes = map[oam.AssetType]reflect.Type{
	oam.Account          : reflect.TypeOf(oam_account.Account{}),
	oam.AutnumRecord     : reflect.TypeOf(oam_reg.AutnumRecord{}),
	oam.AutonomousSystem : reflect.TypeOf(oam_net.AutonomousSystem{}),
	oam.ContactRecord    : reflect.TypeOf(oam_contact.ContactRecord{}),
	oam.DomainRecord     : reflect.TypeOf(oam_reg.DomainRecord{}),
	oam.File             : reflect.TypeOf(oam_file.File{}),
	oam.FQDN             : reflect.TypeOf(oam_dns.FQDN{}),
	oam.FundsTransfer    : reflect.TypeOf(oam_financial.FundsTransfer{}),
	oam.Identifier       : reflect.TypeOf(oam_general.Identifier{}),
	oam.IPAddress        : reflect.TypeOf(oam_net.IPAddress{}),
	oam.IPNetRecord      : reflect.TypeOf(oam_reg.IPNetRecord{}),
	oam.Location         : reflect.TypeOf(oam_contact.Location{}),
	oam.Netblock         : reflect.TypeOf(oam_net.Netblock{}),
	oam.Organization     : reflect.TypeOf(oam_org.Organization{}),
	oam.Person           : reflect.TypeOf(oam_people.Person{}),
	oam.Phone            : reflect.TypeOf(oam_contact.Phone{}),
	oam.Product          : reflect.TypeOf(oam_pf.Product{}),
	oam.ProductRelease   : reflect.TypeOf(oam_pf.ProductRelease{}),
	oam.Service          : reflect.TypeOf(oam_pf.Service{}),
	oam.TLSCertificate   : reflect.TypeOf(oam_cert.TLSCertificate{}),
	oam.URL              : reflect.TypeOf(oam_url.URL{}),
}

type assetHandler struct {}

func (h *assetHandler) Hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world"))
}
func (h *assetHandler) upsertAsset(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		http.Error(w, "no body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var input AssetInput
	
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&input); err != nil {
		http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	T, ok := assetTypes[input.Type]
	if !ok {
		http.Error(w, "unsupported asset type", http.StatusBadRequest)
		return
	}

	neo_store, err := neo4j.New(neo4j.Neo4j, "bolt://neo4j:password@localhost:7687/neo4j")
	if err != nil {
		http.Error(w, "Unable to connect to db: "+err.Error(), http.StatusBadRequest)
		return		
	}

	asset := reflect.New(T)

	if err := json.Unmarshal(input.Payload, asset.Interface()); err != nil {
		http.Error(w, "invalid FQDN: "+err.Error(), http.StatusBadRequest)
		return
	}
	
	ctx := context.Background()
	entity, err := neo_store.CreateAsset(ctx, asset.Interface().(oam.Asset))
	if err != nil {
		http.Error(w, "Failed to create asset: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte(entity.ID))
	
}

func (h *assetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet:
		h.Hello(w, r)
		return
	case r.Method == http.MethodPost:
		h.upsertAsset(w, r)
		return
	default:
		return
	}
}

func main() {
	mux := http.NewServeMux()

	mux.Handle("/asset", &assetHandler{})
	mux.Handle("/asset/", &assetHandler{})

	http.ListenAndServe(":8080", mux)
}
