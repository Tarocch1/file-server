import axios from 'axios';
import { startUpload, uploadChunk, endUpload } from './service';

const CancelToken = axios.CancelToken;
const defaultChunkSize = 10 * 1024 * 1024;
const defaultThread = 3;

function Uploader(options) {
  this.options = options;
  this.chunks = [];
  this.uploadingChunks = 0; // 正在上传的分块数量
  this.curChunk = 0; // 分块指针
  this.cancelTokens = [];
  this._createChunks();
  return this;
}

Uploader.prototype._createChunks = function () {
  const size = this.options.chunkSize || defaultChunkSize;
  let cur = 0;
  while (cur < this.options.file.size) {
    this.chunks.push({
      file: this.options.file.slice(cur, cur + size),
      loaded: 0,
    });
    cur += size;
  }
};

Uploader.prototype.start = function () {
  startUpload(this.options.path).then(res => {
    if (!res.erred) {
      for (let i = 0; i < (this.options.thread || defaultThread); i++) {
        this._uploadChunk();
      }
    } else {
      typeof this.options.onError === 'function' &&
        this.options.onError(res.message);
    }
  });
};

Uploader.prototype._uploadChunk = function () {
  if (this.uploadingChunks >= (this.options.thread || defaultThread)) return;
  if (this.curChunk >= this.chunks.length) {
    if (this.uploadingChunks === 0) {
      this._end();
    }
    return;
  }
  const id = this.curChunk;
  const chunk = this.chunks[id];
  const cancelSource = CancelToken.source();
  this.cancelTokens.push(cancelSource);
  uploadChunk(this.options.path, id, chunk.file, cancelSource.token, e => {
    console.log(e);
    this.chunks[id].loaded = e.loaded;
    this._onProgress();
  }).then(res => {
    if (!res.erred) {
      this._onProgress();
      this.uploadingChunks--;
      this._uploadChunk();
    } else {
      this._cancel();
      typeof this.options.onError === 'function' &&
        this.options.onError(res.message);
    }
  });
  this.uploadingChunks++;
  this.curChunk++;
};

Uploader.prototype._end = function () {
  endUpload(this.options.path).then(res => {
    if (!res.erred) {
      typeof this.options.onSuccess === 'function' && this.options.onSuccess();
    } else {
      typeof this.options.onError === 'function' &&
        this.options.onError(res.message);
    }
  });
};

Uploader.prototype._cancel = function () {
  this.cancelTokens.forEach(source => {
    source.cancel();
  });
};

Uploader.prototype._onProgress = function () {
  const loaded = this.chunks.reduce((acc, chunk) => acc + chunk.loaded, 0);
  typeof this.options.onProgress === 'function' &&
    this.options.onProgress({ loaded });
};

export default Uploader;
