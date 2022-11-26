---
title: "Quick Start"
slug: quick-start
description: "Learn to start Poke tests in minutes"
lead: "Learn to start Poke tests in minutes."
weight: 120
toc: false
---

## Prerequisites

Before you begin you must [install Poke]({{< ref "/docs/getting-started/installation" >}}). You must also be comfortable
working from the command line.

### Create the first test

Using your favorite IDE or a text editor create a file with the name `test.js`. Write into the file the following code:

```javascript
const response = get(`https://httpbin.org/get`)

console.debug(response.headers) // print response headers
console.debug(response.json())  // print response body as JSON

mustBe.equals(response.status, 200)   // check response status code
assert.true(response.body.length > 0) // check response body length
```

{{< alert icon="📑" >}}
Read the [language reference]({{< ref "/docs/lang/" >}}) to learn more about [assertions]({{< ref "/docs/lang/asserts" >}}) and built-in helpers.
{{< /alert >}}

Save it, and run:

```bash
$ poke --log-level debug run test.js
```

If all asserts will be passed, the exit code will be `0`. Otherwise, the exit code will be `1`. Additionally, you can
use `--log-level` option to set the log level. Available values are `debug`, `info`, `warn`, and `error`.

{{< alert icon="📑" >}}
All supported CLI commands, flags, and so on are described [here]({{< ref "/docs/getting-started/cli" >}}).
{{< /alert >}}

If you need to test a large project, you may want to organize your tests in different directories. For example:

```bash
./tests/
├── api
│   ├── pets # pets API tests
│   │   ├── create.poke.js
│   │   ├── move.poke.js
│   │   └── search.poke.js
│   └── session # session API tests
│       ├── create.poke.js
│       └── drop.poke.js
└── static
    ├── index_page.poke.js # test for index page
    └── robots_txt.poke.js # test for robots.txt
```

And run them all using the following command:

```bash
$ poke run ./tests/**/*.poke.js
```

Or only a "pets" API tests:

```bash
$ poke run ./tests/api/pets/*.poke.js
```

## Ask for help

Feel free to ask questions in the [GitHub discussions][discussions], or report bugs/to ask for features in the [GitHub issues][new-issue].

[discussions]:https://github.com/tarampampam/poke/discussions
[new-issue]:https://github.com/tarampampam/poke/issues/new/choose
