import { expect as expectCDK, matchTemplate, MatchStyle } from '@aws-cdk/assert';
import * as cdk from '@aws-cdk/core';
import * as MySlackWatcher from '../lib/my-slack-watcher-stack';

test('Empty Stack', () => {
    const app = new cdk.App();
    // WHEN
    const stack = new MySlackWatcher.MySlackWatcherStack(app, 'MyTestStack');
    // THEN
    expectCDK(stack).to(matchTemplate({
      "Resources": {}
    }, MatchStyle.EXACT))
});
