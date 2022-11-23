#!/usr/bin/env poke run

// convert javascript object into json-string
const payload = JSON.stringify({foo: 123})

// send the request
const resp = fetchSync('https://httpbin.org/anything', {
  method: 'POST',
  headers: {
    'Content-type': 'application/json',
    'X-Foo': 'bar',
  },
  body: JSON.stringify({foo: 123}) // body data type must match "Content-Type" header
})

// now we can assert the remote server response data, like payload:
assert.equals(resp.json().data, payload)
// or an HTTP status code:
assert.equals(resp.json().headers['X-Foo'], 'bar')

// want to set custom user-agent? why not:
assert.equals(
  get('https://httpbin.org/headers', {headers: {'User-Agent': 'jack/daniels'}}).json().headers['User-Agent'],
  'jack/daniels'
)

// additionally, the following helper functions are available:
assert.equals(get('https://httpbin.org/get').status, 200)
assert.equals(post('https://httpbin.org/post', {body: payload}).status, 200)
assert.equals(put('https://httpbin.org/put', {headers: {'User-Agent': 'curl/7.68.0'}}).status, 200)
assert.equals(del('https://httpbin.org/delete').status, 200)
assert.equals(patch('https://httpbin.org/status/444').status, 444)
assert.equals(head('https://httpbin.org/get').status, 200)

// need to go through basic auth? no problem:
assert.equals(
  get('https://httpbin.org/basic-auth/user1/password1', {
    headers: {Authorization: 'Basic ' + encoding.base64encode('user1:password1')}
  }).status,
  200
)

// bearer? sure:
assert.equals(
  get('https://httpbin.org/bearer', {
    headers: {Authorization: 'Bearer user1'}
  }).status, 200
)

// digest authentication? a bit harder, but still possible:
const knock = get('https://httpbin.org/digest-auth/auth/user1/password1').headers['Www-Authenticate']
const [realm, nonce] = [
  knock.match(/realm=\\?"([^\\"]+)\\?"/)[1],
  knock.match(/nonce=\\?"([^\\"]+)\\?"/)[1],
]
assert.equals(
  get('https://httpbin.org/digest-auth/auth/user1/password1', {
    headers: {
      Authorization: `Digest username="user1", realm="${realm}", nonce="${nonce}", uri="/digest-auth/auth/user1/password1", `+
        `response="${hashing.md5([
          hashing.md5(['user1', realm, 'password1'].join(':')),
          nonce,
          hashing.md5(['GET', '/digest-auth/auth/user1/password1'].join(':'))
        ].join(':'))}"`,
    }
  }).status, 200
)
