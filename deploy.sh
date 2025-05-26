#!/usr/bin/env bash
set -euo pipefail

PROJECT_ID=${GCP_PROJECT:-your-gcp-project}
REGION=${GCP_REGION:-us-central1}

# Build & push Docker images for all services
make build-all
make docker.build-all

docker push gcr.io/$PROJECT_ID/sim-auth-token-broker-broker:latest
docker push gcr.io/$PROJECT_ID/sim-auth-token-broker-partner:latest
docker push gcr.io/$PROJECT_ID/sim-auth-token-broker-cellcom:latest
docker push gcr.io/$PROJECT_ID/sim-auth-token-broker-pelephone:latest

# Deploy services to Cloud Run
gcloud run deploy sim-auth-token-broker-broker \
  --image gcr.io/$PROJECT_ID/sim-auth-token-broker-broker:latest \
  --region $REGION --platform managed --port 8080 \
  --set-secrets BROKER_JWT_SECRET=projects/$PROJECT_ID/secrets/BROKER_JWT_SECRET:latest

gcloud run deploy sim-auth-token-broker-partner \
  --image gcr.io/$PROJECT_ID/sim-auth-token-broker-partner:latest \
  --region $REGION --platform managed --port 8081

gcloud run deploy sim-auth-token-broker-cellcom \
  --image gcr.io/$PROJECT_ID/sim-auth-token-broker-cellcom:latest \
  --region $REGION --platform managed --port 8082

gcloud run deploy sim-auth-token-broker-pelephone \
  --image gcr.io/$PROJECT_ID/sim-auth-token-broker-pelephone:latest \
  --region $REGION --platform managed --port 8083
