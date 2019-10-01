FROM golang:1.13-alpine
WORKDIR /crud
COPY go.mod go.sum ./crud
RUN go mod download
COPY . /crud
CMD ["go", "run", "-connectURI=localhost:27017", "-database_name=person", "-dbtype=mongo"]
EXPOSE 8000 8888