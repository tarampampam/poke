declare global {
  const process: {
    env: { [key: string]: string }
  }

  const io: {
    stdOut(...v): void // Send something to the standard output
    stdErr(...v): void // Send something to the errors output
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
