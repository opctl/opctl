---
title: Run a go service
sidebar_label: run a go service
---

We'll now look at an op to run a sample Go application. Please see the code at [https://github.com/opctl/opctl/tree/master/examples/run-a-go-service](https://github.com/opctl/opctl/tree/master/examples/run-a-go-service)
The sample application we have requires a mysql database, therefore to run it locally we need to:
1. Start a MySQL database with some DB credentials
2. (optional) Seed the DB with sample data
3. Start the application and provide the MySQL DB credentials as inputs - our application will read those from environment variables

The ops to make that happen are explained below.

#### mysql
```yaml
name: mysql
description: runs a mysql container, optionally seeding it with sample data
inputs:
  MYSQL_USER:
    string:
      default: testuser
      description: username for MySQL user to create
  MYSQL_PASSWORD:
    string:
      default: testpassword
      description: password for MySQL user to create
  MYSQL_DATABASE:
    string:
      default: testapp
      description: Database to create
  MYSQL_HOST:
    string:
      default: mysql
  doSeed:
    boolean:
      default: false
      description: if true, sample data will be inserted in the database
run:
  parallel:
    - container:
        image: { ref: 'mysql:8' }
        name: mysql # this value will be provided as input to the 'local' op so our code knows what mysql host to connect to
        envVars: {MYSQL_USER , MYSQL_PASSWORD , MYSQL_DATABASE , MYSQL_RANDOM_ROOT_PASSWORD: "yes"}
        ports: { '3306': '3306' }
    - container:
        image: { ref: 'mysql:8' }
        envVars: {MYSQL_USER, MYSQL_PASSWORD, MYSQL_HOST, MYSQL_DATABASE, doSeed}
        files:
          /seed.sql: # this shorthand form will copy seed.sql from the root of this op to the root of the container
        workDir: /db
        cmd: 
          - sh
          - -ce
          - |
            echo "starting seed script"
            sleep 25 # we'll sleep while the MySQL DB starts
            if [ $doSeed = "true" ]
            then
              echo "connecting to load sql script"
              mysql -u $MYSQL_USER -p$MYSQL_PASSWORD -h $MYSQL_HOST $MYSQL_DATABASE < /seed.sql
            fi
            exit 0
```
This op will run a MySQL database in a Docker container. It has an optional input, `doSeed`, which is of type `Boolean`. When `true`, it will run the `seed.sql` script in a second, parallel container to seed the MySQL database in first container with a single table and a single row of data.

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
    string:
      default: testuser
      description: username for MySQL user to create
  MYSQL_PASSWORD:
    string:
      default: testpassword
      description: password for MySQL user to create
  MYSQL_DATABASE:
    string:
      default: testapp
      description: Database to create
  MYSQL_HOST:
    string:
      default: mysql # the host for mysql is the container name in the mysql op
run:
  parallel:
    - op:
        ref: ../mysql # we reference the mysql op using its relative path to this op
        inputs: { MYSQL_USER, MYSQL_PASSWORD, MYSQL_HOST, MYSQL_DATABASE, doSeed: true} # we pass the relevant inputs through to the mysql op
    - container:
        image: { ref: 'golang:1.10.3' }
        name: go-svc
        dirs:
          /go/src/github.com/golang-ops-example: $(/app) # IMPORTANT: we've created a symlink in the root of the op (i.e. at /.opspec/dev) to the source code of our app (i.e. at /app) - this is so we can encapsulate the op, meaning that any files the op needs to reference now are accessible from within its directory. we reference that symlink here as /app because that refers to the root of the op - see https://opctl.io/docs/reference/op-definition-format/op.yml/expressions/dir-entry-ref.html#embedded-json-file for more details
        envVars: {MYSQL_USER, MYSQL_PASSWORD, MYSQL_HOST, MYSQL_DATABASE} # the same inputs are needed to let our code know how to connect to the DB
        workDir: /go/src/github.com/golang-ops-example
        ports: { '8080': '8080' }
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