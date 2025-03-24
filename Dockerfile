FROM golang:1.22

RUN apt-get update && \
        apt-get install -y libx11-dev && \
        apt-get clean && \
        rm -rf /var/lib/apt/lists/*

# Creates an app directory to hold your app’s source code
WORKDIR /app

# Copies everything from your root directory into /app
COPY . .

# Installs Go dependencies
RUN cd src && go mod download && go get

# Builds your app with optional configuration
RUN cd src && go build -o fl

# Specifies the executable command that runs when the container starts
CMD [ “/fl” ]
