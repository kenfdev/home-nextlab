'use strict';

const fs = require('fs');
const { WebClient } = require('@slack/client');

module.exports = (context, callback) => {
  callback(undefined, {
    body: {
      payload: {
        google: {
          expectUserResponse: false,
          richResponse: {
            items: [
              {
                simpleResponse: {
                  textToSpeech: '受け付けました'
                }
              }
            ]
          }
        }
      }
    }
  });
  const oAuthToken = fs
    .readFileSync(`/run/secrets/bot-user-oauth-access-token`)
    .toString();

  const { data } = JSON.parse(context);
  const body = JSON.parse(Buffer.from(data.body, 'base64'));
  console.error('body', body);

  const query = body.queryResult.queryText;
  const web = new WebClient(oAuthToken);
  web.chat
    .postMessage({ channel: 'C994KQKG9', text: query })
    .then(res => {
      // `res` contains information about the posted message
      console.error('Message sent: ', res.ts);
    })
    .catch(err => {
      console.error(err);
    });
};
