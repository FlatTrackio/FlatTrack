- [API](#sec-1)
    - [e2e tests](#sec-1-0-1)
- [Frontend](#sec-2)
    - [e2e tests](#sec-2-0-1)
    - [unit tests](#sec-2-0-2)

Learn how to test FlatTrack.

# API<a id="sec-1"></a>

The API uses [Ginkgo](https://onsi.github.io/ginkgo) as the e2e testing framework.

### e2e tests<a id="sec-1-0-1"></a>

```sh
ginkgo -r --randomizeAllSpecs --randomizeSuites --failOnPending --cover --trace --progress test/backend/e2e
```

# Frontend<a id="sec-2"></a>

The frontend uses [Mocha](https://mochajs.org/) as the testing framework.

### e2e tests<a id="sec-2-0-1"></a>

```sh
npm run test:e2e
```

### unit tests<a id="sec-2-0-2"></a>

```sh
npm run test:unit
```
