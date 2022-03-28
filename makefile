# all of the following commands are relative to the root folder of the project

run:
	go run cmd/main.go -STORAGE=$(storage) 

upload:
	curl -X POST localhost:5000/upload \
	-H 'Authorization: Bearer $(token)' \
	-F "filename=@./client/$(file)"

# NOTE: the response is stored in the file, wether successful or not
# could improve with a bash script -> store to file in case of success, output to terminal if it fails
download:
	curl localhost:5000/file?filename=$(file) \
	-H 'Authorization: Bearer $(token)' \
	-o ./client/$(file)

authenticate:
	curl -XPOST localhost:5000/authenticate \
	-H 'Content-Type:application/json' \
	-d '{"email":"francisco.calixto@globant.com", "pwd":"admin"}'

register:
	curl -XPOST localhost:5000/register \
	-H 'Content-Type:application/json' \
	-d '{"email":"$(email))", "pwd":"$(pwd))"}'