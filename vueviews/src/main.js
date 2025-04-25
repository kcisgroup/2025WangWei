import Vue from 'vue'
import App from './App.vue'
import axios from "axios"
import router from '@/router/index'
import ElementUI from 'element-ui'
import 'element-ui/lib/theme-chalk/index.css'

Vue.config.productionTip = false

//全局配置axios的请求根路径
axios.defaults.baseURL='http://47.100.71.46:8088'
//将axios挂载到Vue.prototype上，供全部vue组件实例直接使用
Vue.prototype.$http=axios
//elementUI配置
Vue.use(ElementUI)

new Vue({
  render: h => h(App),
  router
}).$mount('#app')
