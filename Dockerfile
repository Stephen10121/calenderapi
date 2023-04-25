FROM golang:1.19

# Set the Current Working Directory inside the container
WORKDIR $GOPATH/src/github.com/stephen10121/calendarapi

# Copy everything from the current directory to the PWD (Present Working Directory) inside the container
COPY . .

ARG PORT=9090
ARG SECRET=ekjwbfkb32kjhbdjknf32jkd3n2erkj

# Download all the dependencies
# RUN go get -d -v ./...

# Install the package
RUN go install -v ./...

# This container exposes port 8080 to the outside world
EXPOSE 9090

# Run the executable
CMD ["go-sample-app"]