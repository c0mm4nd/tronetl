# Currently using ubuntu for usability.
# BUILDER
FROM golang:latest as builder

ARG GOPROXY

COPY . /tronetl
WORKDIR /tronetl

# RUN apt install gcc -y
RUN GOPROXY=$GOPROXY go build .

# MAIN
FROM ubuntu:latest

COPY --from=builder /tronetl/tronetl /usr/local/bin/
WORKDIR /workspace

ENTRYPOINT ["tronetl"]
