<template>
  <div id="building">
    <el-row id="elrowform"><!--layout布局-->
      <el-col :span="10" id="elrowcol" :offset="7">
        <div class="grid-content"><!--偏移-->
          <el-form :model="ruleForm" status-icon hide-required-asterisk="true" label-width="auto" :rules="rules" ref="ruleForm">

              <el-row>
                <el-col :span="24">
                    <h2 style="color: forestgreen;font-family: 仿宋;margin-bottom: 35px; text-align: center">
                        代码产权可信互用平台
                    </h2>
                </el-col>
            </el-row>

            <el-form-item label="用户名" prop="username"> <!--配合el-form标签中:rules="rules"使用，用于传递rules中参数，从而实现输入合法性校验-->
              <el-input type="text" v-model="ruleForm.username" autocomplete="off" clearable></el-input>
            </el-form-item>
            <el-form-item label="密码" prop="pass" style="margin-top: 35px;">
              <el-input
                type="password"
                v-model="ruleForm.pass"
                autocomplete="off"
                clearable
                @keyup.enter.native="submitForm('ruleForm')">
              </el-input>
            </el-form-item>
              
            <el-form-item>
              <el-row type="flex" justify="center">
                <el-col :span="6">
                  <el-button type="primary" @click="submitForm('ruleForm')">登&nbsp&nbsp录</el-button>
                </el-col>
                <el-col :span="6" :offset="1">
                  <el-button @click="resetForm('ruleForm')">重&nbsp&nbsp置</el-button>
                </el-col>
              </el-row>
            </el-form-item>

              <!--注册功能：嵌套form的抽屉-->
            <el-row type="flex" justify="end">
              <el-col :span="6">
                <el-button
                  type="text"
                  @click="dialog=true"
                  size="small"
                  style="font-size: medium;color: darkgrey">
                  点击注册
                </el-button>
              </el-col>
            </el-row>

            <el-drawer
                title="账号注册"
                :before-close="handleClose"
                :visible.sync="dialog"
                direction="rtl"
                custom-class="demo-drawer"
                ref="drawer"
                size="35%"
            >
              <div>
                <el-form :model="form" :rules="rules" ref="form" status-icon><!--:model双向绑定名为form的表单，v-bind:rules绑定规则:校验输入内容是否符合规范-->
                  <el-form-item label="用户名" prop="username" :label-width="formLabelWidth" id="drawerForm"><!--prop配合上行:rules="rules"使用，用于传递rules中规则，从而实现校验-->
                    <el-input v-model="form.username" placeholder="请使用您的工号" autocomplete="off" clearable></el-input>
                  </el-form-item>
                  <el-form-item label="密码" prop="pwd" :label-width="formLabelWidth" id="drawerForm">
                    <el-input type="password" v-model="form.pwd" autocomplete="off" placeholder="请输入密码" show-password></el-input>
                  </el-form-item>
                  <el-form-item label="确认密码" prop="checkPwd" :label-width="formLabelWidth" id="drawerForm">
                    <el-input type="password" v-model="form.checkPwd" autocomplete="off" placeholder="请再次输入密码" show-password></el-input>
                  </el-form-item>
                  <el-form-item label="姓名" prop="name" :label-width="formLabelWidth" id="drawerForm">
                    <el-input v-model="form.name" placeholder="请输入您的姓名" autocomplete="off" clearable></el-input>
                  </el-form-item>
                  <el-form-item label="公司" prop="role" :label-width="formLabelWidth" id="drawerForm">
                    <el-select v-model="form.role" clearable style="width: 100%;" placeholder="请选择您就职的公司">
                      <el-option
                          v-for="item in form.options"
                          :key="item.role"
                          :label="item.label2"
                          :value="item.role">
                      </el-option><!--v-for="(value,key,index) in obj"访问对象-->
                    </el-select>
                  </el-form-item>
                  <el-form-item label="性别" prop="gender" :label-width="formLabelWidth" id="drawerForm">
                    <el-radio v-model="form.gender" label="男">男</el-radio>
                    <el-radio v-model="form.gender" label="女">女</el-radio>
                  </el-form-item>
                  <el-form-item label="手机号" prop="tel" :label-width="formLabelWidth" id="drawerForm">
                    <el-input v-model="form.tel" placeholder="请输入手机号" autocomplete="off" clearable></el-input>
                  </el-form-item>
                  <el-form-item label="身份证号" prop="idnum" :label-width="formLabelWidth" id="drawerForm">
                    <el-input v-model="form.idnum" placeholder="请输入身份证号" autocomplete="off" clearable></el-input>
                  </el-form-item>
                  <el-form-item label="地址" prop="addr" :label-width="formLabelWidth" id="drawerForm">
                    <el-input v-model="form.addr" placeholder="请输入您的所在地" autocomplete="off" clearable></el-input>
                  </el-form-item>
                </el-form>
