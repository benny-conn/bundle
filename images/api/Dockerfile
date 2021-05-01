FROM golang
WORKDIR /bundle/api
COPY . .
RUN go install cmd/api/api.go
EXPOSE 8020
CMD ["api"]