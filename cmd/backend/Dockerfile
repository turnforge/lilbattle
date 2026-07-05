FROM golang:latest
# FROM panyam/golangdev:latest

WORKDIR /app

COPY locallinks ./locallinks
COPY Makefile go.mod go.sum ./
RUN mkdir -p cmd clients datastore services gen locallinks
COPY cmd ./cmd
COPY assets ./assets
COPY gen ./gen
COPY locallinks ./locallinks
COPY services ./services
COPY utils ./utils
COPY lib ./lib
COPY main.go ./main.go
# NOTE: env files (configs/.env, configs/.env.dev) are NOT copied at
# build time — they're bind-mounted at runtime via docker-compose so
# secrets stay off the image and edits don't require rebuild. main.go
# tolerates a missing configs/.env.dev in dev mode (see LILBATTLE_ENV).

RUN mkdir -p web web/server web/gen web/static
COPY web/server ./web/server ./web/static
COPY web/gen ./web/gen
COPY web/templates ./web/templates
COPY web/assets ./web/assets

RUN go mod download

# WORKDIR stays at /app so `go run main.go` resolves to /app/main.go
# and godotenv finds configs/.env.dev relative to the working dir.

# Uncomment below line to eee if the build works as a test
# Usually this will be useful if you have added or removed folders
# RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /tmp/main main.go

CMD ["go", "run", "main.go"]
