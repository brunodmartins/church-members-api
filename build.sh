#!/bin/bash
GOOS=linux CGO_ENABLED=0 go build
zip function.zip church-members-api