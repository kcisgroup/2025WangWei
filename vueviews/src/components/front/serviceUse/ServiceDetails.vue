<template>
  <div>
    <el-row id="breadcrumb">
      <el-breadcrumb separator-class="el-icon-arrow-right">
        <el-breadcrumb-item style="font-size: small">组件服务</el-breadcrumb-item>
        <el-breadcrumb-item style="font-size: small" :to="{ name: 'ServiceList' }">服务列表</el-breadcrumb-item>
        <el-breadcrumb-item>服务详情</el-breadcrumb-item>
      </el-breadcrumb>
    </el-row>
    <el-row id="main_range">
      <br>

      <el-row>
        <el-col :span="20" :offset="2">
            <el-card shadow="hover" class="box-card">
                <div slot="header" class="clearfix">
                    <span>服务详细信息</span>
                </div>
                <el-descriptions :column="2" border>
                    <el-descriptions-item label="项目名">{{this.serviceInfo.filename}}</el-descriptions-item>
                    <el-descriptions-item label="项目版本号">{{this.serviceInfo.version}}</el-descriptions-item>
                    <el-descriptions-item label="乙方联系人">{{this.serviceInfo.author}}</el-descriptions-item>
                    <el-descriptions-item label="乙方单位">{{this.serviceInfo.company}}</el-descriptions-item>
                    <el-descriptions-item label="容器名" :span="2">{{this.serviceInfo.dockerName}}</el-descriptions-item>
                    <el-descriptions-item label="容器创建时间" :span="2">{{this.serviceInfo.createdTime}}</el-descriptions-item>
                    <el-descriptions-item label="IP地址">{{this.serviceInfo.ip}}</el-descriptions-item>
                    <el-descriptions-item label="端口">{{this.serviceInfo.port}}</el-descriptions-item>
                    <el-descriptions-item label="状态">{{this.serviceInfo.status}}</el-descriptions-item>
                    <el-descriptions-item label="使用时长（小时）">{{this.serviceInfo.term}}</el-descriptions-item>
                    <el-descriptions-item label="价格（元/天）">{{this.serviceInfo.realPrice}}</el-descriptions-item>
                    <el-descriptions-item label="当前费用（元）">{{this.serviceInfo.cost}}</el-descriptions-item>
                    <el-descriptions-item label="功能描述" :span="2">{{this.serviceInfo.description}}</el-descriptions-item>
                </el-descriptions>
            </el-card>

            <br>

            <el-row>
                <el-col :span="24" :offset="18">
                    <span style="font-size: small;">数据截止至  {{this.serviceInfo.currentTime}}</span>
                </el-col>
            </el-row>
          <br>
        </el-col>
      </el-row>
    </el-row>
  </div>
</template>

<script>
export default {
  name: "ServiceDetails",
  props: ['id'],
  data(){
    return{
      serviceInfo:{}
    }
  },
  created() {
    this.getDetails();
  },
  methods: {
    async getDetails() {
      const {data: res} = await this.$http({
        url: '/product/query',
        method: 'POST',
        headers: {
          "Authorization": 'Bearer ' + localStorage.getItem("userLogin_jwtToken")
        },
        data: {
          id:this.id
        }
      })
      if (res.code === 0) {
        this.$message.error(res.msg);
      } else {
        this.serviceInfo = res.data;
      }
    },
    backToMySubs(){
      this.$router.push({ name: 'ServiceList'});
    }
  }
}
</script>

<style scoped>
</style>