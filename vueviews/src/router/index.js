import Vue from 'vue'
import VueRouter from 'vue-router'
import Login from "@/components/Login";
import Home from "@/components/front/Home";
import AdminHome from "@/components/rear/AdminHome";
import MyCips from "@/components/front/cipManage/MyCips";
import UserHome from "@/components/front/UserHome";
import CipsOpened from "@/components/front/componentSub/CipsOpened";
import UploadDraft from "@/components/front/cipManage/DraftUpload";
import DraftBox from "@/components/front/cipManage/DraftBox";
import DraftDetails from "@/components/front/cipManage/DraftDetails";
import CipOpenedDetails from "@/components/front/componentSub/CipOpenedDetails";
import MyCipDetails from "@/components/front/cipManage/MyCipDetails";
import SolveSub from "@/components/front/componentSub/SolveSub";
import ApplyForCip from "@/components/front/cipManage/ApplyForCip";
import ApplySteps from "@/components/front/cipManage/ApplySteps";
import AchieveStatistical from "@/components/front/AchieveStatistical";
import OrderStatistical from "@/components/front/OrderStatistical";
import CartList from "@/components/front/componentSub/CartList";
import SubscribeCip from "@/components/front/componentSub/SubscribeCip";
import MySubs from "@/components/front/componentSub/MySubs";
import MySubCipDetails from "@/components/front/componentSub/MySubCipDetails";
import Solve from "@/components/front/componentSub/Solve";
import ConfirmSub from "@/components/front/componentSub/ConfirmSub";
import ServiceCreate from "@/components/front/serviceUse/ServiceCreate";
import ServiceList from "@/components/front/serviceUse/ServiceList";
import ServiceDestroy from "@/components/front/serviceUse/ServiceDestroy.vue";
import ServiceDetails from "@/components/front/serviceUse/ServiceDetails.vue";
import UserDetails from "@/components/front/userInfo/UserDetails.vue";
import UserUpdate from "@/components/front/userInfo/UserUpdate.vue";

//调用Vue.user()将VueRouter安装为Vue的插件
Vue.use(VueRouter)

//创建路由的实例对象
const router=new VueRouter({
    routes:[
        {
            path:'/',   // 重定向到登录组件
            redirect:'/login'
        },
        {
            path:'/login',   // 跳转登录组件（该组件包含注册功能，抽屉形式）
            name:'Login',
            component:Login
        },
        {
            path:'/home', // 首页的导航菜单组件
            name:'Home',
            redirect: '/userhome',
            component:Home,
            children:[
                {
                    path:'/userhome',
                    name:'UserHome',
                    component:UserHome
                },
                {
                    path:'/draft_upload', // 成果上传
                    name:'UploadDraft',
                    component:UploadDraft
                },
                {
                    path:'/draft_box', // 草稿箱
                    name:'DraftBox',
                    component:DraftBox
                },
                {
                    path:'/draft_details/:id', // 项目初稿详情
                    name:'DraftDetails',
                    component:DraftDetails,
                    props: true  // 启用将路由参数传递为组件的 props
                },
                {
                    path:'/apply_for_cip', // 申请代码产权
                    name:'ApplyForCip',
                    component:ApplyForCip,
                },
                {
                    path:'/apply_steps', // 代码产权申请流程
                    name:'ApplySteps',
                    component:ApplySteps
                },
                {
                    path:'/my_cips', // 查验代码产权
                    name:'MyCips',
                    component:MyCips
                },
                {
                    path:'/my_cip_details/:cpi', // 产权详情
                    name:'MyCipDetails',
                    component:MyCipDetails,
                    props: true
                },
                {
                    path:'/cips_opened', // 浏览开放组件
                    name:'CipsOpened',
                    component:CipsOpened
                },
                {
                    path:'/cip_opened_details/:id', // 组件详情
                    name:'CipOpenedDetails',
                    component:CipOpenedDetails,
                    props:true
                },
                {
                    path:'/cart_list', // 我的购物车
                    name:'CartList',
                    component:CartList,
                    props:true
                },
                {
                    path:'/subscribe_cip/:id', // 订阅申请表
                    name:'SubscribeCip',
                    component:SubscribeCip,
                    props:true
                },
                {
                    path:'/my_subs', // 我的订阅记录
                    name:'MySubs',
                    component:MySubs
                },
                {
                    path:'/my_sub_cip_details/:id', // 我订阅的组件的详情
                    name:'MySubCipDetails',
                    component:MySubCipDetails,
                    props:true
                },
                {
                    path:'/confirm_sub/:id', // 我订阅的组件的详情
                    name:'ConfirmSub',
                    component:ConfirmSub,
                    props:true
                },
                {
                    path:'/solve_sub', // 我的订单，待处理订单列表
                    name:'SolveSub',
                    component:SolveSub
                },
                {
                    path:'/solve/:id', // 处理订单
                    name:'Solve',
                    component:Solve,
                    props:true
                },
                {
                    path:'/service_destroy', // 销毁服务
                    name:'ServiceDestroy',
                    component:ServiceDestroy
                },
                {
                    path:'/service_list', // 服务列表
                    name:'ServiceList',
                    component:ServiceList
                },
                {
                    path:'/service_create', // 创建服务
                    name:'ServiceCreate',
                    component:ServiceCreate
                },
                {
                    path:'/service_details/:id', // 服务详情
                    name:'ServiceDetails',
                    component:ServiceDetails,
                    props:true
                },
                {
                    path:'/user_details', // 我的资料
                    name:'UserDetails',
                    component:UserDetails
                },
                {
                    path:'/user_update', // 更改资料
                    name:'UserUpdate',
                    component:UserUpdate
                },
                {
                    path:'/achieve_statistical', // 成果报表
                    name:'AchieveStatistical',
                    component:AchieveStatistical
                },
                {
                    path:'/order_statistical', // 订单报表
                    name:'OrderStatistical',
                    component:OrderStatistical
                },
            ]
        },
        {
            path:'/adminHome', // 管理员首页
            component:AdminHome
        }
    ]
})

router.beforeEach((to,from,next)=>{ // beforeEach是router的钩子函数，在进入路由前执行
    document.title='代码产权可信互用平台' // 修改页面标题
    next() // 执行进入路由，如果不写就不会进入目标页
})

// 向外共享路由的实例对象
export default router
