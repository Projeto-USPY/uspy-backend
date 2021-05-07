# USPY 🕵️ - Backend

This is the official repository for the [USPY](https://uspy.me) Backend! Here you can find how to run the application youself and some brief explanation on the code repository.

## Package organization

This repository is organized in these packages:

```
config
db
entity
iddigital
main
server/
├── controllers/
│   ├── account
│   ├── private
│   ├── public
│   └── restricted
├── middleware
└── models/
    ├── account
    ├── private
    ├── public
    └── restricted
utils
```

Their respective responsibilities are the following:

#### **config**

    - Environment configuration/initialization

#### **db**

    - Database Access Object and database initialization

#### **entity**

    - All object definitions and their bindings to HTTP requests and their database objects

#### **iddigital**

    - Wrapper functions for interacting with the USP iddigital API and Records' PDF parsing.

#### **main**

    - Entrypoint for the web server and backend endpoints definitions

#### **server**

    - Endpoint closures and their implementations
    - Middleware contains useful middleware functions, such as JWT validation and data binding
    - Endpoint implementation uses a Model-Controller structure
    - This M-C can be divided into four categories:
        - Account: all operations related to the user's account management, such as login, signup, delete, etc
        - Private: all operations related to the user's data management, such as their grades and reviews
        - Public: all operations related to data that is public (including non-registered users), such as subject data
        - Restricted: all operations related to data that is only accessible by registered users (aggregated data)

#### **utils**

    - Utility functions such as hashing functions and encoding stuff

## Deployment & Execution

To deploy and/or run this application, there are a few requisites:

### Environment variables

| Name                   | Description                                     |    Required?     | Possible values |  Default Value  |
| :--------------------- | :---------------------------------------------- | :--------------: | :-------------: | :-------------: |
| **USPY_DOMAIN**        | Domain to run the web server                    |     **Yes**      |                 |   `localhost`   |
| **USPY_PORT**          | Port to run the web server                      |     **Yes**      |                 |   `localhost`   |
| **USPY_JWT_SECRET**    | Private key to be used to generate `JWT` Tokens |     **Yes**      |                 |   `my_secret`   |
| **USPY_MODE**          | Which mode to run the web server                |     **Yes**      |  `[prod, dev]`  |      `dev`      |
| **USPY_AES_KEY**       | Private AES key to be used for AES Encryption   |     **Yes**      |     AES key     |   `71deb5...`   |
| **USPY_RATE_LIMIT**    | `Frequency:Time` string for the rate-limiter    |      **No**      |  `F:P` string   |                 |
| **USPY_FIRESTORE_KEY** | Path to firestore access key                    | **Only locally** |                 |                 |
| **USPY_PROJECT_ID**    | GCP Project ID                                  | **In the Cloud** |                 |                 |

### Testing

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

`chmod u+x test.sh && ./test.sh`

### Cloud Services

The following services are used by the backend application:

### Firestore:

    - Non relational database. Used to store all persistent data.
    - Must be accessed with an IAM key when running locally or just with the project ID if in production

### Cloud run:

    - Serverless application that will run the web server
    - Can be set up manually, but also through cloud build using the cloubuild.yaml configuration file
    - Runs the web server by building the container using the Dockerfile in the repository

## How to contribute

### Features, requests, bug reports

If this is the case, please submit an issue through the [contributions repository](github.com/Projeto-USPY/uspy-contributions/issues).

### Actual code

Although we are not yet ready for community contributions, you **could** submit a pull requests and we'll analyze it through =).
