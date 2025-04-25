<template>
  <div>
    <el-row id="breadcrumb">
      <el-breadcrumb separator-class="el-icon-arrow-right">
        <el-breadcrumb-item style="font-size: small">订阅中心</el-breadcrumb-item>
        <el-breadcrumb-item :to="{ name: 'SolveSub' }">我处理的订单</el-breadcrumb-item>
        <el-breadcrumb-item>处理订单</el-breadcrumb-item>
      </el-breadcrumb>
    </el-row>
    <el-row id="main_range">
      <br>
      <el-row style="margin-bottom: 30px;">
        <el-col :span="10" :offset="1">
          <el-row>
            <el-card shadow="always" class="box-card">
              <div class="scroll-container">
                <el-descriptions :column="4" size="medium" direction="vertical" title="组件代码知识产权" border>
                <el-descriptions-item label="名称" :span="2">{{this.cipInfo.filename}}</el-descriptions-item>
                <el-descriptions-item label="版本号" :span="2">{{this.cipInfo.version}}</el-descriptions-item>
                <el-descriptions-item label="代码产权数字摘要" :span="4"><el-tag effect="plain">{{this.cipInfo.cpi}}</el-tag></el-descriptions-item>
                <el-descriptions-item label="确权日期时间">{{this.cipInfo.regTime}}</el-descriptions-item>
              </el-descriptions>
              </div>
            </el-card>
          </el-row>

          <br><br>

          <el-row>
            <el-card shadow="always" class="box-card">
              <div class="scroll-container">
                <el-descriptions :column="4" size="medium" title="组件信息" direction="vertical" border>
                  <el-descriptions-item label="组件标签" :span="1">{{this.cipInfo.tag}}</el-descriptions-item>
                  <el-descriptions-item label="程序设计语言" :span="1">{{this.cipInfo.language}}</el-descriptions-item>
                  <el-descriptions-item label="初次上传时间" :span="2">{{this.cipInfo.createdTime}}</el-descriptions-item>
                  <el-descriptions-item label="近期更新时间" :span="2">{{this.cipInfo.updatedTime}}</el-descriptions-item>
                  <el-descriptions-item label="组件描述" :span="2">{{this.cipInfo.description}}</el-descriptions-item>
                  <el-descriptions-item label="源代码文件存储地址" :span="4"><el-tag effect="light">{{this.cipInfo.addr}}</el-tag></el-descriptions-item>
                </el-descriptions>
              </div>
            </el-card>
          </el-row>
        </el-col>

        <el-col :span="11" :offset="1">
          <el-row>
            <el-card shadow="hover" class="box-card">
              <div class="scroll-container">
                <el-descriptions :column="3" size="medium" title="订阅申请信息" direction="vertical" border>
                  <template slot="extra">
                    <el-button type="danger" size="small" @click="handleRefuse">拒绝交易</el-button>
                  </template>
                  <el-descriptions-item label="订阅方姓名">{{this.partyAInfo.name}}</el-descriptions-item>
                  <el-descriptions-item label="手机号"><el-tag effect="plain">{{this.partyAInfo.tel}}</el-tag></el-descriptions-item>
                  <el-descriptions-item label="所属组织">{{this.partyAInfo.role}}</el-descriptions-item>
                  <el-descriptions-item label="订阅生效时间">{{this.subDetails.startTime}}</el-descriptions-item>
                  <el-descriptions-item label="乙方初拟价格">{{this.cipInfo.price}}</el-descriptions-item>
                  <el-descriptions-item label="甲方意向价格"><el-tag effect="light" type="danger">{{this.subDetails.realPrice}}</el-tag></el-descriptions-item>
                  <el-descriptions-item label="下单时间">{{this.subDetails.createdTime}}</el-descriptions-item>
                  <el-descriptions-item label="订单当前状态">{{this.subDetails.status}}</el-descriptions-item>
                </el-descriptions>
              </div>
            </el-card>
          </el-row>

          <br>

          <el-row>
            <el-card shadow="hover" class="box-card">
                <div slot="header" class="clearfix">
                  <span>修改订阅信息</span>
                  <el-button style="float: right; padding: 3px 0" type="text" @click="handleAgree('newSub')">同意交易</el-button>
                </div>
                <div class="scroll-container">
                <div>
                <el-form :model="newSub" :rules="newSubRules" ref="newSub" label-width="150px">
                  <el-form-item label="订阅价格" prop="price">
                    <el-input v-model="newSub.price" placeholder="请输入成交价格">
                        <template slot="append">（元/天）</template>
                    </el-input>
                  </el-form-item>
                  <el-form-item label="订阅生效时间" prop="startTime">
                    <el-date-picker
                        v-model="newSub.startTime"
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
                        v-model="newSub.note">
                    </el-input>
                  </el-form-item>
                </el-form>
                </div>
              </div>
            </el-card>
          </el-row>
        </el-col>
      </el-row>
    </el-row>
  </div>
