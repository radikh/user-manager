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

##### Database migration
PostgreSQL database is migrate from oficial image postgres:12.2. Initial structure of database is done with  ["Database migrations. CLI and Golang library." ](https://github.com/golang-migrate/migrate ) using [Docker image for golang-migrate binary](https://hub.docker.com/r/migrate/migrate/). 

Queries are placed in folder "migrations" and must named as N____.up.sql or N___.down.sql, where N means the order of executing. 

Execution of migration  
```bash
docker-compose up -d
```

#### Admin panel

Service have admin command line tool to manipulate accounts with admin rights.
CLI tool is capable:
- Create a user
- Delete a user by login
- Disable a user by login
- Get user information by login except of password hash and salt 

CLI tool is implemented both as a part of the service and as a separate tool.  

USAGE:
-    umcli command [arguments...]

COMMANDS:
-    create    Create new account in database
-    delete    Delete account in database
-    disable   Disable account in database
-    activate  Activate previously disabled account in database
-    update    Update account in database
-    info      Show information about user stored in database
-    help, h   Shows a list of commands or help for one command

GLOBAL OPTIONS:
-    --help, -h  show help (default: false)

Execution of CLI in starting docker container  
```bash
docker-compose exec  usermanager sh -c "/opt/services/user-manager/cmd/umcli create login=UserLogin pwd=UserPassword email=UserEmail@company.com phone=7777 name=UserName lastname=UserLastName"
```
```bash
docker-compose exec  usermanager sh -c "/opt/services/user-manager/cmd/umcli update login=UserLogin pwd=UserPassword email=UserEmail@company.com phone=7777 name=UserName lastname=UserLastName"
```
```bash
docker-compose exec  usermanager sh -c "/opt/services/user-manager/cmd/umcli delete login=UserLogin"
```
```bash
docker-compose exec  usermanager sh -c "/opt/services/user-manager/cmd/umcli info login=UserLogin"
```

#### Start all services

You must create .env  file in the root of project with environment variables as in 

    env.example

**Note!!! Your .env file should be without comments**

Commands for build all infrastructure

    make up

Some other commands  for working with this project your can find in Makefile