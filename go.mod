module github.com/0ppliger/open-asset-gateway

go 1.24.0

replace github.com/owasp-amass/asset-db => ../asset-db

require (
	github.com/owasp-amass/asset-db v0.23.1
	github.com/owasp-amass/open-asset-model v0.15.0
)

require (
	github.com/google/uuid v1.6.0 // indirect
	github.com/neo4j/neo4j-go-driver/v5 v5.28.4 // indirect
)
