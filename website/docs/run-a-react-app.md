---
title: Run a react app
sidebar_label: run a react app
---

We'll now look at an op to run a sample React application. Please see the code at [https://github.com/opctl/opctl/tree/main/examples/run-a-react-app](https://github.com/opctl/opctl/tree/main/examples/run-a-react-app)

We'll stick to simple node.js conventions of including the run command in an `npm start` script in `package.json`. Because we used `create-react-app` to bootstrap our project, the start script is `react-scripts start`, which will launch the webpack dev server to serve our app for development.

Let's say our frontend React app needs to call the a go api from our previous example. When running that stack we want our services running locally (React app running via webpack, our go service, and the mysql database). We however prefer to have the ops running each service to live with the source code it runs, rather than in a separate place.

What we need our `dev` op to do then is to:
1. call `go-svc`'s `dev` op by remote reference
3. run our React app in a container

So our ops in `run-a-react-app` would look like this

#### dev
```yaml
name: dev
description: runs react-app for development
run:
  parallel:
    - op:
        # reference run-a-go-service dev op
        ref: ../../../run-a-go-service/.opspec/dev
    - container:
        image:
          ref: 'node:15-alpine'
        cmd:
          - sh
          - -ce
          - yarn && yarn run start
        dirs:
          /src: $(../..)
        workDir: /src
        ports:
          3000: 3000
```

going to http://localhost:8080 should show us the `go-svc` api being served, and http://localhost:3000 should show us the react app, which in turns is making a call to `go-svc` and fetching data to show.

Notice the following:
1. we're referencing the `dev` op from the go-svc.
2. we can run our ops in any combination of `parallel` and `serial` blocks, composing them as needed. for our case the `dev` op can run in the background while we init then run our react app
3. networking between our services "just works" by referencing containers by name, thanks to how they all are in the same Docker container network, so the webpack dev server proxy configuration in `package.json` targets `go-svc`:

``` json
"proxy": {
    "/api": {
      "target": "http://go-svc:8080"
    }
  }
```
