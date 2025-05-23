# Web Server

A multithreaded web server in Go for a Wiki. 

At first, I started by following the [official](https://go.dev/doc/articles/wiki/) Golang tutorial.

In that, a web server is built to allow a user to create, edit, and view text files on a wiki.

I then enhanced it by doing the following:
- Changed it from storing text files that clump up the directory to writing to a PostgreSQL server on my PC.
- Abstracted away HTML template executions to make the code more modular & easier to maintain, keeping the frontend inside of "templates" and backend data management inside of "data".

### Setup

To run the web server successfully the primary prerequisite is Golang, and a PostgreSQL server.

Once PostgreSQL is installed, you need to create a password for the superuser `posgres`.

Then connect to the server via the `psql -U postgres` command and enter the password when prompted.

Finally, to setup the database and table:

```SQL
CREATE DATABASE webserver_data;

\c my_database

CREATE TABLE pages (
    title VARCHAR(255) PRIMARY KEY,
    body TEXT
);
```

### Usage

To run, simply do:

```
go run main.go
```

Then visit [http://localhost:8080/view/ANewPage](http://localhost:8080/view/ANewPage).