var iopipe = require("../js/iopipe")

describe("define", function() {
  it("returns a function", function() {
    var output = iopipe.define(function() { })
    expect(typeof output).toEqual("function")
  })
})

describe("defined-function", function() {
  it("can pass no args", function(done) {
    var fun = iopipe.define(function(i, cb) {
      expect(i).toEqual(undefined)
      cb()
      done()
    })
    fun()
  })
  it("can pass one arg", function(done) {
    var expected = "hello world"
    var fun = iopipe.define(function(i, cb) {
      expect(i).toEqual(expected); cb()
    }, done)
    fun(expected)
  })
  it("can trigger callback", function(done) {
    var fun = iopipe.define(function(_, cb) { cb() })
    fun("", done)
  })
  it("passes result to callback", function(done) {
    var input = 2
    var expected = 3
    var fun = iopipe.define(function(i, cb) {
      cb(i + 1)
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
    iopipe.map(function(i, cb) { cb(i + 2) })(input, function(output) {
      expect(input.length).toEqual(output.length);
      done()
    })
  });
  it("preserves order", function(done) {
    var input = [0, 1, 2]
    iopipe.map(function(i, cb) { cb(i) })(input, function(output) {
      expect(output).toEqual(input);
      done()
    })
  });
  it("transforms each input element", function(done) {
    var input = [0, 1, 2]
    var expected = [1, 2, 3]
    iopipe.map(function(i, cb) { cb(i + 1) })(input, function(output) {
      expect(output).toEqual(expected)
      done()
    })
  });
});

describe("tee", function() {
  it("has as many outputs as functions", function(done) {
    var input = [0, 1, 2, 3, 4]
    var echo = function(i, cb) { cb(i) }
    iopipe.tee(echo, echo)(input, function(output) {
      expect(output.length).toEqual(2);
      done()
    })
  });
  it("preserves order", function(done) {
    var input = [0, 1, 2]
    var echo = function(i, cb) { cb(i) }
    var ret2 = function(i, cb) { cb(2) }
    iopipe.tee(echo, ret2)(input, function(output) {
      echo(input, function(e) {
        ret2(input, function(r) {
          expect(output).toEqual([e, r])
          done()
        })
      })
    })
  });
});

describe("reduce", function() {
  it("can sum all input elements", function(done) {
    var input = [1, 2, 3]
    var sum = function(prev, next) {
      return prev + next
    }
    iopipe.reduce(sum)(input, function(output) {
      expect(output).toEqual(6)
      done()
    })
  });
})

describe("exec", function() {
  it("can chain functions", function(done) {
    iopipe.exec(function(_, cb) { cb("hello world") }
                ,function(input, cb) {
                  expect(input).toEqual("hello world")
                  done()
                  cb()
                }
                ,function(input, cb) { done(); cb() })
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
