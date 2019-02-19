const fs = require('fs');
const splitFileStream = require('split-file-stream');
const basePath = '/storage/videoFiles';
const fileSize = 2048; //Save files in 2GB chunks

module.exports = CameraStorage;

function CameraStorage(cam) {
  let self = this;
  this.cam = cam;
  this.storagePath = basePath + cam.id;

  // Create the camera's storage path if it does not exist
  if (!fs.existsSync(self.storagePath)) {
    fs.mkdirSync(self.storagePath, 0755);
  }

  // Whenever there's data coming through the readStream, write it to disk.
  // When the file reaches a certain size, split it into a new file.
  splitFileStream.split(self.cam.videoStream, fileSize, this.storagePath, (filePaths) => {
    logger.log('debug', "Wrote video files for camera " + this.cam.name + ": " + filePaths);
  });
}
