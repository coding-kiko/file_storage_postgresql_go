A simple file storage api that can use either /pkg/repository/file_storage or PostgresSql as a database for files.

## How to use:

### Spinning up locally (standing on the root folder):
`docker-compose up -d` if you have docker installed in your machine, the postgres and adminer image will be pulled automatically 
`make storage=(storage) run`  (*storage=fs* for directory storage | *storage=postgres* for postgres storage)


### Register a new user:

    make email=(email) pwd=(pwd) register

### Authenticate:

     make email=(email) pwd=(pwd) register
    {"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJOYW1lIjoiZnJhbmNpc2NvLmNhbGl4dG9AZ2xvYmFudC5jb20iLCJleHAiOjE2NDg0OTQyNzgsImlhdCI6MTY0ODQ5Mzk3OCwiaXNzIjoibG9jYWxob3N0OjUwMDAvIn0.2A7PPAKBrrwl8eRlf5Eb7_ir481OGB388XKmDvAOM10"}

### Upload a file:
Using the token given after authentication.

`make token=12n...d&g file=(filename) upload` Notice: The token runs out after 5 minutes OR 5 requests. And the file must be inside the *./client* folder

### Download a file: 
Using the token given after authentication.
`make token=12n...d&g file=(filename) download` Notice: the file also gets downloaded inside *./client*. In case of error the file content will be the Json response.
