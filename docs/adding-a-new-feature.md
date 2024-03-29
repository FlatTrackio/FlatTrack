# Adding a new feature

> A flow for adding a new feature to FlatTrack

## Considering the feature

Make an issue in GitLab to disuss the feature and to plan how it should be.

## Database

Consider how the data will be structured.
Will new tables be needed? If so, using [golang-migrate](https://github.com/golang-migrate/migrate) to migrate the tables in place.
Checkout the structure for other tables. Common features are required (with exception), such as:

-   id (refer to the row)
-   author (who changed the data)
-   authorLast (who changed it before the current change)
-   creationTimestamp (when was it created)
-   modificationTimestamp (when was it modified)
-   deletionTimestamp (when was it deleted)

For a new change to the database with a new feature, migrations (up) must be merged into a single table.

## Backend

### Feature functions

In the `pkg` folder, `yourfeature/yourfeature.go` should be created with the feature functions.
FlatTrack is generally resource oriented, so what&rsquo;s put in is the same as what&rsquo;s returned (except passwords).
For a fully flexible feature, make sure that all of the regular HTTP verbs are covered with resource modifications:

-   create
-   get
-   update
-   patch
-   delete

Sometimes it can be useful to have options as a parameter to pass to the feature function (e.g: ShoppingListOptions, ShoppingTagOptions). These can be useful for things like filtering, sorting, and selecting.

### Handlers

A handler should be responsible for taking the HTTP request turning the data into something that the feature handler can use; Such as unmarshalling the JSON, getting the vars.

All I/O should be JSON based for consistency and metadata.

Handlers are stored in `pkg/routes/routes.go`

### Endpoints

Endpoints are where the hander is linked up to a route. The routes most likely will be restricted by the requirement of authentication. Some are also restricted by group, commonly `admin`.

Endpoints are stored in `pkg/routes/endpoints.go`

## Frontend

### Request handlers

Request handlers are the linking of the frontend to the backend.
The handlers are based on the HTTP verbs for managing the feature&rsquo;s resources.
They are responsible for formulating the request in the way that the backend understands.

Request handlers are located in `web/src/requests`, then in either `admin`, `authenticated`, and/or `public`.

### Page

The feature page will be (at least the first) page where the feature will live.

App pages are located in `web/src/views`, then in either `admin`, `authenticated`, and/or `public`.

### Apps page

The final piece is adding the listing to the apps page.

The general apps page is located in `web/src/view/authenticated/apps.vue`.
The admin apps page is located in `web/src/view/admin/apps.vue`

## Tests

New features require tests to be written for them.
Without tests, the features may not have the expectation of consistent behaviour.

Tests are located in `test/backend/e2e/api_e2e.go`

## Docs

If there have been changes that require updating of docs, such as how FlatTrack is configured or deployed (or requirements therein) new documentation is required.

## Expectations

### Consistency

New features that are brought in should be styled and written in a similar way.

This in mind, feel free to suggest better ways of doing things.

### Self-hosting

No portion should depend on external services (especially default).
Examples:

-   fonts, CSS are only located on the instance
-   authentication is only performed on the instance (proxy OAuth may be allowed)

