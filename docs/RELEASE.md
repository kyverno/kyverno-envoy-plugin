# Release docs

This doc contains information for releasing a new version.

## Create a release

Creating a release can be done by pushing a tag to the GitHub repository (beginning with `v`).

The [release workflow](../../.github/workflows/release.yaml) will take care of creating the GitHub release and will publish artifacts.

```shell
VERSION="v0.0.1"
TAG=$VERSION

git tag $TAG -m "tag $TAG" -a
git push origin $TAG
```

## Publish documentation

Publishing the documentation for a release is decoupled from cutting a release.

To publish the documentation push a tag to the GitHub repository (beginning with `docs-v`).

```shell
VERSION="v0.0.1"
TAG=docs-$VERSION

git tag $TAG -m "tag $TAG" -a
git push origin $TAG
```

## Misc

- Add to the drop-down list in the bug template
