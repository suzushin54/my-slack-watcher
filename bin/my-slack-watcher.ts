#!/usr/bin/env node
import 'source-map-support/register';
import * as cdk from '@aws-cdk/core';
import {MySlackWatcherStack} from '../lib/my-slack-watcher-stack';

const util = require('util')
const exec = util.promisify(require('child_process').exec)

async function deploy() {
  // Build
  await exec('go get -v -t -d ./lambdaSource/... && ' +
    'GOOS=linux GOARCH=amd64 ' +
    'go build -o ./lambdaSource/main ./lambdaSource/**.go')

  const app = new cdk.App();
  new MySlackWatcherStack(app, 'MySlackWatcherStack')
  app.synth()

  // clean
  await exec('rm -f ./lambdaSource/main')
}

deploy()
