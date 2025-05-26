.PHONY: dev-broker dev-partner dev-cellcom dev-all
.SILENT:

dev-broker:
	# Run broker locally
	go run services/broker/main.go

dev-partner:
	# Run partner locally
	go run services/partner/main.go

dev-cellcom:
	# Run cellcom locally
	go run services/cellcom/main.go

dev-pelephone:
	# Run pelephone locally
	go run services/pelephone/main.go

dev-all:
	# Launch all services concurrently
	$(MAKE) dev-partner & \
	$(MAKE) dev-cellcom & \
	$(MAKE) dev-broker & \
	$(MAKE) dev-pelephone & \
	wait
