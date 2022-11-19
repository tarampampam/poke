declare global {
  let io: {
    stdOut(...v): void // Send something to the standard output
    stdErr(...v): void // Send something to the errors output
  }

  let process: {
    env: {[key: string]: string}
    gosched(): void
  }
}

export {}
