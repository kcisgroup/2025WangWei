<template>
  <div>
    <el-row id="breadcrumb">
      <el-breadcrumb separator-class="el-icon-arrow-right">
        <el-breadcrumb-item style="font-size: small">订阅中心</el-breadcrumb-item>
        <el-breadcrumb-item :to="{ name: 'MySubs' }">订阅申请记录</el-breadcrumb-item>
        <el-breadcrumb-item>确认合同条款</el-breadcrumb-item>
      </el-breadcrumb>
    </el-row>
    <el-row id="main_range">
      <br>
      <el-row>
        <el-col :span="22" :offset="1">
          <el-card shadow="always" class="box-card">
            <el-descriptions :column="5" size="medium" title="组件信息" direction="vertical" border>
              <el-descriptions-item label="组件名" :span="1">{{this.cipInfo.filename}}</el-descriptions-item>
              <el-descriptions-item label="版本号" :span="1">{{this.cipInfo.version}}</el-descriptions-item>
              <el-descriptions-item label="组件标签" :span="1"><el-tag effect="light" type="warning">{{this.cipInfo.tag}}</el-tag></el-descriptions-item>
              <el-descriptions-item label="程序设计语言" :span="1">{{this.cipInfo.language}}</el-descriptions-item>
              <el-descriptions-item label="创作者" :span="1">{{this.cipInfo.name}}</el-descriptions-item>
              <el-descriptions-item label="乙方手机号" :span="1">{{this.cipInfo.tel}}</el-descriptions-item>
              <el-descriptions-item label="组件获权时间" :span="2">{{this.cipInfo.regTime}}</el-descriptions-item>
              <el-descriptions-item label="组件更新时间" :span="2">{{this.cipInfo.updatedTime}}</el-descriptions-item>
              <el-descriptions-item label="组件描述" :span="4">{{this.cipInfo.description}}</el-descriptions-item>
            </el-descriptions>
          </el-card>
        </el-col>
      </el-row>

      <br>

      <el-row>
        <el-col :span="22" :offset="1">
          <el-card shadow="always" class="box-card">
            <el-descriptions :column="3" size="medium" title="订阅申请信息" direction="vertical" border>
              <template slot="extra">
                <el-button type="primary" size="small" @click="handleConfirm" plain>确认</el-button>
                <el-button type="danger" size="small" @click="handleRefuse" plain>拒绝</el-button>
              </template>
              <el-descriptions-item label="订阅价格(元/天)"><el-tag effect="light" type="danger">{{this.subDetails.realPrice}}</el-tag></el-descriptions-item>
              <el-descriptions-item label="合同生效时间">{{this.subDetails.startTime}}</el-descriptions-item>
              <el-descriptions-item label="订单当前状态" :span="2">{{this.subDetails.status}}</el-descriptions-item>
              <el-descriptions-item label="备注">{{this.subDetails.note}}</el-descriptions-item>
            </el-descriptions>
          </el-card>
        </el-col>
      </el-row>
    </el-row>
  </div>
</template>

<script>
export default {
  name: "ConfirmSub",
  props: ['id'],
  data(){
    return{
      subDetails:{},
      cipInfo:{},
    }
  },
  mounted() {
    this.getInfo();
  },
  methods:{
    getInfo(){
      this.querySubInfo().then(()=>{
        this.queryCipInfo();
      });
    },
    async querySubInfo(){
      const {data:res}=await this.$http({
        url:'/subscribe/query',
        method:'POST',
        headers: {
          "Authorization": 'Bearer ' + localStorage.getItem("userLogin_jwtToken")
        },
        data:{
          id:this.id
        }
      })
      if(res.code===1){
        this.subDetails=res.data;
      } else {
        this.$message.error(res.msg);
      }
    },
    async queryCipInfo(){
      console.log(this.subDetails);
      const {data:res}=await this.$http({
        url:'/code/queryCipInfo',
        method:'POST',
        headers:{
          "Authorization":'Bearer '+localStorage.getItem("userLogin_jwtToken")
        },
        data:{
          id:this.subDetails.codeId
        }
      })
      if(res.code===1){
        this.cipInfo=res.data;
      } else {
        this.$message.error(res.msg);
      }
    },
    handleRefuse(){
      this.$confirm('您确定拒绝订阅条款', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        this.refuse();
      });
    },
    async refuse(){
      const {data:res}=await this.$http({
        url:'/subscribe/confirm',
        method:'POST',
        headers:{
          "Authorization":'Bearer '+localStorage.getItem("userLogin_jwtToken")
        },
        data:{
          id:this.subDetails.id,
          status:'甲方拒绝订阅条款'
        }
      })
      if(res.code===1){
        this.$message({
          showClose: true,
          message: '您拒绝了本次订阅申请',
          type: 'success'
        });
        var self = this;
        setTimeout(function() {
          self.$router.push({ name: 'MySubs'});
        }, 1000);
        this.$notify({
          title: '操作成功',
          message: '您拒绝了订阅条款，交易关闭',
          type: 'warning'
        });
      } else {
        this.$message.error(res.msg);
      }
    },
    handleConfirm(formName){
      this.$confirm('您确定接受订阅条款并签订合同吗?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        this.confirm();
      });
    },
    async confirm(){
        const loading = this.$loading({
            lock: true,
            text: '操作已提交，请耐心等待',
            spinner: 'el-icon-loading',
            background: 'rgba(0, 0, 0, 0.7)'
        });

        const {data:res}=await this.$http({
        url:'/subscribe/confirm',
        method:'POST',
        headers:{
          "Authorization":'Bearer '+localStorage.getItem("userLogin_jwtToken")
        },
        data:{
          id:this.subDetails.id,
          status:'双方已签订合同'
        }
      })
      if(res.code===0){
        this.$message.error(res.msg);
      } else {
        var self = this;
        setTimeout(function() {
          self.$router.push({ name: 'MySubs'});
        }, 1000);
        this.$notify({
          title: '操作成功',
          message: '恭喜您，成功订阅了此组件',
          type: 'success'
        });
      }

      loading.close();
    }
  }
}
</script>

<style scoped>

</style>
