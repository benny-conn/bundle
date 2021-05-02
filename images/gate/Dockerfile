FROM golang
WORKDIR /bundle/gate
COPY . .
RUN go install cmd/gate/gate.go
EXPOSE 8020
CMD ["gate"]