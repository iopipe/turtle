var agent = require('..')

describe('metrics agent', () => {
  it('should return a function', () => {
    var wrapper = agent()
    expect(typeof wrapper).toEqual('function')
  })
})