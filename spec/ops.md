## Ops

An op is an orchestration of (a) containerized process(es).

Ops SHOULD be deterministic i.e, assuming external state remains
constant, calling the same package w/ the same input(s) SHOULD always
result in the same output(s).
