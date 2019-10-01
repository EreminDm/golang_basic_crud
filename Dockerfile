

FROM golang:1.13-alpine
WORKDIR /crud
COPY . /crud
CMD ["go", "run", "./main.go", "-connectURI=192.168.99.100:27017", "-database=person", "-dbtype=mongo"]
EXPOSE 8000 8888
