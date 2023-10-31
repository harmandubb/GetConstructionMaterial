FROM ubuntu:latest
WORKDIR /app

RUN apt update
RUN apt install -y golang-go 

RUN apt install -y git ca-certificates \
    && update-ca-certificates \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Copy go mod and sum files
COPY go.mod go.sum ./

# Install dependencies
RUN go mod download

# Copy the source code
COPY . .

# explicitally copy the local .env file (This should be removed/looked into when running the actual server)
# --env-file


# CMD ["bash"]

CMD ["go", "run", "main.go"]