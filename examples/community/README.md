# How to contribute a community example

Currently we are accepting PRs for community contributed shader examples. If lots of people contribute, we will find a better way to do this, but let's keep it simple for now.

There are two things you will need to do:
- The first is creating a folder inside `/examples/community/` with at least a `go.mod`, `go.sum`, a `main.go` and a `shader.kage` file. This is similar to the example programs shown on the intro tutorial. For an actual example, see the [examples/intro/mirror](https://github.com/tinne26/kage-desk/blob/main/examples/intro/mirror) folder. Minimal dependencies where possible.
- The second is modifying the `community_examples.md` file to add a line with the link to your example. The markdown may look like this:
```
- [Shader link](https://github.com/tinne26/kage-desk/blob/main/examples/community/subfolder) by [Author](https://github.com/username), v2.X: short description.
```
Replace the `X` in `v2.X` with whatever version of Ebitengine you used for the example.
