<template>
    <div style="background: #f1f1f1; height: 100vh;">
        <!-- 顶部导航栏 -->
        <el-row>
            <el-menu
                    :default-active="activeIndex"
                    class="el-menu-demo"
                    mode="horizontal"
                    @select="handleSelect"
                    background-color="#324057"
                    text-color="#fff"
                    active-text-color="#fff">
                <el-col :span="5">
                    <el-menu-item style="font-size: x-large; color: #ffd04b; font-family: 楷体">代码产权可信互用平台</el-menu-item>
                </el-col>
                <el-col :span="1" offset="17">
                    <el-button type="text" icon="el-icon-remove-outline" id="el_icon_exitlogin" @click="open">注销</el-button>
                </el-col>
            </el-menu>
        </el-row>

        <el-row style="height: calc(100vh - 60px); overflow: hidden;"> <!-- 减去顶部导航栏的高度 -->
              <!-- 侧边栏 -->
            <el-col :span="3" style="height: 100vh; background-color: #FFFFFF; display: flex; flex-direction: column;"> <!-- 铺满整个界面 -->
                <el-menu
                  router="true"
                  style="flex: 1; overflow-y: auto;"
                  class="el-menu-vertical-demo"
                  @open="handleOpen"
                  @close="handleClose"
                  unique-opened="true"
                  background-color="#FFFFFF"
                  text-color="grey"
                  active-text-color="dodgerblue">

          <el-menu-item index="1" route="/userhome">
            <i class="el-icon-s-home"></i>
            <span style="font-size: medium">首页</span>
          </el-menu-item>

          <el-submenu index="2">
            <template slot="title">
              <i class="el-icon-s-claim"></i>
              <span style="font-size: medium">成果管理</span>
            </template>
            <el-menu-item index="2-1">在线创作</el-menu-item>
            <el-menu-item index="2-2" route="/draft_upload">上传成果</el-menu-item>
            <el-menu-item index="2-3" route="/draft_box">草稿箱</el-menu-item>
            <el-menu-item index="2-4" route="/apply_for_cip">申请代码产权</el-menu-item>
            <el-menu-item index="2-5" route="/my_cips">查验代码产权</el-menu-item>
          </el-submenu>

          <el-submenu index="3">
            <template slot="title">
              <i class="el-icon-s-shop"></i>
              <span style="font-size: medium">订阅中心</span>
            </template>
            <el-menu-item index="3-1" route="/cips_opened">浏览开放组件</el-menu-item>
            <el-menu-item index="3-2" route="/cart_list">我的购物车</el-menu-item>
            <el-menu-item index="3-3" route="/my_subs">订阅申请记录</el-menu-item>
            <el-menu-item index="2-6" route="/solve_sub">我处理的订单</el-menu-item>
          </el-submenu>

          <el-submenu index="4">
            <template slot="title">
              <i class="el-icon-menu"></i>
              <span style="font-size: medium">组件服务</span>
            </template>
            <el-menu-item index="4-1" route="/service_create">创建服务</el-menu-item>
            <el-menu-item index="4-2" route="/service_list">服务列表</el-menu-item>
            <el-menu-item index="4-3" route="/service_destroy">销毁服务</el-menu-item>
          </el-submenu>

          <el-submenu index="5">
            <template slot="title">
              <i class="el-icon-s-management"></i>
              <span style="font-size: medium">个人信息</span>
            </template>
            <el-menu-item index="5-1" route="/user_details">我的资料</el-menu-item>
            <el-menu-item index="5-2" route="/user_update">更改资料</el-menu-item>
          </el-submenu>

          <el-submenu index="6">
            <template slot="title">
              <i class="el-icon-s-data"></i>
              <span style="font-size: medium">数据中心</span>
            </template>
            <el-menu-item index="6-1" route="/achieve_report">成果报表</el-menu-item>
            <el-menu-item index="6-2" route="/order_report">订单报表</el-menu-item>
          </el-submenu>
        </el-menu>
      </el-col>

      <el-col :span="21" style="height: calc(100vh - 60px); overflow-y: auto;">
        <router-view></router-view>
      </el-col>
    </el-row>
  </div>
</template>

<script>
export default {
  name: "Home",
  data() {
    return {
      activeIndex: '1',
      // routeIndexMap: { // 路由与菜单index映射表
      //     'UserHome': '1',
      //     'ServiceCreate': '4-1',
      //     'ServiceList': '4-2',
      //     'ServiceDestroy': '4-3',
      //     'CipsOpened': '3-1',
      //     'CartList': '3-2',
      //     'MySubs': '3-3',
      //     'SolveSub': '2-6',
      //     'UserInfoQuery': '5-1',
      //     'UserInfoUpdate': '5-2'
      // }
    };
  },
  // watch: {
  //   $route(newVal) {
  //     console.log('当前路由名称:', newVal.name); // 添加调试输出
  //     this.activeIndex = this.routeIndexMap[newVal.name] || '1'
  //   }
  // },
  // created() {
  //     this.activeIndex = this.routeIndexMap[this.$route.name] || '1'
  // },
  methods: {
    handleSelect(key, keyPath) {
      console.log(key, keyPath);
    },
    handleOpen(key, keyPath) {
      console.log(key, keyPath);
    },
    handleClose(key, keyPath) {
      console.log(key, keyPath);
    },
    open() {
      this.$confirm('退出登录后，下次登录时需要重新填写登录信息，您确定要继续吗?', '安全退出', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        window.localStorage.removeItem("userLogin_jwtToken")
        this.$message({
          type: 'success',
          message: '删除成功!'
        });
        this.$router.push('/login')
      }).catch(() => {});
    },
  }
}
</script>

<!--设为全局-->
<style>
#el_icon_exitlogin{
  background: #324057;
  margin-top: 8px;
  font-size: 17px;
  color: lightgoldenrodyellow;
}

#main_range {
  height: calc(100vh - 76px); /* 精确高度计算 */
  margin: 8px 20px;
  background: #FFFFFF;
  box-sizing: border-box;
  overflow-y: auto;
  display: flex;
  flex-direction: column;

  /* 可选：添加最小高度保证内容展示 */
  min-height: 300px;

  /* 可选：美化滚动条 */
  scrollbar-width: thin;
  scrollbar-color: #888 transparent;
}

/* 滚动条样式美化（可选） */
#main_range::-webkit-scrollbar {
  width: 6px;
}
#main_range::-webkit-scrollbar-thumb {
  background-color: #888;
  border-radius: 3px;
}

#breadcrumb{
  margin-left: 25px;
  margin-top: 10px
}
</style>
