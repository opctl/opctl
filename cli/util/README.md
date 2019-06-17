util/ contains utility packages used by the opctl/opctl project.

External consumers should not depend on utility packages as they may
have breaking changes or cease to exist without notice.

Utility packages are kept separate from the core codebase to keep it as
small and concise as possible. If some utilities grow larger and their
APIs stabilize, they may be moved to their own repository under the
opspec organization to facilitate re-use by other projects however that
is not the priority.
