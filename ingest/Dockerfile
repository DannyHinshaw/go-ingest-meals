FROM golang:1.13.6

RUN apt-get -y update && apt-get -y upgrade && \
    apt-get -y install bash git openssh-client openssh-server

# Add the source
RUN mkdir /app
RUN ls -al

ADD . /app/
WORKDIR /app

# Build the go executable
RUN go build -o ingest .
RUN ls -al

# Make create app user and run it.
RUN chmod +x ingest
RUN chmod +x wait-for-it.sh
RUN useradd -ms /bin/bash appuser

USER appuser
