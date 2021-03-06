#!/usr/bin/env node
import 'source-map-support/register';
import * as cdk from '@aws-cdk/core';
import { MySlackWatcherStack } from '../lib/my-slack-watcher-stack';

const app = new cdk.App();
new MySlackWatcherStack(app, 'MySlackWatcherStack');
