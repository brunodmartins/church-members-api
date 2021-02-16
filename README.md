# church-members-api

![Go](https://github.com/BrunoDM2943/church-members-api/workflows/Go/badge.svg)
[![codecov](https://codecov.io/gh/BrunoDM2943/church-members-api/branch/master/graph/badge.svg)](https://codecov.io/gh/BrunoDM2943/church-members-api)  [![Go Report Card](https://goreportcard.com/badge/github.com/BrunoDM2943/church-members-api?style=flat-square)](https://goreportcard.com/report/github.com/BrunoDM2943/church-members-api)
![Deploy](https://github.com/BrunoDM2943/church-members-api/workflows/Deploy/badge.svg)
![Docker Image CI](https://github.com/BrunoDM2943/church-members-api/workflows/Docker%20Image%20CI/badge.svg)


Its a simple API to manage a church's members. The API is written on Go 1.13 and the data is saved on MongoDB. It also supports Auth0, protecting against unathorized acess. 

## Installation

The API is containerized thanks to Docker. So, its so simple as build the image into a Docker Host and run. You only have to provided some environment parameters:

- APP_LANG - app language for i18n (supports only pt-BR and en)
- SCOPE - apps scope for resources usage (for production, use 'prod')
- VPR_CHURCH_MEMBERS_DATABASE_URL - URL to mongo database
- VPR_CHURCH_NAME - Church's name
- VPR_AUTH_ENABLE - enabled Auth0 support (true/false)
- VPR_AUTH_JWK, VPR_AUTH_ISS, VPR_AUTH_TOKEN, VPR_AUTH_AUD - Auth0 keys

## Improvements, Issue

You can create a PR for it =) 

## Questions

You can email me: bdm2943@gmail.com
