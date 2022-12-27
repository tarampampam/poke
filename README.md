<div align="center">
  <img src="https://user-images.githubusercontent.com/7326800/203780071-a963f064-e8bd-4d0c-bbf3-37bdbde03547.png" alt="Logo" width=300" />
  <p>A <strong>simple and powerful</strong> tool for <strong>E2E-testing</strong>.</p>

![Project language][badge_language]
[![Build Status][badge_build]][link_build]
[![Coverage][badge_coverage]][link_coverage]
[![Image size][badge_size_latest]][link_docker_hub]

  Just imagine, that you can write E2E tests using familiar javascript syntax without any external dependencies,
  like node.js, jest, Axios, etc. Sounds like something impossible? 😉
</div>

> 🚧 The project currently is under active development. Do not use it somewhere except the testing purposes.

## 🔥 Features list

- Write your tests in familiar ECMAScript 5.1 syntax - `let`, `const`, arrow functions and classes are supported ([JS VM](https://github.com/dop251/goja) embedded in the app)
- Tests are **parallelized** by running them in their own thread to minimize execution time
- Ready to work out of the box - you won't need to configure anything
- Like any Go application, you can run Poke on Linux, windows, or macOS. Additionally, it ready to integrate with **CI/CD** and **Docker**
- The jest-like syntax (`describe`, `test`, `beforeEach`) for grouping tests cases and asserts
- You can use **fake data generator** (strings, numbers, IP addresses, UUIDs, and so on), **hashing**, and **encoding** functions - they are embedded in the app
- Easy to use by developers and QA engineers

<h2 align="center"><a href="https://poke.is-an.app/">📖 Read full documentation at our site</a></h2>

## 🗒 TODO

- [ ] `require(<js-or-json-file>)`
- [ ] `Language reference generation`

## Support

[![Issues][badge_issues]][link_issues]
[![Issues][badge_pulls]][link_pulls]

If you find any package errors, please, [make an issue][link_create_issue] in current repository.

## License

This is open-sourced software licensed under the [MIT License][link_license].

[badge_language]:https://img.shields.io/github/go-mod/go-version/tarampampam/poke?longCache=true
[badge_build]:https://img.shields.io/github/actions/workflow/status/tarampampam/poke/tests.yml?branch=master&maxAge=30&logo=github
[badge_coverage]:https://img.shields.io/codecov/c/github/tarampampam/poke/master.svg?maxAge=30
[badge_size_latest]:https://img.shields.io/docker/image-size/tarampampam/poke/latest?maxAge=30
[badge_issues]:https://img.shields.io/github/issues/tarampampam/poke.svg?maxAge=45
[badge_pulls]:https://img.shields.io/github/issues-pr/tarampampam/poke.svg?maxAge=45

[link_coverage]:https://codecov.io/gh/tarampampam/poke
[link_build]:https://github.com/tarampampam/poke/actions
[link_docker_hub]:https://hub.docker.com/r/tarampampam/poke/
[link_docker_tags]:https://hub.docker.com/r/tarampampam/poke/tags
[link_license]:https://github.com/tarampampam/poke/blob/master/LICENSE
[link_releases]:https://github.com/tarampampam/poke/releases
[link_commits]:https://github.com/tarampampam/poke/commits
[link_issues]:https://github.com/tarampampam/poke/issues
[link_create_issue]:https://github.com/tarampampam/poke/issues/new/choose
[link_pulls]:https://github.com/tarampampam/poke/pulls
[link_ghcr]:https://github.com/users/tarampampam/packages/container/package/poke

