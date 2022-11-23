#!/usr/bin/env poke run

beforeAll(() => console.log('> before all'))
beforeEach((testName) => console.log(`>> before test "${testName}"`))
afterEach((testName) => console.log(`<< after test "${testName}"`))
afterAll(() => console.log('< after all'))

describe('boolean values', () => {
  test('true', () => assert.true(true))

  it('false', () => {
    require.false(false)
  })
})

describe('strings', () => {
  test('equals', () => assert.equals('bar', 'bar'))

  describe('sub-strings', () => {
    test('length', () => require.true('foo'.length === 3))
    test('numbers', () => require.equals(1, 1))
  })
})
