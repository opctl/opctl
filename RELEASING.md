# Releasing
Opctl's GitHub Actions `build` workflow automatically does most of the work required to create a new release for the software. As a developer there are a few aspects of releasing new versions of opctl that you'll need to be aware of.

## Updating the CHANGELOG
In order for opctl to automatically create a new release with your changes, you'll need to update `CHANGELOG.md` with a new version and a description of your changes. The `CHANGELOG` is opctl's source of truth for versioning information. If you _don't_ update the `CHANGELOG` you'll still be able to merge your PRs, but a job will fail in your PR that checks whether you've updated the file. This job is intended to be an indicator that you've forgotten a step in the update workflow.

## Promoting a release
As of this writing, when a change is merged to main that includes a new version definition in `CHANGELOG.md`, a new **draft** release will automatically be created. The following are your responsibilties in order to promote a release out of draft state.
1. Install the new version of opctl on your machine and perform some level of manual testing by running ops in various projects using the new binary
1. Checkout the main branch and pull the latest code
1. Run the `release/promote` op in order to update the existing release and remove its draft designation and also to release a docker image tagged with the new version to docker hub
