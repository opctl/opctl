---
title: UI
---

### Access the UI
1. Make sure you have a node running; an easy way is to run an op e.g. `opctl run github.com/opspec-pkgs/uuid.v4.generate#1.1.0`.
1. Open a web browser and navigate to [http://localhost:42224](http://localhost:42224).

### Visualize an op
1. Obtain either a local or remote ref to mount.
1. [Access the UI](#access-the-ui) with `?mount=ref`. Remote refs need to be URL encoded, Local refs need to be absolute paths.
1. From the explorer, navigate to the directory containing an `op.yml` and open it.

> If any remote private refs are encountered and aren't cached, you'll be prompted for pullCreds.

#### Example visualize github.com/opctl/opctl#0.1.34/.opspec/build
1. Navigate to [http://localhost:42224/?mount=github.com%2Fopctl%2Fopctl%230.1.34%2F.opspec%2Fbuild](http://localhost:42224/?mount=github.com%2Fopctl%2Fopctl%230.1.34%2F.opspec%2Fbuild)
1. expand `/` > `opctl` > `opctl#0.1.34` > `.opspec` > `build`
1. open `op.yml`
