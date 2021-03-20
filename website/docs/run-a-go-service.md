---
title: Run a go service
sidebar_label: run a go service
---

We'll now look at an op to run a sample Go application. Please see the code at [https://github.com/opctl/opctl/tree/main/examples/run-a-go-service](https://github.com/opctl/opctl/tree/main/examples/run-a-go-service)
The sample application we have requires a mysql database, therefore to run it locally we need to:
1. Start a MySQL database with some DB credentials
2. Seed the DB with sample data
3. Start the application and provide the MySQL DB credentials as inputs - our application will read those from environment variables

The ops to make that happen are explained below.

#### mysql
```yaml
name: mysql
description: runs a mysql container, seeding it with sample data
inputs:
  MYSQL_USER:
    description: username for MySQL user to create
    string:
      default: testuser
  MYSQL_PASSWORD:
    description: password for MySQL user to create
    string:
      default: testpassword
  MYSQL_DATABASE:
    description: name of mysql database to create
    string:
      default: testapp
  MYSQL_HOST:
    string:
      default: run-a-go-service-mysql
run:
  container:
    dirs:
      # mount our data seed scripts
      /docker-entrypoint-initdb.d:
    image:
      ref: 'mysql:8'
    # this sets an opctl overlay network DNS A record so the containers available via this name.
    name: run-a-go-service-mysql
    envVars:
      MYSQL_USER:
      MYSQL_PASSWORD:
      MYSQL_DATABASE:
      MYSQL_RANDOM_ROOT_PASSWORD: "yes"
    ports:
      3306: 3306
```
This op will run a MySQL database in a Docker container and seeds the MySQL database with a single table and a single row of data.

Note that instead of writing the `cmd` inline, we could have also put the run script in a file in the op directory and referenced it.
e.g.
```yaml 
cmd: ["/run.sh"]
```
This becomes useful for readability if the script we're running in a container is large.

Note that while we could have also built a custom docker image that contains that script and ran it, that's not a recommended practice. While the op remains portable with that approach, it is less transparent and harder to maintain, since we would have the extra step of building and pushing that image whenever we make a change to the script. Instead, we usually opt for the leanest Docker image that fulfills our use case, and include any scripts or binaries we need in the op directory.

#### dev
```yaml
name: dev
description: runs go-svc for development
inputs:
  MYSQL_USER:
    description: username for MySQL user to create
    string:
      default: testuser
  MYSQL_PASSWORD:
    description: password for MySQL user to create
    string:
      default: testpassword
  MYSQL_DATABASE:
    description: Database to create
    string:
      default: testapp
  MYSQL_HOST:
    string:
      default: mysql # the host for mysql is the container name in the mysql op
run:
  parallel:
    - op:
        ref: ../mysql # we reference the mysql op using its relative path to this op
        inputs:
          # we pass the relevant inputs through to the mysql op
          MYSQL_USER:
          MYSQL_PASSWORD:
          MYSQL_HOST:
          MYSQL_DATABASE:
    - container:
        image:
          ref: 'golang:1.15'
        name: go-svc
        dirs:
          # mount the source code of our app to the container
          /src: $(../..) 
        envVars:
          # let our code know how to connect to the DB
          MYSQL_USER:
          MYSQL_PASSWORD:
          MYSQL_HOST:
          MYSQL_DATABASE:
        workDir: /src
        ports:
          8080: 8080
        cmd:
          - sh
          - -ce
          - |
            go get -u github.com/cespare/reflex # get reflex to watch and hot-reload our main.go
            sleep 30 # we'll sleep while the MySQL DB starts
            /go/bin/reflex -g 'main.go' -s -- sh -c 'go run main.go'
```

This op will do the following, in parallel:
1. Run the mysql op, passing in the inputs needed to create it
2. run our service in the `golang` container. We'll use [reflex](https://github.com/cespare/reflex) to make development easier and avoid having to restart the whole operation after every change to main.go

Now to run the service locally, we only run `opctl run dev`

Notice how we can run ops in a `parallel` block. Our app runs parallel to the MySQL DB container (via the `mysql` op), because we need mysql running while our app runs.
