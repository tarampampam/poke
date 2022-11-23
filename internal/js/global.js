/** Send HTTP request by GET method. */
const get = (url, options) => fetchSync(url, {method: 'GET', ...options})
/** Send HTTP request by POST method. */
const post = (url, options) => fetchSync(url, {method: 'POST', ...options})
/** Send HTTP request by PUT method. */
const put = (url, options) => fetchSync(url, {method: 'PUT', ...options})
/** Send HTTP request by DELETE method. */
const del = (url, options) => fetchSync(url, {method: 'DELETE', ...options})
/** Send HTTP request by PATCH method. */
const patch = (url, options) => fetchSync(url, {method: 'PATCH', ...options})
/** Send HTTP request by HEAD method. */
const head = (url, options) => fetchSync(url, {method: 'HEAD', ...options})

/** @internal */
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

  /**
   * @param {Map<string|Symbol, Function>} m
   * @param {function(function, string|Symbol)} executor
   */
  reduceMap(m, executor) {
    while (m.size > 0) {
      const [name, fn] = m.entries().next().value

      m.delete(name)

      executor(fn, name)
    }
  }

  /** Run all the tests. */
  run() {
    this.reduceMap(this.describeQueue, (fn) => fn())

    if (this.testsQueue.size > 0) { // run tests
      this.reduceMap(this.beforeAll, (fn) => fn())

      this.reduceMap(this.testsQueue, (fn, name) => {
        this.beforeEach.forEach((fn) => fn(name))

        fn()

        this.afterEach.forEach((fn) => fn(name))
      })

      this.reduceMap(this.afterAll, (fn) => fn())
    }

    if (this.describeQueue.size > 0) {
      this.run() // recursive run
    }
  }

  /**
   * @param {*} object
   * @returns {boolean}
   */
  isEmpty(object) {
    if (!object) {
      return true
    }

    switch (typeof object) {
      case 'undefined':
        return true

      case 'string' && object === "":
        return true

      case 'boolean' && object === false:
        return true

      case 'number' && object === 0:
        return true

      case 'bigint' && object === 0:
        return true

      case 'object' && object === null:
        return true

      case 'object': {
        if (object !== null) {
          switch (true) {
            case Array.isArray(object) && object.length === 0: // empty array
              return true

            case object instanceof Date && Number.isNaN(object.getTime()): // is invalid date
              return true

            case object.constructor === Object && Object.keys(object).length === 0: // empty collection
              return true

            case object instanceof Set && object.size === 0: // empty set
              return true

            case object instanceof Map && object.size === 0: // empty map
              return true
          }
        }
      }
    }

    return false
  }

  /**
   * @param {*} message
   * @return {boolean}
   */
  notEmptyString(message) {
    return typeof message === 'string' && message.trim() !== ""
  }

  /**
   * @param {string} message
   * @param {boolean} interrupt
   * @param {*[]} consoleArgs
   */
  triggerError(message, interrupt, ...consoleArgs) {
    console.error(message, ...consoleArgs)
    events.push({level: 'error', message: message})

    if (interrupt === true) {
      process.interrupt(message)
    }
  }
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

/** Assertion functions. */
const assert = new class {
  /**
   * @param {*} mustBeTrue
   * @param {string?} message
   * @param {boolean?} interrupt
   */
  true(mustBeTrue, message, interrupt) {
    if (mustBeTrue === true) {
      return
    }

    message = tests.notEmptyString(message)
      ? message
      : 'Expected true but got ' + String(mustBeTrue)

    tests.triggerError(message, interrupt, mustBeTrue)
  }

  /**
   * @param {*} mustBeFalse
   * @param {string?} message
   * @param {boolean?} interrupt
   */
  false(mustBeFalse, message, interrupt) {
    if (mustBeFalse === false) {
      return
    }

    message = tests.notEmptyString(message)
      ? message
      : 'Expected false but got ' + String(mustBeFalse)

    tests.triggerError(message, interrupt, mustBeFalse)
  }

  /**
   * @param {*} actual
   * @param {*} expected
   * @param {string?} message
   * @param {boolean?} interrupt
   */
  equals(actual, expected, message, interrupt) {
    if (actual === expected) {
      return
    }

    message = tests.notEmptyString(message)
      ? message
      : String(actual) + ' and ' + String(expected) + ' are not the same'

    tests.triggerError(message, interrupt, actual, expected)
  }

  /**
   * @param {*} actual
   * @param {*} expected
   * @param {string?} message
   * @param {boolean?} interrupt
   */
  notEquals(actual, expected, message, interrupt) {
    if (actual !== expected) {
      return
    }

    message = tests.notEmptyString(message)
      ? message
      : String(actual) + ' and ' + String(expected) + ' are the same, but they should not be'

    tests.triggerError(message, interrupt, actual, expected)
  }

  /**
   * @param {*} object
   * @param {string?} message
   * @param {boolean?} interrupt
   */
  empty(object, message, interrupt) {
    if (tests.isEmpty(object)) {
      return
    }

    message = tests.notEmptyString(message)
      ? message
      : String(object) + ' is not empty'

    tests.triggerError(message, interrupt, object)
  }

  /**
   * @param {*} object
   * @param {string?} message
   * @param {boolean?} interrupt
   */
  notEmpty(object, message, interrupt) {
    if (!tests.isEmpty(object)) {
      return
    }

    message = tests.notEmptyString(message)
      ? message
      : String(object) + ' is empty, but should not be'

    tests.triggerError(message, interrupt, object)
  }
}

const mustBe = new class {
  /**
   * @param {*} mustBeTrue
   * @param {string?} message
   */
  true(mustBeTrue, message) {
    assert.true(mustBeTrue, message, true)
  }

  /**
   * @param {*} mustBeFalse
   * @param {string?} message
   */
  false(mustBeFalse, message) {
    assert.false(mustBeFalse, message, true)
  }

  /**
   * @param {*} actual
   * @param {*} expected
   * @param {string?} message
   */
  equals(actual, expected, message) {
    assert.equals(actual, expected, message, true)
  }

  /**
   * @param {*} actual
   * @param {*} expected
   * @param {string?} message
   */
  notEquals(actual, expected, message) {
    assert.notEquals(actual, expected, message, true)
  }

  /**
   * @param {*} object
   * @param {string?} message
   */
  empty(object, message) {
    assert.empty(object, message, true)
  }

  /**
   * @param {*} object
   * @param {string?} message
   */
  notEmpty(object, message) {
    assert.notEmpty(object, message, true)
  }
}

/**
 * This function will be called at the end of the script execution by the Go runtime.
 *
 * @internal
 */
const init = () => {
  tests.run()
}
