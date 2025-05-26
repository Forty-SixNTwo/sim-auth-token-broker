.PHONY: go.build-broker go.build-partner go.build-cellcom go.build-pelephone go.build-all \
        docker.build-broker docker.build-partner docker.build-cellcom docker.build-pelephone docker.build-all

# Go build targets
go.build-broker:
	go build -o bin/broker services/broker/main.go

go.build-partner:
	go build -o bin/partner services/partner/main.go

go.build-cellcom:
	go build -o bin/cellcom services/cellcom/main.go

go.build-pelephone:
	go build -o bin/pelephone services/pelephone/main.go

go.build-all: go.build-broker go.build-partner go.build-cellcom go.build-pelephone

# Docker build targets
docker.build-broker:
	docker build \
	  --build-arg SERVICE_NAME=broker \
	  --build-arg PORT=8080 \
	  -t gcr.io/$(PROJECT_ID)/sim-auth-token-broker-broker:latest .

docker.build-partner:
	docker build \
	  --build-arg SERVICE_NAME=partner \
	  --build-arg PORT=8081 \
	  -t gcr.io/$(PROJECT_ID)/sim-auth-token-broker-partner:latest .

docker.build-cellcom:
	docker build \
	  --build-arg SERVICE_NAME=cellcom \
	  --build-arg PORT=8082 \
	  -t gcr.io/$(PROJECT_ID)/sim-auth-token-broker-cellcom:latest .

docker.build-pelephone:
	docker build \
	  --build-arg SERVICE_NAME=pelephone \
	  --build-arg PORT=8083 \
	  -t gcr.io/$(PROJECT_ID)/sim-auth-token-broker-pelephone:latest .

docker.build-all: docker.build-broker docker.build-partner docker.build-cellcom docker.build-pelephone
