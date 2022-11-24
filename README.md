<div align="center">
  <img src="https://user-images.githubusercontent.com/7326800/203780071-a963f064-e8bd-4d0c-bbf3-37bdbde03547.png" alt="Logo" width=300" />
  <p>A <strong>simple and powerful</strong> tool for <strong>E2E-testing</strong>.</p>

![Project language][badge_language]
[![Build Status][badge_build]][link_build]
[![Coverage][badge_coverage]][link_coverage]
[![Image size][badge_size_latest]][link_docker_hub]

Just imagine, that you can write E2E tests using familiar javascript syntax without any external dependencies, like
node.js, jest, Axios, etc. Sounds like something impossible? 😉
</div>

> 🚧 The project currently is under active development. Do not use it somewhere except the testing purposes.

## 🔥 Features list

- Write your tests in ECMAScript 5.1 (`let`, `const`, arrow functions and classes are supported)
  syntax ([JS VM](https://github.com/dop251/goja) embedded in the app)
- Tests are parallelized by running them in their own thread to maximize performance
- Ready to work out of the box, config free
- Cross-platform (like any go-based app)
- The jest-like syntax for grouping tests cases and asserts
- Built-in **faker** (fake data generator, like strings, numbers, IP addresses, UUIDs, and so on), **hashing** and **encoding** functions
- Easy to use by developers and QA engineers

## ⚙ Usage

Let's start from the simplest example. Create a file `./test.js` with the following content:

```js
#!/usr/bin/env poke run

describe('test API', () => {
  test('data posting', () => {
    const resp = post('https://httpbin.org/post', {
      body: JSON.stringify({string: faker.string()})
    })

    console.log(resp)

    assert.equals(resp.status, 200)
  })
})
```

And run it:

![shell](https://user-images.githubusercontent.com/7326800/203788784-80b791f9-03c3-4c8a-a9fc-c8b267f16d65.gif)

Wildcards (`?*`) and double stars (`**`) are allowed. For example, you can run all tests from the `./tests` (including all nested) directory:

```shell
./poke run ./tests/**/*.js
```

## 🧩 Installation (WIP)

Download the latest binary file for your arch (to run on macOS use the `linux/arm64` platform) from
the [releases page][link_releases]. For example, let's install it on **amd64** arch (e.g.: Debian, Ubuntu, etc):

```shell
$ curl -SsL -o ./poke https://github.com/tarampampam/poke/releases/latest/download/poke-linux-amd64
$ chmod +x ./poke

# optionally, install the binary file globally:
$ sudo install -g root -o root -t /usr/local/bin -v ./poke
$ rm ./poke
$ poke --help
```
Additionally, you can use the docker image:

| Registry                               | Image                      |
|----------------------------------------|----------------------------|
| [GitHub Container Registry][link_ghcr] | `ghcr.io/tarampampam/poke` |
| [Docker Hub][link_docker_hub]          | `tarampampam/poke`         |

> Using the `latest` tag for the docker image is highly discouraged because of possible backward-incompatible changes
> during **major** upgrades. Please, use tags in `X.Y.Z` format

## 🔌 Language reference ([d.ts](internal/js/global.d.ts))

For the more details take a look the [language reference](docs/language.md) and [examples](examples).

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
[badge_build]:https://img.shields.io/github/workflow/status/tarampampam/poke/tests?maxAge=30&logo=github
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

