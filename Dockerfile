FROM golang as build
WORKDIR /code
COPY . .
# Need to manually bump this before each release
ENV ACTION_VERSION=0.12.1
ENV CGO_ENABLED=0
RUN go build -ldflags "-X 'github.com/equinix-labs/metal-project-action/internal.version=${ACTION_VERSION}'"

FROM alpine
COPY --from=build /code/metal-project-action /bin/metal-project-action
CMD ["metal-project-action"]
