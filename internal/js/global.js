/** Simple replacement for the Console implementation. */
const console = new class {
  /**
   * @param {*} v
   * @return {boolean}
   */
  isPromise(v) {
    return v && typeof v === 'object' && v instanceof Promise
  }

  /**
   * @param {*} v
   * @return {string|number}
   */
  fmtValue(v) {
    const type = typeof v

    switch (typeof v) {
      case 'function':
        return `ƒ(…)`

      case 'object':
        if (Array.isArray(v)) {
          return '[…]'
        } else if (this.isPromise(v)) {
          return '<Promise>'
        } else if (v === null) {
          return 'null'
        }
        return '{…}'

      case 'string':
        return `"${v}"`

      case 'number':
      case 'bigint':
        return v

      case 'boolean':
        return v ? 'true' : 'false'

      case 'symbol':
        return '<Symbol>'

      default:
        return type
    }
  }

  /**
   * @param {*} v
   * @return {String}
   */
  fmt(...v) {
    const parts = new Array(v.length)

    for (let i = 0; i < v.length; i++) {
      const current = v[i]

      switch (true) {
        case Array.isArray(current):
          parts[i] = `[${current.map(this.fmtValue).join(', ')}]`
          break

        case this.isPromise(current):
          parts[i] = this.fmtValue(current)
          break

        case current instanceof Error:
          parts[i] = current.toString()
          break

        case typeof current === 'object' && current !== null: // watch 1 level deep
          let props = []

          for (let id in current) {
            const type = typeof current[id]
            const value = current[id]

            switch (type) {
              case 'function':
                props.unshift(`${id}: ${this.fmtValue(value)}`) // always first
                break

              default:
                props.push(`${id}: ${this.fmtValue(value)}`)
            }
          }

          parts[i] = `{${props.join(', ')}}`
          break

        default:
          parts[i] = this.fmtValue(current)
      }
    }

    return parts.join(', ').toString() + '\n'
  }

  /**
   * @param {String} logLevel
   * @return {number}
   */
  logLevelToInt(logLevel) {
    switch (logLevel) {
      case 'debug':
        return -1

      case 'info':
        return 0

      case 'warn':
        return 1

      case 'error':
        return 2

      default:
        return 0  // as an info level
    }
  }

  /**
   * @param {'debug' | 'info' | 'warn' | 'error'} logLevel
   * @return {boolean}
   */
  checkLevel(logLevel) {
    return this.logLevelToInt(logLevel) >= this.logLevelToInt(io.logLevel())
  }

  /** Log a message at debug level. */
  debug(...v) {
    if (this.checkLevel('debug')) {
      io.stdOut(this.fmt(...v))
    }
  }

  /** Log a message at info level. */
  log(...v) {
    if (this.checkLevel('info')) {
      io.stdOut(this.fmt(...v))
    }
  }

  /** Log a message at info level. */
  info(...v) {
    if (this.checkLevel('info')) {
      io.stdOut(this.fmt(...v))
    }
  }

  /** Log a message at warn level. */
  warn(...v) {
    if (this.checkLevel('warn')) {
      io.stdOut(this.fmt(...v))
    }
  }

  /** Log a message at error level. */
  error(...v) {
    if (this.checkLevel('error')) {
      io.stdErr(this.fmt(...v))
    }
  }
}

/**
 * @internal
 */
const tests = new class {
  /** @type {Map<Symbol, Function>} */
  beforeAll = new Map()
  /** @type {Map<Symbol, Function>} */
  beforeEach = new Map()
  /** @type {Map<Symbol, Function>} */
  afterEach = new Map()
  /** @type {Map<Symbol, Function>} */
  afterAll = new Map()
  /** @type {Map<string, Function>} */
  testsQueue = new Map()
  /** @type {Map<string, Function>} */
  describeQueue = new Map()
}

/** Runs a function before any of the tests in this file run. */
const beforeAll = (fn) => tests.beforeAll.set(Symbol(), fn)

/** Runs a function before each of the tests in this file runs. */
const beforeEach = (fn) => tests.beforeEach.set(Symbol(), fn)

/** Runs a function after each one of the tests in this file completes. */
const afterEach = (fn) => tests.afterEach.set(Symbol(), fn)

/** Runs a function after all the tests in this file have completed. */
const afterAll = (fn) => tests.afterAll.set(Symbol(), fn)

/** All you need in a test file is the test method which runs a test. */
const test = (name, fn) => tests.testsQueue.set(name, fn)

/** Is an alias for the test() function. */
const it = (name, fn) => test(name, fn)

/** Creates a block that groups together several related tests. */
const describe = (name, fn) => tests.describeQueue.set(name, fn)

/**
 * This function will be called at the end of the script execution by the Go runtime.
 *
 * @internal
 */
const __afterScript = () => {
  const bootstrapTests = () => {
    while (tests.describeQueue.size > 0) { // run code inside describe groups
      const [groupName, describeFn] = tests.describeQueue.entries().next().value

      tests.describeQueue.delete(groupName)

      console.debug(`» Running group: ${groupName}`)
      describeFn()
    }

    if (tests.testsQueue.size > 0) { // run tests
      if (tests.beforeAll.size > 0) { // run "before all" hooks (once)
        console.debug(`Running "before all" hooks (${tests.beforeAll.size})`)

        tests.beforeAll.forEach((fn) => fn())
        tests.beforeAll.clear()
      }

      while (tests.testsQueue.size > 0) {
        const [testName, testFn] = tests.testsQueue.entries().next().value

        tests.testsQueue.delete(testName)

        if (tests.beforeEach.size > 0) { // run "before each test" hooks
          console.debug(`Running "before each" hooks (${tests.beforeEach.size})`)

          tests.beforeEach.forEach((fn) => fn())
        }

        console.debug(`> Running test: ${testName}`)
        testFn()

        if (tests.afterEach.size > 0) { // run "after each test" hooks
          console.debug(`Running "after each" hooks (${tests.afterEach.size})`)

          tests.afterEach.forEach((fn) => fn())
        }
      }

      if (tests.afterAll.size > 0) { // run "after all" hooks (once)
        console.debug(`Running "after all" hooks (${tests.afterAll.size})`)

        tests.afterAll.forEach((fn) => fn())
        tests.afterAll.clear()
      }
    }

    if (tests.describeQueue.size > 0) {
      bootstrapTests()
    }
  }

  bootstrapTests()
}

/** Assertion functions. */
const assert = new class {
  /**
   * @param {*} mustBeTrue
   * @param {string?} message
   */
  true(mustBeTrue, message) {
    if (mustBeTrue === true) {
      return
    }

    if (message === undefined) {
      message = 'Expected true but got ' + String(mustBeTrue)
    }

    events.push({level: 'error', message: 'Assertion failed', error: new Error(message)})
    // throw new Error(message)
  }

  /**
   * @param {*} actual
   * @param {*} expected
   * @param {string?} message
   */
  same(actual, expected, message) {
    if (message === undefined) {
      message = String(actual) + ' and ' + String(expected) + ' are not the same'
    }

    assert.true(actual === expected, message)
  }
}
