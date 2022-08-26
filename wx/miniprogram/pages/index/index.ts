// index.ts
// 获取应用实例
const app = getApp<IAppOption>()

Page({
    data: {
      setting :{
          skew: 0,
          rotate: 0,
          showLocation: true,
          showScale: true,
          subKey: '',
          layerStyle: -1,
          enableZoom: true,
          enableScroll: true,
          enableRotate: false,
          showCompass: false,
          enable3D: false,
          enableOverlooking: false,
          enableSatellite: false,
          enableTraffic: false,
      },
      location: {
          latitude: 31.240232,
          longitude: 121.318533
      },
      scale: 14,
      isOverLooking: false,
      is3D: true,
      minScale: 3,
      maxScale: 20,
      markers: [
        {
          id: 1,
          latitude: 31.240258,
          longitude: 121.320411,
          width: 30,
          height: 30,
          iconPath:  '/resources/sedan.png'
        },
        {
          id: 2,
          latitude: 31.240946,
          longitude: 121.316476,
          width: 30,
          height: 30,
          iconPath:  '/resources/sedan.png'
        },
        {
          id: 3,
          latitude: 31.241201,
          longitude: 121.305924,
          width: 30,
          height: 30,
          iconPath:  '/resources/sedan.png'
        },
    ],
    },
    onMyLocationTap() {
      wx.getLocation({
        type: 'gcj02',
        success: res => {
          this.setData({
            location: {
              latitude: res.latitude,
              longitude: res.longitude,
            },
          })
        }, 
        fail: () => {
          wx.showToast({
            icon: 'none',
            title: '请前往设置页授权',
          })
        }
      })
    },
})
