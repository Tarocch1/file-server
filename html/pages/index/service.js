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
