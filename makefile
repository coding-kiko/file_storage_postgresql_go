run:
	PORT=5000 go run cmd/main.go

client:
	FILENAME=$(filename) go run client/client.go

clean:
	rm ./tmp/*