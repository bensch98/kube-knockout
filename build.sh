#!/bin/bash

go build -o kubectl-knockout cmd/knockout/main.go
sudo mv kubectl-knockout /usr/local/bin
