The data directory is the location where opctl's state is stored such as remote ops, events, run graph data (e.g. container files/dirs), etc...

### Default location
By default, the data dir is located at "[per user app data path](https://github.com/appdataspec/spec/blob/master/index.md#per-user-app-data)/opctl"

### Custom location (since v0.1.25)
A custom data dir location (such as `.opctl`) can be specified by (listed in order of precedence):
- passing the `--data-dir` option to [node create](../cli/node/create.md)
- setting an `OPCTL_DATA_DIR` environment variable

The value of either MUST be a relative or absolute path. 

## Paths

* [pid.lock](pid.lock.md)
* [ops](ops/README.md)
* [dcg](dcg/README.md)
  * [event.db](dcg/event.db.md)
  * [{op-id}](dcg/op-id/README.md)
    * [containers](dcg/op-id/containers/README.md)
      * [{container-id}](dcg/op-id/containers/container-id/README.md)
        * [fs](dcg/op-id/containers/container-id/fs/README.md)
