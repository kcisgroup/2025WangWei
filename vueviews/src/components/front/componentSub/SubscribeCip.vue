<template>
  <div>
    <el-row id="breadcrumb">
      <el-breadcrumb separator-class="el-icon-arrow-right">
        <el-breadcrumb-item style="font-size: small">订阅中心</el-breadcrumb-item>
        <el-breadcrumb-item :to="{ name: 'CartList' }">我的购物车</el-breadcrumb-item>
        <el-breadcrumb-item>订阅申请表</el-breadcrumb-item>
      </el-breadcrumb>
    </el-row>
    <el-row id="main_range">
      <div style="margin-top: 10px;margin-bottom: 10px;">
        <el-row class="row-bg">
          <el-col :span="20" :offset="2">
            <el-card class="box-card" shadow="hover">
              <div slot="header" class="clearfix">
                <span>订阅申请表</span>
                <el-button style="float: right; padding: 3px 0" type="text" @click="applySub('subTable')">提交申请</el-button>
              </div>

              <div>
                <el-form :model="subTable" :rules="subRules" ref="subTable" label-width="150px" size="medium" class="demo-ruleForm">
                  <el-form-item label="订阅项目名">
                    <el-input v-model="projectInfo.filename" disabled></el-input>
                  </el-form-item>
                  <el-form-item label="项目版本号">
                    <el-input v-model="projectInfo.version" disabled></el-input>
                  </el-form-item>
                  <el-form-item label="项目标签">
                    <el-input v-model="projectInfo.tag" disabled></el-input>
                  </el-form-item>
                  <el-form-item label="甲方联系人">
                    <el-input v-model="userInfoData.name" disabled></el-input>
                  </el-form-item>
                  <el-form-item label="甲方联系电话">
                    <el-input v-model="userInfoData.tel" disabled></el-input>
                  </el-form-item>
                  <el-form-item label="甲方单位">
                    <el-input v-model="userInfoData.role" disabled></el-input>
                  </el-form-item>
                  <el-form-item label="乙方联系人">
                    <el-input v-model="projectInfo.name" disabled></el-input>
                  </el-form-item>
                  <el-form-item label="乙方联系电话">
                    <el-input v-model="projectInfo.tel" disabled></el-input>
                  </el-form-item>
                  <el-form-item label="乙方单位">
                    <el-input v-model="projectInfo.company" disabled></el-input>
                  </el-form-item>
                  <el-form-item label="甲方初拟价格">
                    <el-input v-model="projectInfo.price" disabled>
                        <template slot="append">（元/天）</template>
                    </el-input>
                  </el-form-item>
                  <el-form-item label="意向价格" prop="price">
                      <el-input v-model="subTable.price" placeholder="请填写意向价格">
                          <template slot="append">（元/天）</template>
                      </el-input>
                  </el-form-item>
                  <el-form-item label="订阅生效时间" prop="startTime">
                    <el-date-picker
                        v-model="subTable.startTime"
                        type="datetime"
                        value-format="yyyy-MM-dd HH:mm:ss"
                        placeholder="请选择订阅生效时间">
                    </el-date-picker>
                  </el-form-item>
                  <el-form-item label="备注" prop="note">
                    <el-input
                        type="textarea"
                        :rows="3"
                        placeholder="请输入备注"
                        v-model="subTable.note">
                    </el-input>
                  </el-form-item>
                </el-form>
              </div>
            </el-card>
          </el-col>
        </el-row>
      </div>
    </el-row>
  </div>
</template>

<script>
export default {
  name: "SubscribeCip",
  props: ['id'],
  data(){
    return {
      cartInfo:{},
      userInfoData:{},
      projectInfo: {},
      subTable: {
        price:'',
        startTime:'',
        note:''
      },
      subRules:{
        price: [
          { required: true, message: '请输入商议后的价格', trigger: 'blur' },
          {
            validator: (rule, value, callback) => {
              if (!/^\d+$/.test(value)) {
                callback(new Error('输入的价格只能包含数字，不包含小数和其他符合'));
              } else {
                callback();
              }
            },
            trigger: 'blur'
          }
        ],
        startTime:[
          { required: true, message: '请选择订阅生效时间', trigger: 'blur' },
          { validator: this.validateStartTime, trigger: 'blur' }
        ],
        note:[
          { max: 250, message: '输入内容不超过250字', trigger: 'blur' }
        ]
      },
    }
  },
  created() {
    this.getUserInfo();
    this.queryCartAndCipDetails();
  },
  methods: {
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
    },
    async queryCartAndCipDetails(){
      const {data:res}=await this.$http({
        url:'/subscribe/queryCartDetails',
        method:'POST',
        headers:{
          "Authorization":'Bearer '+localStorage.getItem("userLogin_jwtToken")
        },
        data:{
          id:parseInt(this.id)
        }
      })
      if(res.code===0){
        this.$message.error(res.msg);
      } else {
        this.cartInfo=res.data;
        await this.getCipDetails();
      }
    },
    async getCipDetails() {
      const {data: res} = await this.$http({
        url: '/code/queryCipInfo',
        method: 'POST',
        headers: {
          "Authorization": 'Bearer ' + localStorage.getItem("userLogin_jwtToken")
        },
        data: {
          id: parseInt(this.cartInfo.codeId)
        }
      })
      if (res.code === 0) {
        this.$message.error(res.msg);
      } else {
        this.projectInfo = res.data;
        this.assignCipValues();
      }
    },
    applySub(formName){
      this.$refs[formName].validate((valid) => {
        if (valid) {
          this.$confirm('您确定要提交此订阅申请吗', '提示', {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning'
          }).then(() => {
            const loading = this.$loading({
              lock: true,
              text: '正在提交，请耐心等待，请勿关闭',
              spinner: 'el-icon-loading',
              background: 'rgba(0, 0, 0, 0.7)'
            });
            if (this.subTable.note===''){
              this.subTable.note='无'
            }
            this.subscribe().then(()=>{
              loading.close();
            });
          });
        } else {
          return false;
        }
      });
    },
    async subscribe(){
      const {data:res}=await this.$http({
        url:'/subscribe/new',
        method:'POST',
        headers:{
          "Authorization":'Bearer '+localStorage.getItem("userLogin_jwtToken")
        },
        data:{
          id:parseInt(this.id),
          codeId:parseInt(this.projectInfo.id),
          realPrice:parseInt(this.subTable.price),
          startTime:this.subTable.startTime,
          note:this.subTable.note
        }
      })
      if(res.code===0){
        this.$message.error(res.msg);
      } else {
        var self = this;
        setTimeout(function() {
          self.$router.push({ name: 'CartList'});
        }, 1000);
        this.$notify({
          title: '成功',
          message: '您的订阅申请已提交',
          type: 'success'
        });
      }
    },
    validateStartTime(rule, value, callback) {
      if (value && new Date(value) <= new Date()) {
        callback(new Error('注意：您选择的日期时间早于当前中国标准时间，为避免不必要的费用支出，请重新选择日期时间'));
      } else {
        callback();
      }
    }
  }
}
</script>

<style scoped>

</style>
