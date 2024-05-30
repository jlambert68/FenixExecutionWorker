# Compile stage
FROM golang:1.21 AS build-env

ADD . /dockerdev
WORKDIR /dockerdev

COPY go.* ./
RUN go mod tidy

RUN go build -o /fenixWorker .

# Final stage
FROM debian:bookworm
#FROM golang:1.13.8

EXPOSE 6671
#FROM golang:1.13.8
WORKDIR /
COPY --from=build-env /fenixWorker /
#Add data/ data/

#CMD ["/fenixClientServer"]
ENTRYPOINT ["/fenixWorker"]



#// docker build -t  fenix-client-server .
#// docker run -p 5998:5998 -it  fenix-client-server
#// docker run -p 5998:5998 -it --env StartupType=LOCALHOST_DOCKER fenix-client-server

#//docker run --name fenix-client-server --rm -i -t fenix-client-server  bash