## problem statement

lib for pub-sub style messaging

## features

- API exposed via interface
- retroactive subscription
- fake implementation to allow faking interactions

## event stores

Three event stores have been implemented in search of performance gains.
As the performance of each may change over time, they're kept here for easy re-evaluation.

The results of reading/writing 1M events (conducted on the same machine) were as follows:


| implementation | Write 1M events | Read 1M events |
|--|--|--|
|[badgerDBEventStore](badgerDBEventStore.go)|05:53|00:15|
|[boltDBEventStore](boltDBEventStore.go)|04:09|00:28|
|[buntDBEventStore](buntDBEventStore.go)|01:32|00:45|