</template>

<script>
export default {
  name: "Solve",
  props: ['id'],
  data(){
    return{
      subDetails:{},
      cipInfo:{},
      partyAInfo:{},
      newSub:{
        price:'',
        startTime:'',
        note:''
      },
      newSubRules:{
        price: [
          { required: true, message: '请输入您接受的价格', trigger: 'blur' },
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
          { required: true, message: '请选择您接受的订阅生效时间', trigger: 'blur' },
          { validator: this.validateStartTime, trigger: 'blur' }
        ],
        note:[
          { max: 250, message: '输入内容不超过250字', trigger: 'blur' }
        ]
      }
    }
  },
  created() {
    this.getInfo();
  },
  methods:{
    getInfo(){
      this.querySubInfo().then(()=>{
        this.queryCipInfo().then(()=>{
          this.queryPartyAInfo(); // 查看使用者用户数据
        });
      });
    },
    async querySubInfo(){
      const {data:res}=await this.$http({
        url:'/subscribe/query',
        method:'POST',
        headers:{
          "Authorization":'Bearer '+localStorage.getItem("userLogin_jwtToken")
        },
        data:{
          id:this.id
        }
      })
      if(res.code===1){
        this.subDetails=res.data;
        this.newSub.price=this.subDetails.realPrice;
        this.newSub.startTime=this.subDetails.startTime;
        this.newSub.note=this.subDetails.note;
      } else {
        this.$message.error(res.msg);
      }
    },
    async queryCipInfo(){
      const {data:res}=await this.$http({
        url:'/code/queryDetails',
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
        if (this.cipInfo.description.length > 150) {
          this.cipInfo.description = this.cipInfo.description.substring(0, 150) + '······（略）';
        }
      } else {
        this.$message.error(res.msg);
      }
    },
    async queryPartyAInfo(){
      const {data:res}=await this.$http({
        url:'/user/queryTarget',
        method:'POST',
        headers:{
          "Authorization":'Bearer '+localStorage.getItem("userLogin_jwtToken")
        },
        data:{
          username:this.subDetails.partyA
        }
      })
      if(res.code===1){
        this.partyAInfo=res.data;
      } else {
        this.$message.error(res.msg);
      }
    },
    handleRefuse(){
      this.$confirm('您确定拒绝此申请吗', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        this.refuse();
      });
    },
    async refuse(){
      const {data:res}=await this.$http({
        url:'/subscribe/review',
        method:'POST',
        headers:{
          "Authorization":'Bearer '+localStorage.getItem("userLogin_jwtToken")
        },
        data:{
          id:this.subDetails.id,
          status:'乙方已拒绝'
        }
      })
      if(res.code===1){
        this.$message({
          showClose: true,
          message: '您拒绝了此申请',
          type: 'success'
        });
        var self = this;
        setTimeout(function() {
          self.$router.push({ name: 'SolveSub'});
        }, 1000);
      } else {
        this.$message.error(res.msg);
      }
    },
    handleAgree(formName){
      this.$refs[formName].validate((valid) => {
        if (valid) {
          this.$confirm('您确定同意此申请吗', '提示', {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning'
          }).then(() => {
            if (this.newSub.note===''){
              this.newSub.note='无'
            }
            this.agree();
          });
        } else {
          return false;
        }
      });
    },
    async agree(){
      const {data:res}=await this.$http({
        url:'/subscribe/review',
        method:'POST',
        headers:{
          "Authorization":'Bearer '+localStorage.getItem("userLogin_jwtToken")
        },
        data:{
          id:this.subDetails.id,
          realPrice:parseInt(this.newSub.price),
          startTime:this.newSub.startTime,
          note:this.newSub.note,
          status:'乙方同意，待甲方确认订阅条款'
        }
      })
      if(res.code===0){
        this.$message.error(res.msg);
      } else {
        var self = this;
        setTimeout(function() {
          self.$router.push({ name: 'SolveSub'});
        }, 1000);
        this.$notify({
          title: '成功',
          message: '处理成功',
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
.scroll-container {
    max-height: 400px;
    overflow-y: auto;
}

/* 优化滚动条样式 */
.scroll-container::-webkit-scrollbar {
    width: 6px;
}
.scroll-container::-webkit-scrollbar-thumb {
    background: #c1c1c1;
    border-radius: 4px;
}
</style>
