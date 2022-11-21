declare global {
  const process: {
    env: { [key: string]: string }
    delay(ms: number): void // Pause the script execution for the given number of milliseconds (1000 == 1 sec.)
  }

  const io: {
    stdOut(...v: any[]): void // Send something to the standard output
    stdErr(...v: any[]): void // Send something to the errors output
  }

  const events: {
    push(...reports: {
      level?: 'debug' | 'info' | 'warning' | 'error'
      message: string
      error?: Error
    }[]): void
  }

  const assert: {
    true(mustBeTrue: unknown, message?: string)
    same(actual: unknown, expected: unknown, message?: string)
  }

  function fetchSync(url: string, options?: {
    method?: string // *GET, POST, PUT, DELETE, etc.
    headers?: Record<string, string> // {'Content-Type': 'application/x-www-form-urlencoded', ...}
    body?: string // Body data type must match "Content-Type" header
  }): {
    body: string // Body contents
    headers: Record<string, string> // Response headers
    ok: boolean // A boolean indicating whether the response was successful (status in the range 200 – 299) or not
    status: number // The status code of the response
    statusText: string // the status message corresponding to the status code
    url: string // The URL of the response
    arrayBuffer(): ArrayBuffer  // Returns body data as an ArrayBuffer
    json(): unknown // Returns a result of parsing the response body text as JSON
    text(): string // Returns a text representation of the response body
  }
}

export {}
