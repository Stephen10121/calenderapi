FROM golang:1.19

# Set the Current Working Directory inside the container
WORKDIR $GOPATH/src/github.com/stephen10121/calendarapi

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

ARG PORT=9090
ARG SECRET=ekjwbfkb32kjhbdjknf32jkd3n2erkj

# Download all the dependencies
# RUN go get -d -v ./...

# Install the package
RUN go build -o ./out/calendarapi .


# This container exposes port 8080 to the outside world
EXPOSE 9090

# Run the binary program produced by `go install`
CMD ["./out/go-sample-app"]