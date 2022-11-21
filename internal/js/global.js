/**
 * @param {*} v
 * @return {boolean}
 */
const isPromise = (v) => {
  return v && typeof v === 'object' && v instanceof Promise
}

/**
 * Simple replacement for the Console implementation.
 */
const console = new class {
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
        } else if (isPromise(v)) {
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
   * @return {Array<string|number>}
   */
  fmt(...v) {
    const parts = new Array(v.length)

    for (let i = 0; i < v.length; i++) {
      const current = v[i]

      switch (true) {
        case Array.isArray(current):
          parts[i] = `[${current.map(this.fmtValue).join(', ')}]`
          break

        case isPromise(current):
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

    return parts
  }

  log(...v) {
    io.stdOut(this.fmt(...v).join(', ') + '\n')
  }

  error(...v) {
    io.stdErr(this.fmt(...v).join(', ') + '\n')
  }

  debug(...v) {
    this.log(...v)
  }

  info(...v) {
    this.log(...v)
  }

  warn(...v) {
    this.log(...v)
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

/** Is an alias for the test() function. */
const describe = (name, fn) => test(name, fn)

/**
 * This function will be called at the end of the script execution by the Go runtime.
 *
 * @internal
 */
const __afterScript = () => {
  if (tests.testsQueue.size > 0) {
    tests.beforeAll.forEach((fn) => fn())

    tests.testsQueue.forEach((fn, name) => {
      tests.beforeEach.forEach((fn) => fn())
      fn() // execute test
      tests.afterEach.forEach((fn) => fn())
    })

    tests.afterAll.forEach((fn) => fn())
  }
}

/**
 * Assertion functions.
 */
const assert = new class {
  /**
   * @param {*} mustBeTrue
   * @param {string?} message
   *
   * @throws
   */
  true(mustBeTrue, message) {
    if (mustBeTrue === true) {
      return
    }

    if (message === undefined) {
      message = 'Expected true but got ' + String(mustBeTrue)
    }

    throw new Error(message)
  }

  /**
   * @param {*} actual
   * @param {*} expected
   * @param {string?} message
   *
   * @throws
   */
  same(actual, expected, message) {
    if (message === undefined) {
      message = String(actual) + ' and ' + String(expected) + ' are not the same'
    }

    assert.true(actual === expected, message)
  }
}
