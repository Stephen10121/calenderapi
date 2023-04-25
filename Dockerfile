# syntax=docker/dockerfile:1

FROM golang:1.19

# Set destination for COPY
WORKDIR /go/src/github.com/stephen10121/calendarapi

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY *.go ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-gs-ping

ARG PORT=9090
ARG SECRET=ekjwbfkb32kjhbdjknf32jkd3n2erkj

# Optional:
# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/engine/reference/builder/#expose
EXPOSE 9090

# Run
CMD ["/docker-gs-ping"]