# Go-Starter-Kit
Inspired from [go-starter-kit](https://gostarterkit.com/)
This project is a go starter kit / golang boilerplate that uses the golang standards project layout.
It builds a foundation for golang API using packages gin webframework, gorm, postgres and jwt-go for Authentication.

## **Features**
1. Routing using gin
2. Request Validation
3. Filtering / Pagination / Ordering
4. Authentication using jwt-go
5. Crash Alerts via email using gomail
6. Logging
7. Error Handling
8. Dockerfile
9. Makefile
10. Linter


## **Dependencies**
* Golang
* Postgres
* Docker

## **Quick Start**

Clone the repo use env.sample file to set env variables
Start postgres and docker at your machine then run
`make start`

This will automatically build docker image of this project (if it does not exist already) and run it at port 4000

## **Build and Run**

To only build first 
`make build`
then run
`make run`

You can check for other targets in Makefile.

## **Build and Run without Docker**

You can also run this without docker using

`go build -o bin/` 

This will create a binary named `go-starter-kit` in newly created `bin` directory
Run it
`./bin/go-starter-kit` 