'use strict'

function reduce(json, code) {
  if (/^\./.test(code)) {
    const fx = eval(`function fn() { 
      return ${code === '.' ? 'this' : 'this' + code} 
    }; fn`)
    return fx.call(json)
  }

  if ('?' === code) {
    return Object.keys(json)
  }

  if (/yield\*?\s/.test(code)) {
    const fx = eval(`function fn() {
      const gen = (function*(){ 
        ${code.replace(/\\\n/g, '')} 
      }).call(this)
      return [...gen]
      }; fn`)
    return fx.call(json)
  }

  const fx = eval(`function fn() { 
    return ${code} 
  }; fn`)

  const fn = fx.call(json)
  if (typeof fn === 'function') {
    return fn(json)
  }
  return fn
}

module.exports = reduce
