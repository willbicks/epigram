# Configuration

The Epigram server is primarily configured by a YAML file, while certain parameters can be overwritten by environment variables.

On Windows, the default configuration file locaiton is `.\config.yml`, while on Linux, the default locaiton is `/etc/epigram/config.yml`. On both platforms, the config file location can be overwritten by the `EP_CONFIG` environment variable.

## Configuration 'Merging'

Many configuratio parameters can be set from multiple sources. The following list outlines the order in which configuration parameters are merged. If a parameter is specified by multiple sources, the value from the last source (highest number) will be used.

1. Default values
2. Configuration file
3. Environment variables

## Configuration parameters

The following table outlines parameters which can be configured, as well as their corresponding environment variables (if applicable), and their default values.

| Parameter                                                                                                                                                                       | YAML key      | Environment variable | Default value                                                                                                                    |
| ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------- | -------------------- | -------------------------------------------------------------------------------------------------------------------------------- |
| **Address** to listen for incoming requsts on.                                                                                                                                  | `address`     | `EP_ADDRESS`         | 0.0.0.0                                                                                                                          |
| **Port** to listen for incoming requests on.                                                                                                                                    | `port`        | `EP_PORT`            | 80                                                                                                                               |
| **BaseURL** is the complete domain and path to access the root of the web server, used for creating callback URLs                                                               | `baseURL`     | `EP_BASEURL`         |                                                                                                                                  |
| **Title** is the name of the applicaiton to be shown in the frontend.                                                                                                           | `title`       | `EP_TITLE`           | Epigram                                                                                                                          |
| **Description** is a short description of the application to be shown in the frontend.                                                                                          | `description` | `EP_DESCRIPTION`     | Epigram is a simple web service for communities to immortalize the enlightening, funny, or downright dumb quotes that they hear. |
| **Repo** dictates what type of storage the application should use for data persistence. (either 'inmemory' or 'sqlite')                                                         | `repo`        | `EP_REPO`            | inmemory                                                                                                                         |
| **DBLoc** is the location where the database can be found. In the case of an SQLite repository, this is the path to database file. It has no effect on an in-memory repository. | `DBLoc`       | `EP_DBLOC`           | Unix: `/var/epigram/epigram.db`, Windows: `.\epigram.db`                                                                         |
| **TrustProxy** dictates whether `X-Forwarded-For` header should be trusted to obtain the client IP, or if the requestor IP shoud be used instead.                               | `trustProxy`  | `EP_TRUSTPROXY`      | false                                                                                                                            |

### OIDC Provider Configuration

Additionally, an OpenID Connect provider is required to authenticate users. The following parameters are required to configure the OIDC provider. Thsese parameters cannot be set via environment variables, and have no default values. They should be specified in the configruation file as a map under the `OIDCProvider` key.

| Parameter                                                        | YAML key       | Example value                         |
| ---------------------------------------------------------------- | -------------- | ------------------------------------- |
| **Name** of the OIDC provider, used to build it's callback URL.  | `name`         | google                                |
| **IssuerURL** of the OIDC provider.                              | `issuerURL`    | https://accounts.google.com           |
| **ClientID** assigned by the OIDC provider.                      | `clientID`     | 1234567890.apps.googleusercontent.com |
| **ClientSecret** used to authenticate against the OIDC provider. | `clientSecret` | your-client-secret                    |

### Entry Quiz Configuration

The entry quiz is a simple quiz which is presented to users when they first visit the site. These parameters cannot be set via environment variables, and have no default values. They should be specified in the configruation file as a sequence of maps under the `entryQuiz` key.

| Parameter                             | YAML key   | Example value           |
| ------------------------------------- | ---------- | ----------------------- |
| **Question** to be asked to the user. | `question` | What is the best color? |
| **Answer** to the question.           | `answer`   | purple                  |

## Example Configuration

```yaml
port: 8080
address: 127.0.0.1

baseURL: http://localhost:8080

title: Epigram Demo
description: A place to record quotes.

repo: SQLite
DBLoc: ./epigram.db

trustProxy: false

OIDCProvider:
  name: google
  issuerURL: "https://accounts.google.com"
  clientId: "1234567890.apps.googleusercontent.com"
  clientSecret: "your-client-secret"

entryQuestions:
  - question: What is the best color?
    answer: purple
  - question: What is the best animal?
    answer: dog
```
