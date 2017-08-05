module.exports = {
    menus: [
      {
        key: 5,
        name: 'Home',
        icon: 'home',
        url: '/home'
      },
      {
        key: 1,
        name: 'Pages',
        icon: 'user',
        child: [
        {
            name: 'Machine',
            key: 101,
            url: '/machine'
          },
          {
            name: 'Zone',
            key: 102,
            url: '/zone'
          }
 
       ]
      },
      {
        key: 2,
        name: 'Components',
        icon: 'laptop',
        child: [
       ]
      },
    ]
  }
