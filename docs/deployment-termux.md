
# Termux on Android

> Running FlatTrack on Android with Termux


# Install Termux

Termux can be found in the [F-Droid store](https://f-droid.org/en/packages/com.termux/) or Google Play.


# Obtaining FlatTrack


## Getting a release

Head over to [GitLab pipelines](https://gitlab.com/flattrack/flattrack/-/pipelines) and find the latest tagged release.
The pipeline will contain a `build-zip` job which has artifacts to download.

e.g:

    wget https://gitlab.com/flattrack/flattrack/-/jobs/artifacts/0.16.1/download?job=build-zip


## Manual installation


### Install Dependencies

    pkg add postgresql nodejs golang git python2


### Clone the repo

    git clone https://gitlab.com/flattrack/flattrack


### Build the frontend

Go to the web folder

    cd web

Fetch dependencies:

    npm i

Build:

    npm run build


### Build the backend

    go build -o flattrack ./src/backend/main.go


# Initialize and start Postgres

To bring up FlatTrack on Termux, simply run the following command:

    ./hack/start-on-termux.sh

This will:

-   prepare, init, and launch Postgres
-   set default for and write a .env file
-   bring up FlatTrack

