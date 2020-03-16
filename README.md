## User management service (UM)

User management service stores user related context and credentials. It provides a REST API to perform a set of CRUD to manage users and an endpoint to authenticate. All users data will be stored in a database.

#### User model:

User is the central model of the service: 

```
    User {
        id              // Unique identifier
        username        // Unique login
        password        // Strong enough password
        email           // Valid email           
        first name      // Obviously first name
        last name       // Obviously last name
        phone           // Valid phone
        created at      // Time when user was created
        updated at      // Time of last changes made
    }
```

#### Storage

Service will use PostgreSQL as a storage for user data. All passwords have to be saved securely using hashing algorithms and saults. Access rights for the service SQL user have to be exactly the same that required to cover needed queries requirements. Database connection have to be able to recover after disconnect.

###### Database migration
PostgreSQL database is migrate from oficial image postgres:12.2. Initial structure of database is done with "Database migrations. CLI and Golang library." https://github.com/golang-migrate/migrate using Docker image for golang-migrate binary  https://hub.docker.com/r/migrate/migrate/.
Queries are placed in folder "migrations" and must named as N____.up.sql or N___.down.sql, where N means the order of executing.
Execution of migration
    #docker-compose up -d

#### Admin panel

Service should have admin command line tool to manipulate accounts with admin rights.
CLI tool should be capable:
- Create a user
- Delete a user by login
- Disable a user by login
- Get user information by login except of password hash and salt 

CLI tool can be implemented both as a part of the service or as a separate tool.
