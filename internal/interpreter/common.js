const a = 1

// const resp = await(fetch('https://httpbin.org/get'))
const resp = await(fetch('https://httpbin.org/get',  {
  method: 'get', // *GET, POST, PUT, DELETE, etc.
  mode: 'cors', // no-cors, *cors, same-origin
  cache: 'no-cache', // *default, no-cache, reload, force-cache, only-if-cached
  credentials: 'same-origin', // include, *same-origin, omit
  headers: {
    'Content-type': 'application/json'
    // 'Content-Type': 'application/x-www-form-urlencoded',
  },
  redirect: 'follow', // manual, *follow, error
  referrerPolicy: 'no-referrer', // no-referrer, *no-referrer-when-downgrade, origin, origin-when-cross-origin, same-origin, strict-origin, strict-origin-when-cross-origin, unsafe-url
  body: JSON.stringify({foo: 123}) // body data type must match "Content-Type" header
}))

console.log(resp)
console.log(resp.body)
console.log(resp.json())

console.log(await(new Promise((resolve, reject) => { resolve(1111111111) })))

// console.log(new Promise(resolve => {}))

// for (const property in process.env) {
//   console.log(`${property}: ${process.env[property]}`)
// }


