"use strict"

module.exports = (context, callback) => {
    console.error(context);
    callback(undefined, {headers: {"Content-Type": "text/plain"}, body: "success"});
}
