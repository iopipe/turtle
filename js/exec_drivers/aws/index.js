var AWS = require('aws-sdk')
var ctxutils = require('../../ctxutils')

module.exports = function(opts) {
  return new LambdaDriver(opts)
}

function LambdaDriver(opts) {
  this._aws_lambda = new AWS.Lambda({
    apiVersion: '2015-03-31',
    region: opts.region,
    accessKeyId: opts.access_key,
    secretAccessKey: opts.secret_key
  })
}

LambdaDriver.prototype.invoke = function(event, context) {
  var driver = this._aws_lambda
  var funcID = event.id
  // api docs: http://docs.aws.amazon.com/AWSJavaScriptSDK/latest/AWS/Lambda.html#invoke-property
  return function(prevResult) {
    var params = {
      FunctionName: funcID
      ,Payload: prevResult
    }
    driver.invoke(params, ctxutils.callback(context))
  }
}

LambdaDriver.prototype.listFunctions = function(event, context) {
  this._aws_lambda.listFunctions({}, ctxutils.callback(context))
}

LambdaDriver.prototype.getFunction = function(event, context) {
  // api docs: http://docs.aws.amazon.com/AWSJavaScriptSDK/latest/AWS/Lambda.html#getFunction-property
  this._aws_lambda.getFunction(params, ctxutils.callback(context))

  /* data:
   configuration {
     FunctionName: string
     FunctionArn: string
     Runtime: string
     Role: string
     CodeSize: integer (bytes)
     Description: string
     Timeout: integer (secs)
     MemorySize: integer (64mb chunks)
     LastMOdified: string
     CodeSha256: string (sha256 hash of the function deployment package)
     VpcConfig: {} 
   }
   code {
     RepositoryType: string
     Location: string (presigned URL, valid for 10 minutes)
   }
  */
}

LambdaDriver.prototype.createFunction = function(event, context) {
  /*params = {
    FunctionName:
    ,Role:
    ,Handler: "exports.handler"
    ,Code: {
      Zipfile: zip_data
      ,S3Bucket:
      ,S3Key
      ,S3ObjectVersion: ??
    }
    ,Runtime: "nodejs"
  }*/
}
