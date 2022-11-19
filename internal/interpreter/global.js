const isPromise = (v) => {
  return v && typeof v === 'object' && typeof v.then === 'function' && typeof v.catch === 'function'
}

const console = new class {
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
    io.stdOut(this.fmt(...v).join(', '), '\n')
  }

  error(...v) {
    io.stdErr(this.fmt(...v).join(', '), '\n')
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
 * Simple `await` keyword wrapper.
 *
 * @param {*} args
 * @return {*}
 */
const await = function (...args) {
  if (args.length !== 1) {
    throw new Error('await must be called with exactly 1 argument')
  }

  if (isPromise(args[0])) {
    let resolved = false
    let result = undefined

    args[0].then(v => {
      resolved = true
      result = v
    })

    while (!resolved) {
      process.gosched()
    }

    return result
  }

  return args[0]
}
