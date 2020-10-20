FROM golang as build
WORKDIR /code
COPY . .
ENV CGO_ENABLED=0
RUN go build .

FROM alpine
COPY --from=build /code/metal-project-action /bin/metal-project-action
CMD ["metal-project-action"]
