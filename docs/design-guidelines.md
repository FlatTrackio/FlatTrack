# Design guidelines

How FlatTrack must be built.


## Project

The project must be:

-   community driven: all efforts are for the community and be the community
-   freely\* available and freely\* licensed: code and artwork
-   collaborative: all apps in FlatTrack are for collaboration and should grow the interactions between those who use them
-   useful: the features provided should be exactly what people need
-   reliable: it must be built for strong-reliablity
-   secure: it must be secure by default; there must be large security focus
-   private: it must not communicate with third-parties for app resources and functions (exception for optional things like authentication via OAuth); information in the attached database must not leave where it resides
-   accessible: there must be many ways to run it or get an instance
-   portable: it must be easy to pick up an instance and take it somewhere else
-   conformant: no matter where it runs, it should be expected to be the same
-   tested: it should be well-tested and audited for quality and consistency


## API

The API must be:

-   stateless: it doesn't have any moving parts; nothing is stored in memory which is unique to the instance - so it can scale
-   dependency: it should live out of the database
-   performant: it must be built to be fast
-   declarative: all resources (accounts, shopping lists, etc...) must read and write mostly the same data (expect credentials and secrets)
-   structured: it must respond and accept JSON data
-   split up into packages: each feature or area of FlatTrack must be split up in package for easy reuse and testing


## Frontend

The frontend must be:

-   reactive: the UI adjusts to fit on various screen sizes
-   highly configurable: it should enable advanced settings
-   simple, easy to use, and intuitive: it should be obvious how each action should be done
-   dependency on the API: it should only talk the API from the same source and nothing else; all requests reference the same origin


## Packaging

The packaging must be:

-   minimal: contain only the app and nothing more (expect various stateless configuration files)


## Deployment

The deployments, in production, should be:

-   highly available: downtime should not be expected


## Artwork

The artwork must be:

-   friendly and inviting: it must be a delightful experience

---

\*free as in freedom

