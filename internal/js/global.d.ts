interface FetchSyncOptions {
  /** Request method (`GET` by default) */
  method?: 'GET' | 'HAD' | 'POST' | 'PUT' | 'PATCH' | 'DELETE' | 'CONNECT' | 'OPTIONS' | 'TRACE'
  /** Request headers map. */
  headers?: Record<string, string>
  /** Request body. Body data type must match "Content-Type" header. */
  body?: string
}

interface FetchSyncResponse {
  /** Body contents. */
  readonly body: string
  /** Response headers map. */
  readonly headers: Record<string, string>
  /** A boolean indicating whether the response was successful (status in the range 200 – 299) or not. */
  readonly ok: boolean
  /** The status code of the response. */
  readonly status: number
  /** The status message corresponding to the status code. */
  readonly statusText: string
  /** The URL of the response. */
  readonly url: string
  /** Returns body data as an ArrayBuffer */
  arrayBuffer(): ArrayBuffer
  /** Returns a result of parsing the response body text as JSON. */
  json(): unknown
  /** Returns a text representation of the response body */
  text(): string
}

declare global {
  /**
   * Holds the process details and common functions.
   *
   * @external go Implemented on the Golang side
   */
  const process: {
    /** Process environment variables. */
    readonly env: { [key: string]: string }
    /** Pause the script execution for the (at least) given number of milliseconds (1000 == 1 sec.). */
    delay(ms: number): void
    /** Interrupt the script execution. */
    interrupt(reason?: string): void
  }

  /**
   * Interacting with Input/Output.
   *
   * @external go Implemented on the Golang side
   */
  const io: {
    /** Send something to the standard output. */
    stdOut(...v: unknown[]): void
    /** Send something to the errors output. */
    stdErr(...v: unknown[]): void
  }

  /**
   * Events are the best way to communicate with the script runner (go-side).
   *
   * @external go Implemented on the Golang side
   */
  const events: {
    /** Send an event to the go-side. */
    push(...events: {
      /** Event level ('debug' by default). */
      level?: 'debug' | 'info' | 'warning' | 'error'
      /** Event message. */
      message: string
      /** An error (optional). */
      error?: string | Error
    }[]): void
  }

  // @ts-ignore
  const console: {
    /** Log a message with debug level. */
    debug(...data: unknown[]): void
    /** Log a message with info level. */
    log(...data: unknown[]): void
    /** Log a message with info level. */
    info(...data: unknown[]): void
    /** Log a message with warning level. */
    warn(...data: unknown[]): void
    /** Log a message with error level. */
    error(...data: unknown[]): void
  }

  /**
   * Send an HTTP request (synchronously).
   *
   * @external go Implemented on the Golang side
   */
  function fetchSync(url: string, options?: FetchSyncOptions): FetchSyncResponse

  /** Send HTTP request by GET method. */
  function get(url: string, options?: FetchSyncOptions): FetchSyncResponse
  /** Send HTTP request by POST method. */
  function post(url: string, options?: FetchSyncOptions): FetchSyncResponse
  /** Send HTTP request by PUT method. */
  function put(url: string, options?: FetchSyncOptions): FetchSyncResponse
  /** Send HTTP request by DELETE method. */
  function del(url: string, options?: FetchSyncOptions): FetchSyncResponse
  /** Send HTTP request by PATCH method. */
  function patch(url: string, options?: FetchSyncOptions): FetchSyncResponse
  /** Send HTTP request by HEAD method. */
  function head(url: string, options?: FetchSyncOptions): FetchSyncResponse

  /**
   * Encoding helper functions.
   *
   * @external go Implemented on the Golang side
   */
  const encoding: {
    /** Encode a string to base64 (`std` mode is used by default). */
    base64encode(s: string, options?: {mode: 'std' | 'url'}): string
    /** Decode a base64 string (`std` mode us used is default). `undefined` will be returned on malformed input. */
    base64decode(encoded: string, options?: {mode: 'std' | 'url'}): string | undefined
  }

  /**
   * Hashing helper functions.
   *
   * @external go Implemented on the Golang side
   */
  const hashing: {
    /** Calculate a MD5 hash of the given string. */
    md5(s: string): string
    /** Calculate a Sha256 hash of the given string. */
    sha256(s: string): string
  }

  /**
   * Helper functions to generate a faked data.
   *
   * @external go Implemented on the Golang side
   */
  const faker: {
    /** Returns a random bool. */
    bool(): boolean
    /** Returns a random falsy value. */
    falsy(): false, null, undefined, 0, NaN, ''
    /** Returns a random character (a-zA-Z0-9 by default). */
    character(options?: {pool: string}): string
    /** Returns a random floating point number. */
    floating(): number
    /** Returns a random integer (-2147483648 to 2147483648 by default). */
    integer(options?: {min?: number, max?: number}): number
    /** Returns a random letter (one of "abcdefghijklmnopqrstuvwxyz"). */
    letter(): string
    /** Returns a random string. */
    string(options?: {length?: number, pool?: string}): string
    /** Returns a random paragraph generated from sentences populated by semi-pronounceable random (nonsense) words. */
    paragraph(): string
    /** Returns a semi-pronounceable random (nonsense) word. */
    word(): string
    /** Returns a random domain with a random tld (like `foobar.org`). */
    domain(): string
    /** Returns a random email with a random domain. */
    email(): string
    /** Returns a random IPv4 Address. */
    ip(): string
    /** Returns a random IPv6 Address. */
    ipv6(): string
    /** Returns a random tld (Top Level Domain). */
    tld(): string
    /** Returns a random url. */
    url(): string
    /** Generates a random date. */
    date(): Date
    /** Returns a random hex hash. */
    hash(options?: {length: number}): string
    /** Returns a random UUID. */
    uuid(): string
    /** Returns a randomly picked argument. */
    random<T>(v1: T, v2: T, ...v: T[]): T
  }

  /** Assertion functions. */
  const assert: {
    /** Asserts that the value is truthy. */
    true(mustBeTrue: unknown, message?: string, interrupt?: boolean): void
    /** Asserts that the value is falsely. */
    false(mustBeTrue: unknown, message?: string, interrupt?: boolean): void
    /** Asserts that the values are the same. */
    equals(actual: unknown, expected: unknown, message?: string, interrupt?: boolean): void
    /** Asserts that the values are not the same. */
    notEquals(actual: unknown, expected: unknown, message?: string, interrupt?: boolean): void
    /** Asserts that the specified object is empty. */
    empty(object: unknown, message?: string, interrupt?: boolean): void
    /** Asserts that the specified object is not empty. */
    notEmpty(object: unknown, message?: string, interrupt?: boolean): void
    /** Asserts that the specified string or array contains the required value. */
    contains(what: unknown, where: string | unknown[], message?: string, interrupt?: boolean): void
    /** Asserts that the specified string or array does not contain the required value. */
    notContains(what: unknown, where: string | unknown[], message?: string, interrupt?: boolean): void
  }

  /** Assertion functions that interrupt the script on error. */
  const mustBe: {
    /** Asserts that the value is truthy. */
    true(mustBeTrue: unknown, message?: string): void
    /** Asserts that the value is falsely. */
    false(mustBeTrue: unknown, message?: string): void
    /** Asserts that the values are the same. */
    equals(actual: unknown, expected: unknown, message?: string): void
    /** Asserts that the values are not the same. */
    notEquals(actual: unknown, expected: unknown, message?: string): void
    /** Asserts that the specified object is empty. */
    empty(object: unknown, message?: string): void
    /** Asserts that the specified object is not empty. */
    notEmpty(object: unknown, message?: string): void
    /** Asserts that the specified string or array contains required value. */
    contains(what: unknown, where: string | unknown[], message?: string): void
    /** Asserts that the specified string or array does not contain the required value. */
    notContains(what: unknown, where: string | unknown[], message?: string): void
  }

  /**
   * Runs a function before any of the tests in this file run.
   *
   * @example
   * beforeAll(() => {
   *   console.log('before all tests in the file')
   * })
   */
  function beforeAll(fn: () => void): void

  /**
   * Runs a function before each of the tests in this file runs.
   *
   * @example
   * beforeEach(() => {
   *   console.log('before each test')
   * })
   */
  function beforeEach(fn: (testName: string) => void): void

  /**
   * Runs a function after each one of the tests in this file completes.
   *
   * @example
   * afterEach(() => {
   *   console.log('after each test')
   * })
   */
  function afterEach(fn: (testName: string) => void): void

  /**
   * Runs a function after all the tests in this file have completed.
   *
   * @example
   * afterAll(() => {
   *   console.log('after all tests in the file')
   * })
   */
  function afterAll(fn: () => void): void

  /**
   * All you need in a test file is the test method which runs a test.
   *
   * For example, let's say that the response code should be 200. Your whole test could be:
   *
   * @example
   * test('response code should be 200', () => {
   *   assert.true(fetchSync('https://cdnjs.com/').status === 200)
   * })
   */
  function test(name: string, fn: () => void): void

  /**
   * Is an alias for the test() function.
   *
   * @alias test
   */
  function it(name: string, fn: () => void): void

  /**
   * Creates a block that groups together several related tests.
   *
   * @example
   * describe('boolean tests', () => {
   *   test('is true', () => {
   *     assert.true(true)
   *   })
   *
   *   test('is false', () => {
   *     assert.true(false)
   *   })
   * })
   */
  function describe(name: string, fn: () => void): void
}

export {}
