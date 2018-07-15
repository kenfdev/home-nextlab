"use strict"

module.exports = (context, callback) => {
    console.error(context);
    callback(undefined, {status: "done"});
}