<!--                <div style="margin-top:80px;margin-left: 95px;">-->
<!--                  <el-button type="primary" @click="submitDrawer('form')" style="width: 220px;height: 45px;font-size: medium;">确&nbsp&nbsp定</el-button>-->
<!--                  <el-button @click="cancelForm" style="margin-left: 35px;width: 220px;height: 45px;font-size: medium;">取&nbsp&nbsp消</el-button>-->
<!--                </div>-->
                  <el-row type="flex" justify="center">
                    <el-col :span="8">
                      <el-button
                        type="primary"
                        @click="submitDrawer('form')"
                        style="height: 45px;font-size: medium;width: 100%">
                        确&nbsp&nbsp定
                      </el-button>
                    </el-col>
                    <el-col :span="8" :offset="1">
                      <el-button
                        @click="cancelForm"
                        style="height: 45px;font-size: medium;width: 100%">
                        取&nbsp&nbsp消
                      </el-button>
                    </el-col>
                  </el-row>
              </div>
            </el-drawer>
          </el-form>
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<script>
export default {
  // eslint-disable-next-line vue/multi-word-component-names
  name: "Login",
  data() {
    var validatePwd = (rule, value, callback) => { //注册时，密码校验，必须放到return{};之前
      if (value === '') {
        callback(new Error('请输入密码'));
      } else {
        if (this.form.checkPwd !== '') {
          this.$refs.form.validateField('checkPwd');
        }
        callback();
      }
    };
    var validatePwd2 = (rule, value, callback) => { //注册时，确认密码校验
      if (value === '') {
        callback(new Error('请再次输入密码'));
      } else if (value !== this.form.pwd) {
        callback(new Error('两次输入密码不一致!'));
      } else {
        callback();
      }
    };
    const isCnNewID = (rule, value, callback) => { //注册时，身份证号校验
      var arrExp = [7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2]; //加权因子
      var arrValid = [1, 0, "X", 9, 8, 7, 6, 5, 4, 3, 2]; //校验码
      if (/^\d{17}\d|x$/i.test(value)) {
        var sum = 0, idx;
        for (var i = 0; i < value.length - 1; i++) {
          // 对前17位数字与权值乘积求和
          sum += parseInt(value.substr(i, 1), 10) * arrExp[i];
        }
        // 计算模（固定算法）
        idx = sum % 11;
        // 检验第18为是否与校验码相等
        if (arrValid[idx] == value.substr(17, 1).toUpperCase()) {
          callback()
        } else {
          callback("身份证号格式有误")
        }
      } else {
        callback("身份证号格式有误")
      }
    };
    return {
      ruleForm: { // 登录表单
        username: '',
        pass: ''
      },
      rules: { // 定义校验规则
        username: [
          { required: true, message: '用户名不能为空', trigger: 'blur' },
          { min: 4, max: 12, message: '长度在 4 到 12 个字符', trigger: 'blur' } //trigger: 'blur'失去焦点时校验
        ],
        pass: [
          { required:true,message:'密码不能为空',trigger: 'blur' },
          { min: 5, max: 12, message: '长度在 5 到 12 个字符', trigger: 'blur' },
          { validator: validatePwd, trigger: 'blur' }
        ],
        name:[
          { required:true,message:'姓名不能为空',trigger: 'blur' },
          { min: 2, max: 20, message: '长度在 2 到 20 个字符', trigger: 'blur' }
        ],
        pwd: [
          { required:true,message:'请输入密码',trigger: 'blur' },
          { min: 5, max: 12, message: '长度在 5 到 12 个字符', trigger: 'blur' },
          { validator: validatePwd, trigger: 'blur' }
        ],
        checkPwd: [
          { required:true,message:'请再次输入密码',trigger: 'blur' },
          { min: 5, max: 12, message: '长度在 5 到 12 个字符', trigger: 'blur' },
          { validator: validatePwd2, trigger: 'blur' }
        ],
        tel:[
          { required:true,message:'手机号不能为空',trigger:'blur' },
            // 手机号格式校验
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
        role:[
          { required:true,message:'请选择就职公司',trigger:'blur' }
        ],
        gender:[
          { required:true,message:'请选择您的性别',trigger:'blur' }
        ],
        idnum: [
          { required:true,message:'身份证号不能为空',trigger:'blur' },
          { validator: isCnNewID,trigger: 'blur' }
        ],
        addr: [
          { required:true,message:'地址不能为空',trigger:'blur' }
        ]
      },
      // 抽屉中的表单项
      formLabelWidth: '80px',
      timer: null,
      dialog: false,
      form: { // 抽屉中注册表单
        username: '',
        pwd: '',
        checkPwd:'',
        name:'',
        role:'', // 所在公司
        options:[{ // 公司选项
          role:'companya',
          label2:'companya'
        },{
          role:'companyb',
          label2:'companyb'
        }],
        gender: '男',
        tel:'',
        idnum:'',
        addr:'',
      },
    };
  },
  methods: {
    submitForm(formName) { // 点击登录按钮后验证和提交
      this.$refs[formName].validate((valid) => {
        if (valid) {
          this.login()
        } else {
          this.$message({
            showClose: true,
            message: '请正确填写用户名和密码',
            type: 'error'
          });
          return false;
        }
      });
    },
    async login(){ // 发送登录请求，接收后端返回参数，判断是员工还是管理员，并跳转不同首页
      const {data:res}=await this.$http.post('/user/login',{
        username:this.ruleForm.username,
        pwd:this.ruleForm.pass
      })
      if (res.code===0){
        this.$message({
          message: res.msg,
          type: 'error'
        });
      } else {
        // 将后端传过来的Jwt Token存入Local Storage中
        window.localStorage.setItem("userLogin_jwtToken",res.token)
        if (res.data===null){ // 员工登录
          await this.$router.push('/home')
          this.$notify({
            title: '登录成功',
            message: '欢迎进入代码产权可信互用平台',
            type: 'success',
            showClose: false,
            offset: 100
          });
        } else { // 管理员登录
          await this.$router.push('/adminHome')
          this.$notify({
            title: '登录成功',
            message: '欢迎来到进入代码产权管理后台',
            type: 'success',
            showClose: false,
            offset: 100
          });
        }
      }
    },
    submitDrawer(formName) { // 提交抽屉中的注册信息
      this.$refs[formName].validate((valid) => {
        if (valid) {
          this.register()
        } else {
          this.$message({
            message: '请先填写完成注册信息',
            type: 'warning',
            center:true
          });
          return false;
        }
      });
    },
    async register(){ // 注册信息提交
      const {data:res}=await this.$http.post('/user/register',{
        username:this.form.username,
        pwd:this.form.pwd,
        name:this.form.name,
        role:this.form.role,
        tel:this.form.tel,
        idnum:this.form.idnum,
        addr:this.form.addr,
        gender:this.form.gender
      })
      if (res.code===0){
        this.$message.error(res.msg);
      } else {
        this.dialog = false;
        this.$refs.form.resetFields();
        this.$message.success(res.msg);
      }
    },
    cancelForm() { // Element抽屉的取消按钮
      this.dialog = false;
      clearTimeout(this.timer);
    },
    resetForm(formName) { // 重置按钮
      this.$refs[formName].resetFields();
    },
  }
}
</script>

<style scoped>
#elrowform{
  top: 190px;
}
#elrowcol{
  background: #eff3f8;
  opacity: inherit;
  border-radius: 10px;
}
/*#buttonitem{*/
/*  margin-left: 65px;*/
/*  margin-top: 20px;*/
/*  width: 90px;*/
/*}*/
/*#resetitem{*/
/*  margin-left: 38px;*/
/*  width: 90px;*/
/*}*/
#drawerForm{
  margin-right: 30px;
  margin-left: 15px;
}
.grid-content {
  margin: 20px 40px 5px 38px;
  padding: 15px;
}
#building{
  background: url("../assets/backgroundPage.JPG");
  height: 100vh; /* 100%视窗高度，可根据实际需要调整 */
  background-size: cover;
  background-position: center;
}
</style>
