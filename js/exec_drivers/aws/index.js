var AWS = require('aws-sdk')
var config = require('../config')
var crypto = require('crypto')
var ctxutils = require('../../ctxutils')

var awsLambda = new AWS.Lambda({apiVersion: '2015-03-31'});

exports.init = function() {
  /* Init */
  var cfgbase = config.exec_driver.aws_lambda
  AWS.config.region = cfgbase.region
  AWS.config.update({accessKeyId: cfgbase.access_key,
                     secretAccessKey: cfgbase.secret_key})
}

exports.invoke = function(event, context) {
  var funcID = event.id
  // api docs: http://docs.aws.amazon.com/AWSJavaScriptSDK/latest/AWS/Lambda.html#invoke-property
  return function(prevResult) {
    var params = {
      FunctionName: funcID
      ,Payload: prevResult
    }
    awsLambda.invoke(params, ctxutils.callback(context))
  }
}

exports.listFunctions = function(event, context) {
  awsLambda.listFunctions({}, ctxutils.callback(context))
}

exports.getFunction = function(event, context) {
  // api docs: http://docs.aws.amazon.com/AWSJavaScriptSDK/latest/AWS/Lambda.html#getFunction-property
  awsLambda.getFunction(params, ctxutils.callback(context))

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

exports.createFunction = function(event, context) {
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
