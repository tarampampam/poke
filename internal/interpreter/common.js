// // const resp = await(fetch('https://httpbin.org/get'))
// const resp = await(fetch('https://httpbin.org/get',  {
//   method: 'get', // *GET, POST, PUT, DELETE, etc.
//   mode: 'cors', // no-cors, *cors, same-origin
//   cache: 'no-cache', // *default, no-cache, reload, force-cache, only-if-cached
//   credentials: 'same-origin', // include, *same-origin, omit
//   headers: {
//     'Content-type': 'application/json'
//     // 'Content-Type': 'application/x-www-form-urlencoded',
//   },
//   redirect: 'follow', // manual, *follow, error
//   referrerPolicy: 'no-referrer', // no-referrer, *no-referrer-when-downgrade, origin, origin-when-cross-origin, same-origin, strict-origin, strict-origin-when-cross-origin, unsafe-url
//   body: JSON.stringify({foo: 123}) // body data type must match "Content-Type" header
// }))
//
// console.log(resp.headers)
// console.log(resp.json())

// console.log(await(new Promise((resolve, reject) => { resolve(123123123) })))

new Promise((resolve, reject) => { resolve(123123123) }).catch(console.error)

// console.log(new Promise(resolve => {}))

// console.log(process.env)

// for (const property in process.env) {
//   console.log(`${property}: ${process.env[property]}`)
// }


