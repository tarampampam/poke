#!/usr/bin/env poke run

const resp = fetchSync('https://httpbin.org/get', {
  method: 'get', // *GET, POST, PUT, DELETE, etc.
  headers: {
    'Content-type': 'application/json',
    'X-Foo': 'bar'
  },
  body: JSON.stringify({foo: 123}) // body data type must match "Content-Type" header
})

console.log(resp.headers)
console.log(resp.body)
