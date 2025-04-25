<template>
  <div>
    <el-row id="breadcrumb">
      <el-breadcrumb separator-class="el-icon-arrow-right">
        <el-breadcrumb-item style="font-size: small">个人信息</el-breadcrumb-item>
        <el-breadcrumb-item>更改资料</el-breadcrumb-item>
      </el-breadcrumb>
    </el-row>
    <el-row id="main_range">
      <el-col span="16" offset="4">
        <br><br>
        <el-card class="box-card">
          <div slot="header" class="clearfix">
            <span>修改我的资料</span>
            <el-button style="float: right; padding: 3px 0" type="text" @click="submitForm('updateData')">提交</el-button>
          </div>
          <div>
            <el-form :model="updateData" :rules="rules" ref="updateData" label-width="100px">
              <br>
              <el-form-item label="新密码" prop="pwd">
                <el-input v-model="updateData.pwd" placeholder="请输入您的新密码" clearable></el-input>
              </el-form-item>
              <el-form-item label="姓名" prop="name">
                <el-input v-model="updateData.name" placeholder="请输入您的姓名" clearable></el-input>
              </el-form-item>
              <el-form-item label="新手机号" prop="tel">
                <el-input v-model="updateData.tel" placeholder="请输入您的新手机号" clearable></el-input>
              </el-form-item>
              <el-form-item label="新地址" prop="addr">
                <el-input v-model="updateData.addr" placeholder="请输入您的新地址" clearable></el-input>
              </el-form-item>
              <el-form-item label="性别" prop="gender">
                <el-select v-model="updateData.gender" style="width: 814px;" placeholder="请选择您的性别">
                  <el-option label="男" value="男"></el-option>
                  <el-option label="女" value="女"></el-option>
                </el-select>
              </el-form-item>
            </el-form>
            <br>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script>
export default {
  name: "UserUpdate",
  created() {
    this.getUserInfo()
  },
  data(){
    var validatePwd = (rule, value, callback) => { // 注册时，密码校验，必须放到return{}之前
      if (value === '') {
        callback(new Error('请输入密码'));
      } else {
        callback();
      }
    };
    return{
      updateData:{
        pwd:'',
        name:'',
        tel:'',
        addr:'',
        gender:''
      },
      rules: {
        pwd: [
          { required:true,message:'请设置新密码',trigger: 'blur' },
          { min: 5, max: 12, message: '长度在 5 到 12 个字符', trigger: 'blur' },
          { validator: validatePwd, trigger: 'blur' }
        ],
        name:[
          { required:true,message:'请设置姓名',trigger: 'blur' },
          { min: 2, max: 20, message: '长度在 2 到 20 个字符', trigger: 'blur' }
        ],
        tel:[
          { required:true,message:'请设置手机号',trigger: 'blur' },
          { type: 'number',
            message: '手机号格式不正确',
            trigger: 'blur',
            transform(value) {
              var phonereg = 11 && /^((13|14|15|16|17|18|19)[0-9]{1}\d{8})$/
              if (!phonereg.test(value)) {
                return false
              } else {
                return Number(value)
              }
            }
          }
        ],
        addr:[
          { required:true,message:'请设置地址',trigger: 'blur' },
        ],
        gender:[
          { required:true,message:'请设置新地址',trigger: 'blur' },
        ]
      }
    }
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
        // 将原信息回显至表单中
        this.updateData.pwd=res.data.pwd
        this.updateData.name=res.data.name
        this.updateData.tel=res.data.tel
        this.updateData.addr=res.data.addr
        this.updateData.gender=res.data.gender
      }
    },
    submitForm(formName) { // 提交修改之后的表单
      this.$refs[formName].validate((valid) => {
        if (valid) {
          this.sendUpdateInfo()
        } else {
          this.$message({
            message: '请先填写完成新资料',
            type: 'warning'
          });
          return false;
        }
      });
    },
    async sendUpdateInfo(){ // 发送修改请求
      const {data:res}=await this.$http({
        url:'/user/update',
        method:'PUT',
        data:{
          pwd:this.updateData.pwd,
          name:this.updateData.name,
          tel:this.updateData.tel,
          addr:this.updateData.addr,
          gender:this.updateData.gender
        },
        headers:{
          "Authorization":'Bearer '+localStorage.getItem("userLogin_jwtToken")
        }
      })
      if(res.code===0){
        this.$message.error(res.msg);
      } else {
        this.$message({
          showClose: true,
          message: '修改成功',
          type: 'success'
        });
        // 修改成功，返回登录
        window.localStorage.removeItem("userLogin_jwtToken")
        await this.$router.push('/login')
      }
    },
    resetForm(formName) {
      this.$refs[formName].resetFields();
    }
  }
}
</script>

<style scoped>

</style>
