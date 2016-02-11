var iopipe = require("../js/iopipe")

describe("define", function() {
  it("returns a function", function() {
    var output = iopipe.define(function() { })
    expect(typeof output).toEqual("function")
  })
})

describe("defined-function", function() {
  it("can pass no args", function(done) {
    var fun = iopipe.define(function(i) {
      expect(i).toEqual(undefined)
      done()
    })
    fun()
  })
  it("can pass one arg", function(done) {
    var expected = "hello world"
    var fun = iopipe.define(function(i) {
      expect(i).toEqual(expected)
    }, done)
    fun(expected)
  })
  /* depends on upcoming context patch *
  it("can trigger callback", function(done) {
    var fun = iopipe.define(function() { })
    fun("", done)
  })
  it("passes result to callback", function(done) {
    var input = 2
    var expected = 3
    var fun = iopipe.define(function(i) {
      return i + 1
    })
    fun(input, function(i) { 
      expect(i).toEqual(expected)
      done()
    })
  })
  */
})

describe("map", function() {
  it("has as many outputs as inputs", function() {
    var input = [0, 1, 2]
    var output = iopipe.map(function(i) { return i + 2 })(input)
    expect(input.length).toEqual(output.length);
  });
  it("preserves order", function() {
    var input = [0, 1, 2]
    var output = iopipe.map(function(i) { return i })(input)
    expect(output).toEqual(input);
  });
  it("transforms each input element", function() {
    var input = [0, 1, 2]
    var expected = [1, 2, 3]
    var output = iopipe.map(function(i) { return i + 1 })(input)
    expect(output).toEqual(expected)
  });
});

describe("tee", function() {
  it("has as many outputs as functions", function() {
    var input = [0, 1, 2, 3, 4]
    var echo = function(i) { return i }
    var output = iopipe.tee(echo, echo)(input)
    expect(output.length).toEqual(2);
  });
  it("preserves order", function() {
    var input = [0, 1, 2]
    var echo = function(i) { return i }
    var ret2 = function(i) { return 2 }
    var output = iopipe.tee(echo, ret2)(input)
    expect(output).toEqual([ echo(input), ret2(input) ]);
  });
});

describe("reduce", function() {
  it("can sum all input elements", function() {
    var input = [1, 2, 3]
    var sum = function(prev, next) {
      return prev + next
    }
    var output = iopipe.reduce(sum)(input)
    expect(output).toEqual(6)
  });
})

describe("exec", function() {
  it("can chain functions", function(done) {
    iopipe.exec(function() { return "hello world" }
                ,function(input) {
                  expect(input).toEqual("hello world")
                }
                ,done)
  });
})

describe("apply", function() {
  it("executes function", function(done) {
    iopipe.apply()(done)
  });
})

describe("property", function() {
  it("returns property for arg", function() {
    var obj = { "key": "hello world" }
    var output = iopipe.property("key")(obj)
    expect(output).toEqual(obj["key"])
  });
})
