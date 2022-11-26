## The simplest example

```javascript
describe('test API', () => {
  test('data posting', () => {
    const [user, password] = ['user1', 'password1']

    const response = get(`https://httpbin.org/basic-auth/${user}/${password}`, {
      headers: {Authorization: 'Basic ' + encoding.base64encode('${user}:${password}')}, // basic auth
      body: JSON.stringify({string: faker.string()}), // generate random string
    })

    assert.equals(response.status, 200)
  })
})

// Run it:
// $ poke run ./test.js
```
