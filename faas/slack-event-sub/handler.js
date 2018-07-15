'use strict';

const fs = require('fs');
const { WebClient } = require('@slack/client');

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

  if (data.headers["X-Slack-Retry-Num"] && data.headers["X-Slack-Retry-Reason"] === "http_timeout") {
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

  const web = new WebClient(oAuthToken);
  web.chat
    .postMessage({ channel: body.event.channel, text: 'Hello there' })
    .then(res => {
      // `res` contains information about the posted message
      console.error('Message sent: ', res.ts);
    })
    .catch(console.error);
};
