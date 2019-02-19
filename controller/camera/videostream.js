const util = require('util');
const Duplex = require('stream').Duplex;

module.exports = VideoStream;

function VideoStream(byteLength, options) {
  this.bytes = [];

  if (byteLength) {
    this.size = byteLength;
  } else {
    this.size = 1024
  }

  if (!(this instanceof VideoStream)) {
    return new VideoStream(options);
  } else {
    if (!options) options = {}; // ensure object
    Duplex.call(this, options);
  }
}
util.inherits(VideoStream, Duplex);

VideoStream.prototype._read = function() {
  if (this.bytes >= this.size) {
    this.push(this.bytes)
    this.bytes.splice(0,1023)
  }
}

VideoStream.prototype._write = function(bytes, enc, callback) {
  this.bytes.push(bytes);
  callback();
}
