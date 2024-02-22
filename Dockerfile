FROM golang:1.19-alpine

WORKDIR /payment-app

COPY . .

# Build the Go app
RUN go build -o main main.go

# Command to run the executable
EXPOSE 8080
CMD [ "/payment-app/main" ]