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

    if (typeof message !== 'string' || message === "") {
      message = 'Expected true but got ' + String(mustBeTrue)
    }

    console.error(message, mustBeTrue)
    events.push({level: 'error', message: message})

    if (interrupt === true) {
      process.interrupt(message)
    }
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

    if (typeof message !== 'string' || message === "") {
      message = 'Expected false but got ' + String(mustBeFalse)
    }

    console.error(message, mustBeFalse)
    events.push({level: 'error', message: message})

    if (interrupt === true) {
      process.interrupt(message)
    }
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

    if (typeof message !== 'string' || message === "") {
      message = String(actual) + ' and ' + String(expected) + ' are not the same'
    }

    console.error(message, actual, expected)
    events.push({level: 'error', message: message})

    if (interrupt === true) {
      process.interrupt(message)
    }
  }

  /**
   * @param {*} object
   * @param {string?} message
   * @param {boolean?} interrupt
   */
  empty(object, message, interrupt) {
    if (!object) {
      return
    }

    switch (typeof object) {
      case 'undefined':
        return

      case 'string' && object === "":
        return

      case 'boolean' && object === false:
        return

      case 'number' && object === 0:
        return

      case 'bigint' && object === 0:
        return

      case 'object' && object === null:
        return

      case 'object': {
        if (object !== null) {
          switch (true) {
            case Array.isArray(object) && object.length === 0: // empty array
              return

            case Object.keys(object).length === 0: // empty collection
              return

            case object instanceof Date && Number.isNaN(object.getTime()): // is invalid date
              return

            case object instanceof Set && object.size === 0: // empty set
              return

            case object instanceof Map && object.size === 0: // empty map
              return
          }
        }
      }
    }

    if (typeof message !== 'string' || message === "") {
      message = String(object) + ' is not empty'
    }

    console.error(message, object)
    events.push({level: 'error', message: message})

    if (interrupt === true) {
      process.interrupt(message)
    }
  }
}

const require = new class {
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
   * @param {*} object
   * @param {string?} message
   */
  empty(object, message) {
    assert.empty(object, message, true)
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
