FROM golang:1.15-alpine AS build
ADD . /app
ENV CGO_ENABLED 0
WORKDIR /app
RUN go build -o filer ./cmd/filer
RUN mkdir data && chmod 755 data

FROM scratch
COPY --from=build /app/filer /filer
COPY --from=build /app/data /data
CMD ["/filer"]
EXPOSE 9999