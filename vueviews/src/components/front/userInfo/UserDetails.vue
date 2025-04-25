<template>
  <div>
    <el-row id="breadcrumb">
      <el-breadcrumb separator-class="el-icon-arrow-right">
        <el-breadcrumb-item style="font-size: small">个人信息</el-breadcrumb-item>
        <el-breadcrumb-item>我的资料</el-breadcrumb-item>
      </el-breadcrumb>
    </el-row>
    <el-row id="main_range">
      <br>

      <el-col :span="16" :offset="4">
            <el-card shadow="hover" class="box-card">
                <div slot="header" class="clearfix">
                    <span>账号信息</span>
                </div>
                <el-descriptions :column="1" border>
                    <el-descriptions-item label="用户名">{{this.userInfoData.username}}</el-descriptions-item>
                    <el-descriptions-item label="密码">{{this.userInfoData.pwd}}</el-descriptions-item>
                    <el-descriptions-item label="用户类型">{{this.userInfoData.status === 0 ? '管理员' : '普通员工'}}</el-descriptions-item>
                    <el-descriptions-item label="身份证号">{{this.userInfoData.idnum}}</el-descriptions-item>
                    <el-descriptions-item label="单位">{{this.userInfoData.role}}</el-descriptions-item>
                    <el-descriptions-item label="姓名">{{this.userInfoData.name}}</el-descriptions-item>
                    <el-descriptions-item label="手机号">{{this.userInfoData.tel}}</el-descriptions-item>
                    <el-descriptions-item label="性别">{{this.userInfoData.gender}}</el-descriptions-item>
                    <el-descriptions-item label="地址">{{this.userInfoData.addr}}</el-descriptions-item>
                </el-descriptions>
            </el-card>
        </el-col>
    </el-row>
  </div>
</template>

<script>
export default {
  name: "UserDetails",
  data(){
    return{
      userInfoData: {}
    }
  },
  created() {
    this.getUserInfo()
  },
  methods:{
    async getUserInfo(){ // 发送请求
      const {data:res}=await this.$http({
        url:'/user/query',
        method:'GET',
        headers:{
          "Authorization":'Bearer '+localStorage.getItem("userLogin_jwtToken")
        }
      })
      if(res.code===0){
        this.$message.error('您的个人信息在赶来的路上走丢了，请尝试重新操作或重新登录');
      } else {
          this.userInfoData=res.data;
      }
    }
  }
}
</script>

<style scoped>

</style>
