# Branches
Working with branches is done by rebase-flow.
## `main` - the main release branch with stable code ready for release.
- Injection into the branch only via PullRequest.
- Injection into the branch only by the product owner.
- The final set of changes introduced by this branch must meet the established requirements to code quality.

## `develop` - the main development branch with the most actual product code.
- Merges into the branch only via PullRequest.
- Only product owners and trusted developers can join the branch, it is forbidden to use groups for granting rights.
- This branch is only injected into the main and enhance minor version of the product.
- The final set of changes brought by this branch must meet the specified code quality requirements.

## `team#<issue>-optional-comment` - a feature branch for a task with priority from `P0` to `P2`.
- It is created by the developer and comes from the last commit of the develop branch.
- The optional comment must be in English, spaces are replaced by dashes.
- It is merged into the develop branch by rebase.
- The final set of changes introduced by this branch must meet the specified code quality requirements.

## `hotfix#<issue>-optional-commit` - hotfix branch for a task with `hotfix` priority
- Created by the developer and comes from the last commit of the main branch.
- The optional-comment must be in English, spaces are replaced by dashes.
- This branch is merged into the main branch via rebase and increases the patch version of the product - into the develop branch via merge commit.
- The final set of changes introduced by this branch must meet the established requirements for code quality.

# Commits
A commit must be an atomic set of changes, preferably conforming to all code quality conventions established in the project.
Commit comments should follow the pattern:

```
(#<issue-number>) Comment in English with a capital letter
```
