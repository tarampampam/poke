console.log('test', 123, {}, [], [1, 666, 3], undefined, null)

const a = 1

// const resp = await fetch('https://httpbin.org/get')
// console.log(resp)

// console.log(new Promise(resolve => {}))

for (const property in process.env) {
  console.log(`${property}: ${process.env[property]}`)
}
