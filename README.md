# GoTube

A Go-based application that allows users to upload and watch videos.


## Installation

Install the dependencies:
```
$ go mod download -json
```

## Migrations

Install [Buffalo](https://gobuffalo.io/documentation/getting_started/installation/).

Update the **./database.yml** file.

Create the databases.

Run the migrations:
```
$ soda migrate up -e development
$ soda migrate up -e test
$ soda migrate up -e production
```

## Run the app

Start the web server in development mode on port 3000:
```
$ gin --appPort 8080 --port 3000
```

Finally, load [http://localhost:3000](http://localhost:3000) in your browser.

### Run the tests

```
$ go test ./...
```
