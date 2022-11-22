declare global {
  /**
   * Holds the process details and common functions.
   *
   * @external go Implemented on the Golang side
   * @since 0.0.0
   */
  const process: {
    /** Process environment variables. */
    readonly env: { [key: string]: string }
    /** Pause the script execution for the (at least) given number of milliseconds (1000 == 1 sec.). */
    delay(ms: number): void
  }

  /**
   * Interacting with Input/Output.
   *
   * @external go Implemented on the Golang side
   * @since 0.0.0
   */
  const io: {
    /** Send something to the standard output. */
    stdOut(...v: any[]): void
    /** Send something to the errors output. */
    stdErr(...v: any[]): void
    /** Returns the logging level. */
    logLevel(): 'debug' | 'info' | 'warn' | 'error' // same as log.Level
  }

  /**
   * Events are the best way to communicate with the script runner (go-side).
   *
   * @external go Implemented on the Golang side
   * @since 0.0.0
   */
  const events: {
    /** Send an event to the go-side. */
    push(...events: {
      /** Event level ('debug' by default). */
      level?: 'debug' | 'info' | 'warning' | 'error'
      /** Event message. */
      message: string
      /** An error (optional). */
      error?: Error
    }[]): void
  }

  /**
   * Send an HTTP request (synchronously).
   *
   * @external go Implemented on the Golang side
   * @since 0.0.0
   */
  function fetchSync(url: string, options?: {
    /** *GET, POST, PUT, DELETE, etc. */
    method?: string
    /** Request headers map. */
    headers?: Record<string, string>
    /** Request body. Body data type must match "Content-Type" header. */
    body?: string
  }): {
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

  /**
   * Assertion functions.
   *
   * @since 0.0.0
   */
  const assert: {
    /**
     * Asserts that the value is truthy.
     *
     * @since 0.0.0
     */
    true(mustBeTrue: unknown, message?: string): void

    /**
     * Asserts that the values are the same.
     *
     * @since 0.0.0
     */
    same(actual: unknown, expected: unknown, message?: string): void
  }

  /**
   * Runs a function before any of the tests in this file run.
   *
   * @example
   * beforeAll(() => {
   *   console.log('before all tests in the file')
   * })
   *
   * @since 0.0.0
   */
  function beforeAll(fn: () => void): void

  /**
   * Runs a function before each of the tests in this file runs.
   *
   * @example
   * beforeEach(() => {
   *   console.log('before each test')
   * })
   *
   * @since 0.0.0
   */
  function beforeEach(fn: () => void): void

  /**
   * Runs a function after each one of the tests in this file completes.
   *
   * @example
   * afterEach(() => {
   *   console.log('after each test')
   * })
   *
   * @since 0.0.0
   */
  function afterEach(fn: () => void): void

  /**
   * Runs a function after all the tests in this file have completed.
   *
   * @example
   * afterAll(() => {
   *   console.log('after all tests in the file')
   * })
   *
   * @since 0.0.0
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
   *
   * @since 0.0.0
   */
  function test(name: string, fn: () => void): void

  /**
   * Is an alias for the test() function.
   *
   * @alias test
   * @since 0.0.0
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
   *
   * @since 0.0.0
   */
  function describe(name: string, fn: () => void): void
}

export {}
