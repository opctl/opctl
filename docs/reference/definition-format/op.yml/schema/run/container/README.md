Runs a container.

## Properties

* [dirs](dirs/README.md)
* [envVars](envVars/README.md)
* [files](files/README.md)
* [image](image/README.md)
* [name](name.md)
* [sockets](sockets/README.md)
* [workDir](workDir.md)

# Examples

## Kitchen sink

```yaml
name: kitchenSink
inputs:
  mySocket: 
    string: {}
  registryCreds:
    object:
      constraints:
        properties:
          username: { type: string }
          password: { type: string }
        required: [username, password]
run:
  container:
    dirs:
      /:
    envVars:
      MY_ENV_VAR: my value
    files:
      /op.yml:
      /hello.txt: hello!
    image:
      ref: customBase:1.0.0
      pullCreds:
        username: $(registryCreds.username)
        password: $(registryCreds.password)
    ports: { '80': '80' }
    sockets:
      /mySocket: $(mySocket)
    name: my-kitchen-sink
    workDir: /root
```
