'use strict';

const fs = require('fs');
const { WebClient } = require('@slack/client');
const request = require('request');

module.exports = (context, callback) => {
  callback(undefined, {
    headers: { 'Content-Type': 'text/plain' },
    body: 'success'
  });

  const oAuthToken = fs
    .readFileSync(`/run/secrets/bot-user-oauth-access-token`)
    .toString();
  const verifyToken = fs
    .readFileSync(`/run/secrets/slack-verify-token`)
    .toString();

  const {
    data: { body }
  } = JSON.parse(context);

  if (
    data.headers['X-Slack-Retry-Num'] &&
    data.headers['X-Slack-Retry-Reason'] === 'http_timeout'
  ) {
    // Do nothing if it's a retry.
    return;
  }

  if (body.token !== verifyToken) {
    // callback(undefined, {
    //   statusCode: 400,
    //   headers: {
    //     'Content-Type': 'text/plain'
    //   },
    //   body: 'Token Invalid'
    // });
    // return;
  }

  const url = process.env.GOOGLE_HOME_NOTIFIER_FUNC_URL;
  request.post(url, { text: body.event.text, language: 'ja' }, (err, resp) => {
    let msg = 'success';
    if (err) {
      msg = 'Error sending message.';
      return;
    }

    const web = new WebClient(oAuthToken);
    web.chat
      .postMessage({ channel: body.event.channel, text: msg })
      .then(res => {
        // `res` contains information about the posted message
        console.error('Message sent: ', res.ts);
      })
      .catch(console.error);
  });
};
