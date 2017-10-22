## Intro

In this first tutorial of the "creating a pkg" series, we'll create a
simple pkg which uses a container to print 'hello world!'

## Steps

### step 1

create a dir named `hello-world`

### step 2

add a file named `op.yml` to `hello-world` w/ contents:
[include](op.yml)

#### analysis

##### name

`name` defines the name of the pkg.
> `name` needs to be a valid
> [uri-reference](https://tools.ietf.org/html/rfc3986#section-4.1) but
> doesn't need to be resolvable

##### description

`description` defines the description of the pkg
> see [multiline yaml strings](http://yaml-multiline.info/) if you need
> a multiline description

##### run

`run` defines what the operation runs a.k.a the "run graph".

In this case we have defined the run graph to be a single container
call.

The container call will:

- resolve the [alpine](https://hub.docker.com/_/alpine/) container image
- invoke the echo binary w/ arg `hello world!`

## Outro

Congrats! you've created your first opspec package.
