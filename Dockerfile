FROM golang as build
WORKDIR /code
COPY . .
ENV CGO_ENABLED=0
RUN go build .

FROM alpine
COPY --from=build /code/packet-project-action /bin/packet-project-action
CMD ["packet-project-action"]
