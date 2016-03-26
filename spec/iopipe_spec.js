var iopipe = require("../js/iopipe")()

describe("define", function() {
  it("returns a function", function() {
    var output = iopipe.define(function() { })
    expect(typeof output).toEqual("function")
  })
})

describe("defined-function", function() {
  it("can pass no args", function(done) {
    var fun = iopipe.define(function(i, ctx) {
      expect(i).toEqual(undefined)
      ctx.done()
      done()
    })
    fun()
  })
  it("can pass one arg", function(done) {
    var expected = "hello world"
    var fun = iopipe.define(function(i, ctx) {
      expect(i).toEqual(expected); ctx.done()
    }, done)
    fun(expected)
  })
  it("context is callback", function(done) {
    var fun = iopipe.define(function(_, ctx) { ctx() })
    fun(undefined, done)
  })
  it("can trigger context.done", function(done) {
    var fun = iopipe.define(function(_, ctx) { ctx.done() })
    fun(undefined, done)
  })
  it("can trigger context.succeed", function(done) {
    var fun = iopipe.define(function(_, ctx) { ctx.succeed() })
    fun(undefined, done)
  })
  it("can trigger context.raw", function(done) {
    var fun = iopipe.define(function(_, ctx) { ctx.raw() })
    fun(undefined, done)
  })
  it("passes result to callback", function(done) {
    var input = 2
    var expected = 3
    var fun = iopipe.define(function(i, ctx) {
      ctx.done(i + 1)
    })
    fun(input, function(i) { 
      expect(i).toEqual(expected)
      done()
    })
  })
})

describe("map", function() {
  it("has as many outputs as inputs", function(done) {
    var input = [0, 1, 2]
    iopipe.map(function(i, cxt) { cxt.done(i + 2) })(input, function(output) {
      expect(input.length).toEqual(output.length);
      done()
      //ctx()
    })
  });
  it("preserves order", function(done) {
    var input = [0, 1, 2]
    iopipe.map(function(i, ctx) { ctx.done(i) })(input, function(output, ctx) {
      expect(output).toEqual(input);
      done()
      //ctx()
    })
  });
  it("transforms each input element", function(done) {
    var input = [0, 1, 2]
    var expected = [1, 2, 3]
    iopipe.map(function(i, ctx) { ctx.done(i + 1) })(input, function(output, ctx) {
      expect(output).toEqual(expected)
      done()
      //ctx()
    })
  });
});

describe("tee", function() {
  it("has as many outputs as functions", function(done) {
    var input = [0, 1, 2, 3, 4]
    var echo = function(i, ctx) { ctx.done(i) }
    iopipe.tee(echo, echo)(input, iopipe.make_context(function(output) {
      expect(output.length).toEqual(2);
      done()
    }))
  });
  it("preserves order", function(done) {
    var input = [0, 1, 2]
    var echo = function(i, ctx) { ctx(i) }
    var ret2 = function(i, ctx) { ctx(2) }
    iopipe.tee(echo, ret2)(input, iopipe.make_context(function(output) {
      echo(input, function(e) {
        ret2(input, function(r) {
          expect(output).toEqual([e, r])
          done()
        })
      })
    }))
  });
});

describe("reduce", function() {
  it("can sum all input elements", function(done) {
    var input = [1, 2, 3]
    var sum = function(prev, next) {
      return prev + next
    }
    iopipe.reduce(sum)(input, function(output, ctx) {
      expect(output).toEqual(6)
      done()
      //ctx.done()
    })
  });
})

describe("exec", function() {
  it("can chain functions", function(done) {
    iopipe.exec(function(_, ctx) { ctx.done("hello world") }
                ,function(input, ctx) {
                  expect(input).toEqual("hello world")
                  done()
                  ctx.done()
                }
                ,function(input, ctx) { done(); ctx.done() })
  });
})

describe("apply", function() {
  it("executes function", function(done) {
    iopipe.apply()(done)
  });
})

describe("property", function() {
  it("returns property for arg", function(done) {
    var obj = { "key": "hello world" }
    iopipe.property("key")(obj, function (output) {
      expect(output).toEqual(obj["key"])
      done()
    })
  });
})

describe("callback", function() {
  it("calls function", function(done) {
    iopipe.callback(done)()
  })
  it("passes input to function", function(done) {
    iopipe.callback(function(input) {
      expect(input).toEqual("hello world")
      done()
    })("hello world")
  })
})
