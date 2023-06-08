# USPY 🕵️ - Backend
This is the official repository for the [USPY](https://uspy.me) Backend! Here you can find how to run the application youself and some brief explanation on the code repository.

&nbsp;

## **Package organization**
___
This repository is organized in these packages:

```
├─ config
├─ db
├─ entity
│     │ 
│     ├── controllers/
│     │   ├── account
│     │   ├── private
│     │   ├── public
│     │   └── restricted
│     │
│     ├── models/
│     │   ├── account
│     │   ├── private
│     │   ├── public
│     │   └── restricted
│     │
│     ├── views/
│     │   ├── account
│     │   ├── private
│     │   ├── public
│     │   └── restricted
│     │
│     └── validation
│
├─ config
├─ db
├─ entity
├─ iddigital
├─ server/
│     │ 
│     ├── controllers/
│     │   ├── account
│     │   ├── private
│     │   ├── public
│     │   └── restricted
│     │
│     ├── models/
│     │   ├── account
│     │   ├── private
│     │   ├── public
│     │   └── restricted
│     │
│     ├── views/
│     │   ├── account
│     │   ├── private
│     │   ├── public
│     │   └── restricted
│     │
│     └── middleware
│
└─ utils
```

Their respective responsibilities are the following:

#### **config**

    - Environment configuration/initialization

#### **db**

    - Database Access Object and database initialization

#### **entity**

    - All object definitions
    - follows a MVC architecture, see "server" package for more details
    - Contains subpackage validation, with input sanitization utilities.

#### **iddigital**

    - Wrapper functions for interacting with the USP iddigital API and Records' PDF parsing.

#### **server**

    - Endpoint closures and their implementations

    - Middleware contains useful middleware functions, such as JWT validation, rate limiting, data binding, etc

    - API Handlers and Data Access Objects are organized in a MVC manner:
        - controllers use the entity.controller objects to bind request data
        - models use the entity.models objects to recover data and perform database operations 
        - views use the entity.views objects to represent data the front-end will receive

    - All of these can be divided in the following manner:
        - account: all operations related to the user's account management, such as login, signup, delete, password recovery, etc
        - private: all operations related to the user's data management, such as getting/updating their grades and reviews
        - public: all operations related to data that is public (including non-registered users), such as subject data
        - restricted: all operations related to data that is anonymous yet visible to all registered-users

#### **utils**

    - Utility functions such as hashing functions and encoding stuff
    - Also contains testing utilities like the emulator functions

&nbsp;


## **Deployment & Execution**
___


To deploy and/or run this application, there are a few variables:


### **Environment variables**

| Name                   | Description                                     |    Required?     | Possible values |  Default Value  |
| :--------------------- | :---------------------------------------------- | :--------------: | :-------------: | :-------------: |
| **USPY_DOMAIN**        | Domain to run the web server                    |     **Yes**      |                 |   `localhost`   |
| **USPY_PORT**          | Port to run the web server                      |     **Yes**      |                 |   `8080`   |
| **USPY_JWT_SECRET**    | Private key to be used to generate `JWT` Tokens |     **Yes**      |                 |   `my_secret`   |
| **USPY_MODE**          | Which mode to run the web server                |     **Yes**      |  `[prod, dev, local]`  |      `local`      |
| **USPY_AES_KEY**       | Private AES key to be used for AES Encryption   |     **Yes**      |     AES key     |   `71deb5...`   |
| **USPY_AUTH_ENDPOINT** | Endpoint used to fetch PDF. See Auth section.   |     **Yes**  (*) |                 |                 |
| **USPY_RATE_LIMIT**    | `Frequency:Time` string for the rate-limiter    |      **No**      |  `F:P` string   |                 |
| **USPY_FIRESTORE_KEY** | Path to firestore access key                    | **Only locally** |                 |                 |
| **USPY_PROJECT_ID**    | GCP Project ID                                  | **In the Cloud** |                 |                 |
| **USPY_MAILJET_KEY**   | Mailjet key used for e-mail operations          | **In the Cloud** |                 |                 |
| **USPY_MAILJET_SECRET**| Mailjet secret used for e-mail operations       | **In the Cloud** |                 |                 |

### **(*)** This is only needed for signup. If you use the local deployment, you may have the built-in database population through docker-compose that automatically inserts data without signup.

&nbsp;

## **Running Locally**
___
### **Simple way**

To execute the webserver locally, simply run:

```sh
sudo docker-compose up --build -d
```

**P.S.:** Make sure to run docker-compose as sudo if you're loading previously exported files. These files are saved as root and running without sudo may fail due to permission errors.

**P.S².:** If this fails due to ".../docker/emulator/mount" not existing, try creating the folder before running docker-compose like so:

```sh
mkdir docker/emulator/mount
```

This will run three daemon containers, mapped to local ports:

- **firestore-emulator on 127.0.0.1:8200**
- **firestore-emulator (UI) on 127.0.0.1:4000**
- **uspy-backend on 127.0.0.1:8080**
- **uspy-scraper on 127.0.0.1:8300**

Some things to consider:

1. The firestore-emulator does not cover all features provided by the real database, therefore some things may not work as expected (e.g. anything that involves transactions is not guaranteed to work)
2. After the container initializes, the database will be empty, you can build its data using uspy-scraper by running

```sh
curl -X POST "localhost:8300/build?institute=55"
```

This operation should only take a few seconds to complete and it may fail due to errors on JupiterWeb. You can also use a different value other than 55 for scraping data from other institutes, but some other features may not work correctly and it may consume a lot of memory, since firestore will be hosted in memory.

&nbsp;


### **More advanced way: modifying the scraper!**

If you'd like to mess around with the data scraper source code and test your changes against the Backend, use the docker-compose file **located in the uspy-scraper repository** to spin up firestore and the scraper itself. Then, when running the backend, spin up **only** the uspy-backend container like so:

```sh
sudo docker-compose --build -d up uspy-backend
```

This will prevent the backend from running its own firestore and scraper containers and use the changed container and external network provided by the scraper itself.

&nbsp;

### **Cleanup**

To clean up:

```sh
docker-compose down
```

In order to save all of the data in firestore, run:

``sh
docker/emulator/save_db_data.sh
``

This will save the data inside `docker/emulator/mount/db_data`. Note that `docker/emulator/mount` is mounted into the firestore docker container. When running again the emulator, the saved data will be automatically reused as long as `docker/emulator/mount/db_data` exists. The `IMPORT_DATA` environment variable can be changed to customize the `mount/db_data` path.

&nbsp;

## **Testing**
___


To run tests, you must set up the firestore emulator. Folow these steps:

#### Install the Firebase CLI

Info on how to install here: [Firebase installation reference](https://firebase.google.com/docs/cli#install-cli-mac-linux)

#### Set up a `firebase.json` file (if you don't, the default port 8080 will be used for the emulator)

```json
{
  "emulators": {
    "firestore": {
      "port": <your_port_of_choice>
    },
    "ui": {
      "enabled": <do_you_want_the_ui?>
    }
  }
}
```

#### Run tests

```sh
chmod u+x test.sh && ./test.sh
```


&nbsp;


## **Cloud Services**
---

The following services are used by the backend application

### **Firestore**:

    - Non relational database. Used to store all persistent data.
    - Must be accessed with an IAM key when running locally or just with the project ID if in production

### **Cloud run**:

    - Serverless application that will run the web server
    - Can be set up manually, but also through cloud build using the cloubuild.yaml configuration file
    - Runs the web server by building the container using the Dockerfile in the repository

### **Auth**:

    - Signup is done in a 2-way step that requires an external service (uspy-auth)
    - First step is fetching of user PDF, grades parsing and data insertion
        - This step requires an external service
        - This service should have a root endpoint "/:authCode" that, given an authCode,   responds with the user's PDF. Contact us if you'd like to understand more about this.
    - Second step is auth registering (inserting email, password and verifying user), does not require external services

&nbsp;

## **How to contribute**
---

### **Features, requests, bug reports**

If this is the case, please submit an issue through the [contributions repository](github.com/Projeto-USPY/uspy-contributions/issues).

### **Actual code**

If you'd to directly contribute, fork this repository and create a pull request to merge on `dev` branch. Please do not submit pull requests to the main branch as they will be denied. The main branch is used for releases and we don't really push to it other than through the `deploy.sh` script.

If you'd to directly contribute, fork this repository and create a pull request to merge on `dev` branch. Please do not submit pull requests to the main branch as they will be denied. The main branch is used for releases and we don't really push to it other than through the `deploy.sh` script.

We really appreciate any contributors! This project is from USP students and for USP students! If you have any questions or would simply like to chat, contact us on Telegram @preischadt @lucsturci
