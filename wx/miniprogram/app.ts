import { getSetting, getUserInfo } from "./utils/util"

// app.ts
App<IAppOption>({
  globalData: {},
  onLaunch() {
    // 展示本地存储能力
    const logs = wx.getStorageSync('logs') || []
    logs.unshift(Date.now())
    wx.setStorageSync('logs', logs)

    // 登录
    wx.login({
      success: res => {
        console.log(res.code)
        wx.request({
          url:'http://localhost:8200/v1/auth/login',
          method:'POST',
          data:{
            code:res.code
          },
          success:res => {
              console.log(res)
              wx.request({
                url:'http://localhost:8200/v1/todo',
                method:'POST',
                data:{
                  title:'test'
                },
                header: {
                  authorization: 'Bearer ' + res.data.access_token
                },
                success:console.log,
                fail:console.log,
              })
          },
          fail:console.log,
        })
        // 发送 res.code 到后台换取 openId, sessionKey, unionId
      },
    })
    // getSetting().then(res => {
    //   if (res.authSetting['scope.userInfo']){
    //     return getUserInfo()
    //   }
    //   // return Promise.resolve(undefined)
    //   return undefined
    // }).then(res => {
    //   if (!res) {
    //     return 
    //   }
    //   this.globalData.userInfo = res?.userInfo
    //   //通知页面我获得了用户信息
    //   if(this.userInfoReadyCallback){
    //     this.userInfoReadyCallback(res)
    //   }
    // })
  },
})