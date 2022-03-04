# Project Structure

Epigram is based on a hexagonal architecture (also known as ports and adapters), which permits swapping both the database / storage infrastructure as well as the front end presentation layer with minimal changes to the codebase and preservation of the core application logic. 

To accomplish this, the application is broken down into three primary tiers:

- **`server`**
    - Provides an interface to the services for client interaction.
    - Currently implemented adapters:
        - `http` - uses Go templating to perform server side rendering and provide an interface to the application from a browser.
- **`service`**
    - Core application logic.
    - Shared across all server and storage implementation.
    - Each service requires a unique repository to maintain on it's data.
    - Some services interface with sub-services, and abstract their methods for use by the server.
- **`storage`**
    - Application state storage.
    - Comprised of multiple repositories, each one tasked with storing a specific, consistent type of data.
    - Currently implemented adapters:
        - `inmemory` - stores all data in maps in memory, useful for testing.
        - `sqlite` - uses statically-linked sqlite3 library to store data in local database file.

In addition, each layer depends on the `model` package, which holds all of the structure definitions for objects used throughout the above three tiers. 

The below diagram shows how a `server` implementation utilizes the top level `service` entities, and how each service interacts with others and their repositories in the `storage` package.

```mermaid
classDiagram
    %%class `model.User` {
    %%    +ID         string
    %%    +Name       string
    %%    +Email      string
    %%    +PictureURL string
    %%    +Created    time.Time
    %%    +QuizPassed   bool
    %%    +QuizAttempts int8
    %%    +Banned       bool
    %%    +Admin        bool
    %%    -isAuthorized() bool
    %%    -isAdmin() bool
    %%}

    class `service.User` {
        -ur UserRepository
        -sess service.UserSession
        +GetUserFromIDToken(ctx context.Context, token oidc.IDToken) (model.User, error)
        +CreateUser(ctx context.Context, u *model.User) error
        +FindUserById(ctx context.Context, id string) (model.User, error)
        +UpdateUser(ctx context.Context, u model.User) error
        +CreateUserSession(ctx context.Context, u model.User) (model.UserSession, error)
        +GetUserFromSessionID(ctx context.Context, sessID string) (model.User, error)
    }

    class `service.UserSession` {
        -repo UserSessionRepository
        +CreateUserSession(ctx context.Context, u model.User) (model.UserSession, error)
        +FindSessionByID(ctx context.Context, id string) (model.UserSession, error)
    }

    `service.User` --> `service.UserSession`

    class `service.EntryQuiz`{
        +Questions []QuizQuestion
        +VerifyAnswers(answers map[int]string) (passed bool)
    }
    
    class `server`{

    }

    `server` --> `service.User`
    `server` --> `service.EntryQuiz`

    `service.User` --> `UserRepository`

    class `UserRepository` {
        <<Interface>>
        +Create(ctx context.Context, u model.User) error
        +Update(ctx context.Context, u model.User) error
        +FindByID(ctx context.Context, id string) (model.User, error)
        +FindAll(ctx context.Context) ([]model.User, error)
    }

    `service.UserSession` --> `UserSessionRepository`

    class `UserSessionRepository` {
        <<Interface>>
        +Create(ctx context.Context, us model.UserSession) error
	    +FindByID(ctx context.Context, id string) (model.UserSession, error)
    }

    class `QuoteRepository` {
        <<Interface>>
        +Create(ctx context.Context, q model.Quote) error
        +Update(ctx context.Context, q model.Quote) error
        +FindByID(ctx context.Context, id string) (model.Quote, error)
        +FindAll(ctx context.Context) ([]model.Quote, error)
    }

    class `service.Quote` {
        -repo QuoteRepository
        +CreateQuote(ctx context.Context, q *model.Quote) error
        +GetAllQuotes(ctx context.Context) ([]model.Quote, error)
    }

    `server` --> `service.Quote`
    `service.Quote` --> `QuoteRepository`

    class `service.OIDC` {
        +Name string
        +IssuerURL string
        +ClientID     string
        +ClientSecret string
        -config   oauth2.Config
        -provider *oidc.Provider
        +RedirectURL(state string, nonce string) (url string)
        +CallbackURL() string
        +Init(baseURL string) error
        +ValidateCallback(r http.Request) (oidc.IDToken, error)
    }

    `server` --> `service.OIDC`
```