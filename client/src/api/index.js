var axios = require('axios');
var MockAdapter = require('axios-mock-adapter');
var normalAxios = axios.create();
var mockAxios = axios.create();

// mock 数据
var mock = new MockAdapter(mockAxios);

const machine = [
  { "hostname": "cghost3", "IP": "192.168.1.252", "outIP": "192.168.1.252", "applications": ["zone1"], "online": true, "codeVersion": "38715" }
]


mock.onPut('/login').reply(config => {
  let postData = JSON.parse(config.data).data;
  if (postData.user === 'admin' && postData.password === '123456') {
    return [200, require('./mock/user')];
  } else {
    return [500, { message: "Incorrect user or password" }];
  }
});
mock.onGet('/logout').reply(200, {});
mock.onGet('/my').reply(200, require('./mock/user'));
mock.onGet('/menu').reply(200, require('./mock/menu'));
mock.onGet('/randomuser').reply((config) => {
  return new Promise(function (resolve, reject) {
    normalAxios.get('https://randomuser.me/api', {
      params: {
        results: 10,
        ...config.params,
      },
      responseType: 'json'
    }).then((res) => {
      resolve([200, res.data]);
    }).catch((err) => {
      resolve([500, err]);
    });
  });
});

mock.onGet('/machine').reply(200, machine)
mock.onPost('/machine/add').reply(config => {
  let postData = JSON.parse(config.data).machine;
  machine.push(postData)
  return [200, machine]
})

mock.onPost('/machine/save').reply(config => {
  let postData = JSON.parse(config.data).machine;
  let oldpostData = JSON.parse(config.data).oldmachine;
  machine.forEach((element, index) => {
    if (element.hostname == oldpostData.hostname) {
      machine[index] = postData
    }
  })
  return [200, machine]
})

mock.onPost('/machine/del').reply(config => {
  let postData = JSON.parse(config.data);
  machine.forEach((element, index) => {
    if (element.hostname == postData.hostname) {
      machine.splice(index, 1)
    }
  })
  return [200, machine]
})

//if (process.env.NODE_ENV === 'production') {
//  //export default mockAxios;
//  module.exports = mockAxios;
//} else {
//  module.exports = normalAxios;
//}

export default normalAxios;