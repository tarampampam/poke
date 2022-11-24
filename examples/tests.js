#!/usr/bin/env poke run

beforeAll(() => console.log('> before all'))
beforeEach((testName) => console.log(`>> before test "${testName}"`))
afterEach((testName) => console.log(`<< after test "${testName}"`))
afterAll(() => console.log('< after all'))

describe('boolean values', () => {
  test('true', () => assert.true(true))

  it('false', () => {
    mustBe.false(false)
  })
})

describe('strings', () => {
  test('equals', () => assert.equals('bar', 'bar'))

  describe('sub-strings', () => {
    test('length', () => mustBe.true('foo'.length === 3))
    test('numbers', () => mustBe.equals(1, 1))

    for (const value of [{}, [], ""]) {
      assert.empty(value)
    }

    for (const value of [{foo: 'bar'}, [1], "string"]) {
      mustBe.notEmpty(value)
    }

    assert.contains('foo', 'foobar')
    assert.notContains('foo', 'barbaz')

    assert.contains(true, [3, 2, 1, true])
    assert.notContains(2, ['foo', 'bar'])
  })
})
