## About this project
A file storage and processing app that resizes an image to 1000 width and preserves its ratio. The image and processing records are stored in a database using PostgresSql. The aim of this project was to create a file Api trying to use no third party libraries, which I followed to a certain extent, and implement new technologies. Even though the different "services" are all part of the same repository, the goal was to make use of Go's interfaces to separate concerns the best I could and try to replicate the same that could be achieved if they were developed as truly separated services. \

<div align="center">
  <img alt="JPG" src="readme/diagram.jpg"/>
</div>

## How to use:
The only requirements for running this project are cloning this repo and having docker installed on your machine.
### Spinning up locally (standing on the root folder):

    sudo docker build -t file-services .
    sudo docker-compose up -d
Now you have the service, redis, postgres and rabbitmq instances running on docker containers.
### Making use of the api
For testing the api I decided to go with the curl, the specific command for each endpoint can be found inside the Makefile. 
#### Register a new user:

    make email=(email) pwd=(pwd) register

### Authenticate:

     make email=(email) pwd=(pwd) authenticate
    {"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJOYW1lIjoiZnJhbmNpc2NvLmNhbGl4dG9AZ2xvYmFudC5jb20iLCJleHAiOjE2NDg0OTQyNzgsImlhdCI6MTY0ODQ5Mzk3OCwiaXNzIjoibG9jYWxob3N0OjUwMDAvIn0.2A7PPAKBrrwl8eRlf5Eb7_ir481OGB388XKmDvAOM10"}
This token will last for 15 minutes or 5 requests. 
### Upload a file:
Using the token given after authentication.

`make token=X2n...dEg file=(filename) upload` Notice: The file must be inside the *./client* folder. 

### Download a file: 
Using the token given after authentication.

`make token=X2n...dEg file=(filename) output=(output_filename) download` Notice: the file also gets downloaded inside *./client* and is already resized. In case of error the file content will be the Json response. 
