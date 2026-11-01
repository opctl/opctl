# react-ops-example
example ops alongside a simple react web app

This project was bootstrapped with [Create React App](https://github.com/facebookincubator/create-react-app).

You can find the most recent version of th development guide [here](https://github.com/facebookincubator/create-react-app/blob/master/packages/react-scripts/template/README.md).

Run the UI and Go service together from this directory:

```sh
opctl run -a HOME="$HOME" .opspec/dev
```

`HOME` provides the Go service's module and build caches.
