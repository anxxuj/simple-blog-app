## Run Locally

Make sure you have Go (version 1.22) and MySQL (version 9.0) installed.

1. Clone the repository and navigate to the project directory
    ```
    $ git clone git@github.com:anxxuj/microblog.git
    $ cd microblog
    ```

2. Install dependencies
    ```
    $ go mod tidy
    ```

### Setting up MySQL

1. Log in to MySQL as root
    ```
    $ mysql -u root -p
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
    CREATE TABLE sessions (
        token CHAR(43) PRIMARY KEY,
        data BLOB NOT NULL,
        expiry TIMESTAMP(6) NOT NULL
    );

    CREATE INDEX sessions_expiry_idx ON sessions (expiry);

    CREATE TABLE users (
        id INT PRIMARY KEY AUTO_INCREMENT NOT NULL,
        username VARCHAR(255) NOT NULL,
        email VARCHAR(255) NOT NULL,
        password_hash CHAR(60) NOT NULL
    );
    
    ALTER TABLE users ADD CONSTRAINT users_uc_username UNIQUE (username);
    ALTER TABLE users ADD CONSTRAINT users_uc_email UNIQUE (email);

    CREATE TABLE posts (
        id INT PRIMARY KEY AUTO_INCREMENT NOT NULL,
        title VARCHAR(255) NOT NULL,
        content TEXT NOT NULL,
        created DATETIME NOT NULL
    );
    ```

6. Exit MySQL
    ```sql
    EXIT;
    ```

### Start the application

1. Make sure your MySQL server is running.
    ```
    $ go run ./cmd/web
    ```

2. Open your web browser and navigate to http://localhost:4000
