
# Testing

Learn how to test FlatTrack.


## Backend

The API uses [Ginkgo](https://onsi.github.io/ginkgo) as the e2e testing framework.


### e2e tests

    ginkgo -r --randomizeAllSpecs --randomizeSuites --failOnPending --cover --trace --progress test/backend/e2e


## Frontend

The frontend uses [Mocha](https://mochajs.org/) as the testing framework.


### e2e tests

    npm run test:e2e


### unit tests

    npm run test:unit

