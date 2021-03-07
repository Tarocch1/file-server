import { request } from '../../utils/request';

export function getList(path) {
  return request({
    url: '/api/list',
    method: 'post',
    data: {
      path,
    },
  });
}

export function remove(path) {
  return request({
    url: '/api/remove',
    method: 'post',
    data: {
      path,
    },
  });
}

export function startUpload(path) {
  return request({
    url: '/api/upload/start',
    method: 'post',
    data: {
      path,
    },
  });
}

export function uploadChunk(path, id, data, cancelToken, onUploadProgress) {
  return request({
    url: `/api/upload/chunk?path=${path}&id=${id}`,
    method: 'post',
    data,
    cancelToken,
    onUploadProgress,
  });
}

export function endUpload(path) {
  return request({
    url: '/api/upload/end',
    method: 'post',
    data: {
      path,
    },
  });
}
