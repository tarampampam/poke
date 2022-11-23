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

    for (const value of [
      [],
      {},
      false,
      null,
      undefined,
      "",
      0,
      NaN,
      new Map(),
      new Set(),
      new Date('invalid'),
    ]) {
      assert.empty(value)
    }

    for (const value of [
      [1],
      {foo: 'bar'},
      [undefined],
      1,
      "string",
      true,
      Symbol(),
      () => {},
      new Date(),
      new Promise(() => {})
    ]) {
      assert.notEmpty(value)
    }
  })
})
