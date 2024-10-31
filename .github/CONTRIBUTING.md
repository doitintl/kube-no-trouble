# Contributing Guidelines

*Pull requests, bug reports, and all other forms of contribution are welcomed and highly encouraged!*

### Contents

- [Opening an Issue](#opening-an-issue)
- [Feature Requests](#feature-requests)
- [Triaging Issues](#triaging-issues)
- [Submitting Pull Requests](#repeat-submitting-pull-requests)
- [Branches](#branches)
- [Writing Commit Messages](#writing-commit-messages)
- [Signing Commits](#signing-commits)
- [Code Review](#mark-code-review)
- [Coding Style](#coding-style)
- [Changelog](#changelog)
- [Release](#releases)
- [Credits](#credits)

> **This guide serves to set clear expectations for everyone involved with the project. We are indebted to @stepanstipl for creating this project.**

## Opening an Issue

Before [creating an issue](https://help.github.com/en/github/managing-your-work-on-github/creating-an-issue), check if you are using the latest version of the project. If you are not up-to-date, see if updating fixes your issue first.

### Reporting Security Issues

File a public issue for security vulnerabilities. This way maintainers can quickly act to add a fix.

### Bug Reports and Other Issues

A great way to contribute to the project is to send a detailed issue when you encounter a problem.

- **Do not open a duplicate issue!** Search through existing issues to see if your issue has previously been reported. If your issue exists, comment with any additional information you have. You may simply note "I have this problem too", which helps prioritize the most common problems and requests.

- **Use [GitHub-flavored Markdown](https://help.github.com/en/github/writing-on-github/basic-writing-and-formatting-syntax).** Especially put code blocks and console outputs in backticks (```). This improves readability.

## Feature Requests

Feature requests are welcome! While we will consider all requests, we cannot guarantee your request will be accepted. We want to avoid [feature creep](https://en.wikipedia.org/wiki/Feature_creep). Your idea may be great, but also out-of-scope for the project. If accepted, we cannot make any commitments regarding the timeline for implementation and release. However, you are welcome to submit a pull request to help!

- **Do not open a duplicate feature request.** Search for existing feature requests first. If you find your feature (or one very similar) previously requested, comment on that issue.

- Be precise about the proposed outcome of the feature and how it relates to existing features. Include implementation details if possible.

## Submitting Pull Requests

Before [forking the repo](https://help.github.com/en/github/getting-started-with-github/fork-a-repo) and [creating a pull request](https://help.github.com/en/github/collaborating-with-issues-and-pull-requests/proposing-changes-to-your-work-with-pull-requests) for non-trivial changes, it is usually best to first open an issue to discuss the changes, or discuss your intended approach for solving the problem in the comments for an existing issue.

*Note: All contributions will be licensed under the project's license.*

- **Smaller is better.** Submit **one** pull request per bug fix or feature. A pull request should contain isolated changes pertaining to a single bug fix or feature implementation. **Do not** refactor or reformat code that is unrelated to your change. It is better to **submit many small pull requests** rather than a single large one. Enormous pull requests will take enormous amounts of time to review, or may be rejected altogether.

- **Coordinate bigger changes.** For large and non-trivial changes, open an issue to discuss a strategy with the maintainers. Otherwise, you risk doing a lot of work for nothing!

- **Prioritize understanding over cleverness.** Write code clearly and concisely. Remember that source code usually gets written once and read often. Ensure the code is clear to the reader. The purpose and logic should be obvious to a reasonably skilled developer, otherwise you should add a comment that explains it.

- **Follow existing coding style and conventions.** Keep your code consistent with the style, formatting, and conventions in the rest of the code base. When possible, these will be enforced with a linter. Consistency makes it easier to review and modify in the future.

- **Include test coverage.** Add unit tests or UI tests when possible. Follow existing patterns for implementing tests.

- **Add documentation.** Document your changes with code doc comments or in existing guides.

- **Use the repo's default branch.** Branch from and [submit your pull request](https://help.github.com/en/github/collaborating-with-issues-and-pull-requests/creating-a-pull-request-from-a-fork) to the repo's default branch `main`.

- **[Resolve any merge conflicts](https://help.github.com/en/github/collaborating-with-issues-and-pull-requests/resolving-a-merge-conflict-on-github)** that occur.

- **Promptly address any CI failures**. If your pull request fails to build or pass tests, please ideally amend your current commit with a `git commit --amend` or create a new commit. Your welcome to squash them together.

## Branches

All PRs should be at the HEAD, that is to say we use merge commits with [semi-linear history](https://docs.gitlab.com/ee/user/project/merge_requests/methods/#merge-commit-with-semi-linear-history).

![](https://devblogs.microsoft.com/devops/wp-content/uploads/sites/6/2019/04/semilinear-1.gif)

If not they should be rebased on master, do not merge master into your branches. Do not create branches that already exist, please use a new name. Branches are deleted upon merging. Thank you.

```sh
git fetch -v
git checkout my-branch
git rebase origin/master
```

We also maintain an author merges policy. We will approve your pr then you can merge it.

Additionally we ask that you do not squash your prs or introduce a large number of commits per feature. If your pr becomes too large we may ask you to split the functionality out into more than one pr.

## Writing Commit Messages

We enforce [Conventional Commits][cc] for all commits in the form:

```
<type>: <summary>

[optional body]

[optional footer(s)]
```

Where type is one of:
- **build** - Affects build and/or build system
- **chore** - Other non-functional changes
- **ci** - Affects CI (e.g. GitHub actions)
- **dep** - Dependency update
- **docs** - Documentation only change
- **feat** - A new feature
- **fix** - A bug fix
- **ref** - Code refactoring without functionality change
- **style** - Formatting changes
- **test** - Adding/changing tests

[cc]: https://www.conventionalcommits.org/

Use imperative, present tense (Add, not ~Added~), capitalize first letter of summary, no dot at the and. The body and footer are optional. And are ignored by our release note process.

Relevant GitHub issues should be referenced in the footer in the form `fix(readme): Fixes #456`.

## Signing Commits

Please ensure any contributions are signed with a valid gpg key. We use this to validate that you have committed this and no one else. You can learn how to create a GPG key [here](https://docs.github.com/en/authentication/managing-commit-signature-verification/generating-a-new-gpg-key).

## Code Review

- **Review the code, not the author.** Look for and suggest improvements without disparaging or insulting the author. Provide actionable feedback and explain your reasoning.

- **You are not your code.** When your code is critiqued, questioned, or constructively criticized, remember that you are not your code. Do not take code review personally.

- **Take your time.** Please do not rush to generate functionality, we hold the code in this repo to a high standard taking your time to ensure a high quality result is appreciated.

- Kindly note any violations to the guidelines specified in this document.

## Coding Style

Consistency is the most important. Following the existing style, formatting, and naming conventions of the file you are modifying and of the overall project. Failure to do so will result in a prolonged review process that has to focus on updating the superficial aspects of your code, rather than improving its functionality and performance.

For example, if all private properties are prefixed with an underscore `_`, then new ones you add should be prefixed in the same way. Or, if methods are named using camelcase, like `thisIsMyNewMethod`, then do not diverge from that by writing `this_is_my_new_method`. You get the idea. If in doubt, please ask or search the codebase for something similar.

We use [pre-commit](https://pre-commit.com/) to validate the style of the code along with maintainer review. If you'd like to check that your code matches our style please run and/or install our [pre-commit](https://github.com/doitintl/kube-no-trouble/blob/master/.pre-commit-config.yaml). Branches with commits which do not pass the pre-commit will not be accepted.

```
pip install pre-commit
pre-commit run --all-files
```

### Changelog

Changelog is generated automatically based on merged PRs using
[git-cliff](https://git-cliff.org/). Template can be found in [cliff.toml](https://github.com/doitintl/kube-no-trouble/blob/master/cliff.toml).

PRs are categorized based on their conventional commit groups, into following sections, as seen in [git cliff toml file line 62](https://github.com/doitintl/kube-no-trouble/blob/master/cliff.toml#L62):
- Features - group **feat** - A new feature
- Fixes - group **fix** - A bug fix
- Internal/Other - groups **chore** **build** **ci** **build** **dep**  **docs** **ref** **style** **test** - all other changes

Additionally we will reference any new contributors between the release versions. See an example release note below:

```md
#### Docker Image: ghcr.io/doitintl/kube-no-trouble:latest

## Changelog

### Features:

    feat: Add rego for v1.32 deprecations b4da33a by @dark0dave
    feat: Fix github actions for creating release notes edd2dc3 by @dark0dave

### Fixes:

    fix: Add docker image back 6de101d by @dark0dave
    fix: Add fix for git cliff 9d487c9 by @dark0dave

### Internal/Other:

    dep: Bump lots of deps 7cdf86a by @dark0dave

#### Full Changelog: 1.0.0...2.0.0
### New Contributors

    @dark0dave made their first contribution in #1
```

## Releases

We use [semantic versioning](https://semver.org/) in kubent.

> Given a version number MAJOR.MINOR.PATCH, increment the:
>
> - MAJOR version when you make incompatible API changes
> - MINOR version when you add functionality in a backward compatible manner
> - PATCH version when you make backward compatible bug fixes
>
> Additional labels for pre-release and build metadata are available as extensions to the MAJOR.MINOR.PATCH format.

### Triggering a release

Releases are triggered by pushing a new tag, it is **imperative** you keep to the style of tag defined in the cliff.toml file otherwise git-cliff will not generate valid release notes.

### Expected tag structure

From the cliff toml file:
```toml
tag_pattern = "^[0-9]+.[0-9]+.[0-9]+$"
```
Example:
- 0.7.1
- 0.7.2
- 0.7.3

#### Nightly

Nightly release happen on weekly schedule, this important for good health of the repo. This insures all our ci is always running properly. The intent here, is to release versions so that developers can test recent functionality in the wild. We can not foresee all out comes so we must rely on others.

#### When to release

We create [milestones](https://github.com/doitintl/kube-no-trouble/milestones) in the github repo which guide us towards minor releases, such as 0.7.0 or 0.6.0. @stepanstipl the creator of this repo has always avoided releasing 1.0.0 because that would lock our existing functionality in. Currently there are large outstanding issue which prevent us from having a stable API, ie version 1.0.0. This is in line with the semver doc see [here](https://semver.org/#how-do-i-know-when-to-release-100).

All that to say, our patches tend to be a mixture of small features and bugfixes. Example:

```md
#### Features:
    feat: Add rego for v1.32 deprecations b4da33a by @dark0dave
    feat: Fix github actions for creating release notes edd2dc3 by @dark0dave
#### Fixes:
    fix: Script install.sh in dumb TERM 35927e8 by @FabioAntunes
    fix: warn and fix invalid namespace 83308e1 by @justdan96
    fix: Add docker image back 6de101d by @dark0dave
    fix: Add fix for git cliff 9d487c9 by @dark0dave
#### Internal/Other:
    dep: Bump lots of deps 7cdf86a by @dark0dave
```

Please make an effort to stay with in these guidelines. A **smaller number** releases is preferred to a large number of releases. Additionally, We DO NOT SUPPORT `hotfixing` releases, see [hotfix explained](https://en.wikipedia.org/wiki/Hotfix). Releases stay static and immutable once released. In our view hotfixes cause more problems than they solve and introduce a burden on maintainers attempting to catch bugs and resolve issues.

## Credits

Very heavenly influenced by https://github.com/jessesquires/.github, thank you so much for providing a baseline for this.

Modified and updated by [@dark0dave](https://github.com/dark0dave). Adopted by all maintainers.
