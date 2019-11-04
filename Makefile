prod: ui-prod go-prod

ui-prod:
	cd ui && npm run build

go-prod:
	cd cmd/mailman && go generate && go build -tags prod