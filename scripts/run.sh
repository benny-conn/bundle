#!/bin/bash

export DATABASE_PORT=8040
export REPO_PORT=8060
export WEB_PORT=8080
export GATE_PORT=8020
export GATE_HOST=localhost
export WEB_HOST=localhost
export REPO_HOST=localhost
export DATABASE_HOST=localhost
export MODE=DEV
export AUTH0_ID=MpOxXFrk5XhR7gKcWIhYVZTNDDinx4ZT
export AUTH0_SECRET=jdksX1I0hZ8vej4M6LW-VRxtIiRFVXr2MMVYK0K9FvD8EtsiiRfATnKszcb2SvrG
export AUTH0_AUD=https://bundlemc.io/auth
export MONGO_URL="mongodb+srv://benny-bundle:thisismypassword1@bundle.mveuj.mongodb.net/main?retryWrites=true&w=majority"
export AWS_REGION=us-east-1
export AWS_BUCKET=bundle-repository


go run cmd/web/web.go &
go run cmd/db/db.go &
go run cmd/repo/repo.go &
go run cmd/gate/gate.go && fg
