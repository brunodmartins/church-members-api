# church-members-api

![Go](https://github.com/BrunoDM2943/church-members-api/workflows/Go/badge.svg)
[![codecov](https://codecov.io/gh/BrunoDM2943/church-members-api/branch/master/graph/badge.svg)](https://codecov.io/gh/BrunoDM2943/church-members-api)  [![Go Report Card](https://goreportcard.com/badge/github.com/BrunoDM2943/church-members-api?style=flat-square)](https://goreportcard.com/report/github.com/BrunoDM2943/church-members-api)
![Docker Image CI](https://github.com/BrunoDM2943/church-members-api/workflows/Docker%20Image%20CI/badge.svg)


Its a simple API to manage a church's members. The API is written on Go 1.13 and the data is saved on DynamoDB. The application is supported to be run on AWS, using DynamoDB, Lambda and API Gateway

## Installation

The API is containerized thanks to Docker. So, its so simple as build the image into a Docker Host and run. You only have to provided some environment parameters:

- APP_LANG - app language for i18n (supports only pt-BR and en)
- SCOPE - apps scope for resources usage (for production, use 'prod')

## Improvements, Issue

You can create a PR for it =) 

## Questions

You can email me: bdm2943@gmail.com
