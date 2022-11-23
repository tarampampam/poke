#!/usr/bin/env poke run

console.log(JSON.stringify({
  bool: faker.bool(), // true, false
  falsy: faker.falsy(), // false, null, undefined...
  character: faker.character(), // L
  characterPool: faker.character({pool: 'abc'}), // b
  floating: faker.floating(), // -1913.91796875
  integer: faker.integer(), // 2146049948029865
  integerMinMax: faker.integer({min: 1, max: 10}), // 6
  letter: faker.letter(), // s
  string: faker.string(), // asFlOzy2NoT
  stringPool: faker.string({pool: 'abc'}), // acbaacaaacb
  stringLength: faker.string({length: 64}), // kJ3wlJaA4sUnoyj94AWi1Vt9w1y0g91XzzZs9N0VIVyptrEwL4I4RJuC67VPYyA7
  paragraph: faker.paragraph(), // Aut consequatur sit accusantium perferendis voluptatem.
  word: faker.word(), // quas
  domain: faker.domain(), // iFDyBwW.com
  email: faker.email(), // FBJCQDE@ZpdvrIX.info
  ip: faker.ip(), // 112.245.6.64
  ipv6: faker.ipv6(), // 283f:d66b:30d1:a435:48c3:1f:45e2:a503
  tld: faker.tld(), // org
  url: faker.url(), // http://vjkNNvo.biz/VPpxCkM.php"
  date: faker.date(), // 2009-09-04T07:35:01.209Z
  hash: faker.hash(), // 3ab26351d4c9309899075dd8a854b916ab7ac398
  hashLength: faker.hash({length: 8}), // 4ca57fc7
  uuid: faker.uuid(), // 296fdaf7-9ac0-446a-a59e-fcac325b3c2d
}, undefined, 2))
