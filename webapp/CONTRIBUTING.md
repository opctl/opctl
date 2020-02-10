# Implementation details
The webapp is built using [react](https://facebook.github.io/react/) via [create react app](https://github.com/facebookincubator/create-react-app).


# How do I...

## List operations
1. `opctl ls` from this directory will print a full operation manual.

## Dev
1. `opctl run dev` from this directory to init and start the webapp with live-reload 
1. browse to [localhost](http://localhost) in your web browser
1. make some changes, your browser should live-reload them

## Run a build
1. `opctl run build` from this directory to init, test, and compile the webapp for deployment.


# Contribution guidelines
- DO follow [![JavaScript Style Guide](https://img.shields.io/badge/code_style-standard-brightgreen.svg)](https://standardjs.com).
- DO use style objects and fallback to emotion css object generated classNames when not possible. 
- DO write tests in `arrange`, `act`, `assert` format w/ the given object under test referred to as `objectUnderTest`.
- DO keep tests alongside source code; i.e. place `code.test.ts` alongside `code.ts`.
- DO use functional components & react hooks instead of class components & HOC's (Higher Order Components).