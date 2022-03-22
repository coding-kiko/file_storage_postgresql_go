run:
	go run cmd/main.go -STORAGE=$(storage) 

client:
	FILENAME=$(filename) go run client/client.go

clean:
	rm ./tmp/*