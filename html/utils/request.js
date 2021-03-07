import axios from 'axios';

function request(options) {
  return axios(options)
    .then(res => {
      const { data } = res;
      if (data.code === 0) {
        return data;
      } else {
        return {
          erred: true,
          message: data.message || '',
        };
      }
    })
    .catch(error => {
      return {
        erred: true,
        message: error.message || '',
      };
    });
}

export { request };
