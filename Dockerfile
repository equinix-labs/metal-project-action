FROM golang as build
WORKDIR /code
COPY . .
RUN go build .

FROM scratch
COPY --from=build /code/packet-project-action /bin/packet-project-action
ENTRYPOINT ["packet-project-action"]
