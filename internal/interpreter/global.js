/**
 * Simple `await` keyword wrapper.
 *
 * @param {*} args
 * @return {*}
 */
const await = function(...args) {
  if (args.length !== 1) {
    throw new Error('await must be called with exactly 1 argument')
  }

  return args[0]
}
