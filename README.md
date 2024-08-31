# Microblog

Microblog is a practice project written in Go without using any web framework and ORM. The main purpose of this project was to learn about building web applications.

## Features

- User registration and authentication
- Create, read, update, and delete blog posts
- Session management
- Server-side rendered pages

## Run Locally

Make sure you have Go (version 1.22) and MySQL (version 9.0) installed.

1. Clone the repository
    ```
    git clone https://github.com/anxxuj/microblog.git
    ```

2. Navigate to the project directory
    ```
    cd microblog
    ```

3. Install dependencies
    ```
    go mod tidy
    ```

### Setting up MySQL

1. Log in to MySQL as root (don't use `-p` flag if you haven't set root password)
    ```
    mysql -u root -p
    ```

2. Create a new database for the project
    ```sql
    CREATE DATABASE microblog;
    ```

3. Create a new user and grant privileges (enter username and password as it is)
    ```sql
    CREATE USER 'web'@'localhost' IDENTIFIED BY 'pass';
    GRANT SELECT, INSERT, UPDATE, DELETE ON microblog.* TO 'web'@'localhost';
    ```

4. Switch to the new database
    ```sql
    USE microblog;
    ```

5. Create the necessary tables
    ```sql
    CREATE TABLE users (
        id INT PRIMARY KEY AUTO_INCREMENT NOT NULL,
        username VARCHAR(255) NOT NULL,
        email VARCHAR(255) NOT NULL,
        password_hash CHAR(60) NOT NULL,
        created DATETIME NOT NULL
    );

    ALTER TABLE users ADD CONSTRAINT users_uc_username UNIQUE (username);
    ALTER TABLE users ADD CONSTRAINT users_uc_email UNIQUE (email);
    
    CREATE TABLE posts (
        id INT PRIMARY KEY AUTO_INCREMENT NOT NULL,
        title VARCHAR(255) NOT NULL,
        content TEXT NOT NULL,
        created DATETIME NOT NULL
    );

    CREATE TABLE sessions (
        token CHAR(43) PRIMARY KEY,
        data BLOB NOT NULL,
        expiry TIMESTAMP(6) NOT NULL
    );

    CREATE INDEX sessions_expiry_idx ON sessions (expiry);
    ```

6. Exit MySQL
    ```sql
    EXIT;
    ```

### Start the application

Make sure your MySQL server is running.

```
go run ./cmd/web
```

Open your web browser and navigate to http://localhost:4000
