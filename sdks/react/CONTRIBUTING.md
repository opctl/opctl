# Implementation details
The react SDK is maintained in typescript and released to NPM as javascript (plus type definitions).

# How do I...

## List operations
`opctl ls` from this directory will print a full operation manual.

## Run tests
`opctl run test` from this directory to init and test the SDK.

## Release to NPM
`opctl run release` from this directory to init, test, compile, and publish the SDK to NPM.


# Contribution guidelines
- DO follow [![JavaScript Style Guide](https://img.shields.io/badge/code_style-standard-brightgreen.svg)](https://standardjs.com).
- DO use style objects and fallback to emotion css object generated classNames when not possible. 
- DO write tests in `arrange`, `act`, `assert` format w/ the given object under test referred to as `objectUnderTest`.
- DO keep tests alongside source code; i.e. place `code.test.ts` alongside `code.ts`.