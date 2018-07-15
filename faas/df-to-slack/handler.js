'use strict';

const fs = require('fs');
const { WebClient } = require('@slack/client');

module.exports = (context, callback) => {
  const oAuthToken = fs
    .readFileSync(`/run/secrets/bot-user-oauth-access-token`)
    .toString();

  const { data } = JSON.parse(context);
  const body = Buffer.from(data.body, 'base64');

  const query = body.queryResult.queryText;
  const web = new WebClient(oAuthToken);
  web.chat
    .postMessage({ channel: 'C994KQKG9', text: query })
    .then(res => {
      // `res` contains information about the posted message
      console.error('Message sent: ', res.ts);

      callback(undefined, {
        body: {
          payload: {
            google: {
              expectUserResponse: false,
              richResponse: {
                items: [
                  {
                    simpleResponse: {
                      textToSpeech: '送りました'
                    }
                  }
                ]
              }
            }
          }
        }
      });
    })
    .catch(err => {
      console.error(err);
      callback(undefined, {
        body: {
          payload: {
            google: {
              expectUserResponse: false,
              richResponse: {
                items: [
                  {
                    simpleResponse: {
                      textToSpeech: '送信に失敗しました'
                    }
                  }
                ]
              }
            }
          }
        }
      });
    });
};
