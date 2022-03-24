run:
	go run cmd/main.go -STORAGE=$(storage) 

upload:
	curl -X POST localhost:5000/upload \
	-H 'Authorization: Bearer $(token)' \
	-F "filename=@./client/$(file)"

download:
	curl localhost:5000/file?filename=$(file) \
	-H 'Authorization: Bearer $(token)' \
	--output ./client/$(file)

authenticate:
	curl -XPOST localhost:5000/authenticate \
	-H 'Content-Type:application/json' \
	-d '{"email":"francisco.calixto@globant.com", "pwd":"admin"}' \
