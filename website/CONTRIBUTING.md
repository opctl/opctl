# Implementation details
The website is built using the [Docusaurus 2](https://v2.docusaurus.io/) static site generator.

# How do I...

## List operations
1. `opctl ls` from this directory will print a full operation manual.

## Dev
1. `opctl run dev` from this directory to init and start the website 
1. browse to [localhost:3000](http://localhost:3000) in your web browser
1. make some changes, your browser should live-reload them

## Deploy to https://opctl.io
1. `opctl run deploy` to deploy the website to github pages (hosts https://opctl.io)


# Contribution guidelines
- DO follow [![JavaScript Style Guide](https://img.shields.io/badge/code_style-standard-brightgreen.svg)](https://standardjs.com).
- DO use style objects and fallback to emotion css object generated classNames when not possible. 
- DO write tests in `arrange`, `act`, `assert` format w/ the given object under test referred to as `objectUnderTest`.
- DO keep tests alongside source code; i.e. place `code.test.ts` alongside `code.ts`.