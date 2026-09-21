run:
	go run .

debug:
	dlv debug \
	--headless \
	--listen=127.0.0.1:2345 \
	--api-version=2 \
	--accept-multiclient \
	.
