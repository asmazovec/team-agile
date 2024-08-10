# Tools as dependencies

Based on:
- https://play-with-go.dev/tools-as-dependencies_go119_en;
- https://mariocarrion.com/2021/10/15/learning-golang-versioning-tools-as-dependencies.html.

The following module establish tools as dependencies for current solution.
It lets tools keep their versions over a time for deterministic and stable building and deploying current solution.

## Adding new tools as dependencies

To add new tool as dependency:
1. Append package link into [tools.go](./tools.go).
2. Run `make install/tools` to perform tools installing.
 