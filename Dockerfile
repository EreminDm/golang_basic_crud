FROM golang:1.13-alpine
WORKDIR /crud
COPY . /crud
RUN go-wrapper download
RUN go-wrapper install
CMD ["go-wrapper", "run", "-connectURI=localhost:27017", "-database_name=person", "-dbtype=mongo"]
EXPOSE 8000 8888