# Web Service Go Gin

## Example .env file

```sh
MONGO_URI="mongodb://username:password@localhost:27017/"
MONGO_DATABASE="MyDatabase"

PORT="8080"

# MODES: release, debug, test
MODE="release"
```

## Example Dockerfile

```dockerfile
FROM golang:1.22.0-alpine3.19

WORKDIR /app

# pre-copy/cache go.mod for pre-downloading dependencies and only re-downloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o ./app
RUN chmod +x ./app

CMD ["./app"]
```

## Example docker-compose.yml

```yml
version: "3.9"
services:
  app:
    container_name: go-api
    restart: unless-stopped
    build: .
    depends_on:
      mongodb-api:
        condition: service_healthy
    ports:
      - "8080:8080"
  mongodb-api:
    container_name: mongodb-api
    image: mongo:7.0.5-jammy
    restart: unless-stopped
    environment:
      MONGO_INITDB_ROOT_USERNAME: username
      MONGO_INITDB_ROOT_PASSWORD: password
    healthcheck:
      test: echo 'db.runCommand("ping").ok' | mongosh 127.0.0.1:27017 --quiet
      timeout: 5s
      retries: 10
    volumes:
      - ./appdata/go-api/mongodb:/data/db
      - ./mongo-init.js:/docker-entrypoint-initdb.d/mongo-init.js:ro
```

## Example mongodb-init.js

```javascript
// Authenticate as admin
db.getSiblingDB("admin").auth("username", "password");

// Get GameData database
db = db.getSiblingDB("MyDatabase");

// Create a unique index on username
// will create collection if it does not already exist
db.users.createIndex({ username: "text" }, { unique:true })

// Insert the admin
db.users.insertOne({
  _id: new ObjectId(),
  username: "admin",
  password: "$2a$14$RyxkElOBRzEz3.P4PzLbtuqSdvuJdbRBrFBASSeBbL8MaMovmFaju",
  name: "Admin Example Name",
  email: "admin@example.com",
  role: "admin",
  token: "random-initial-token",
  last_seen: ISODate(),
  since: ISODate(),
  points: 500,
});
```
