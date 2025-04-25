<template>
  <div>
    <el-row id="breadcrumb">
      <el-breadcrumb separator-class="el-icon-arrow-right">
        <el-breadcrumb-item style="font-size: small">成果管理</el-breadcrumb-item>
        <el-breadcrumb-item :to="{ name: 'MyCips' }">代码产权查验</el-breadcrumb-item>
        <el-breadcrumb-item>产权详情</el-breadcrumb-item>
      </el-breadcrumb>
    </el-row>
    <el-row id="main_range">
      <div id="detailsTable">
        <el-col :span="howWidth" v-for="(item, index) in cipArr" :key="index">
          <div class="box">
            <div class="content1">{{ item.key }}</div>
            <div class="content2">{{ item.value === "" ? "待完善" : item.value }}</div>
          </div>
        </el-col>
      </div>
      <!--修改功能，嵌套form的抽屉-->
      <el-row :gutter="20" type="flex" align="middle" justify="center">
        <el-col :span="4" :offset="10" style="margin-bottom: 30px;">
          <el-button
            type="primary"
            @click="dialog = true"
            style="width: 100%"
          >
            修改产权<i class="el-icon-edit el-icon--right"></i>
          </el-button>
        </el-col>
      </el-row>

      <el-drawer
          :visible.sync="dialog"
          direction="rtl"
          custom-class="demo-drawer"
          ref="drawer"
          size="30%">
        <div>
          <el-row type="flex" justify="center">
              <el-col :span="18">
                  <div class="drawer-title">
                      代码产权信息变更
                  </div>
              </el-col>
          </el-row>

          <el-form :model="updateCipForm" :rules="rules" ref="updateCipForm" label-width="110px">
              <el-row type="flex" justify="center">
                <el-col :span="20">
                  <!-- 订阅价格 -->
                  <el-form-item label="订阅价格" prop="price" style="margin-top: 20px;">
                    <el-input
                      v-model="updateCipForm.price"
                      placeholder="请设置新的价格"
                      style ="width: 100%;"
                    ></el-input>
                  </el-form-item>
                </el-col>
              </el-row>

              <el-row type="flex" justify="center">
                <el-col :span="20">
                  <!-- 是否开放订阅 -->
                  <el-form-item label="是否开放订阅" prop="status">
                    <el-select
                      v-model="updateCipForm.status"
                      placeholder="请选择是否开放订阅"
                      style="width: 100%;"
                    >
                      <el-option label="开放" value="开放"></el-option>
                      <el-option label="禁用" value="禁用"></el-option>
                    </el-select>
                  </el-form-item>
                </el-col>
              </el-row>

                  <!-- 按钮组 -->
                    <el-form-item>
                        <el-row type="flex" justify="left" :gutter="20">
                            <!-- 提交按钮 -->
                            <el-col :span="8">
                                <el-button
                                        type="primary"
                                        @click="submitUpdateCipForm('updateCipForm')"
                                        class="form-button"
                                        style ="width: 100%;"
                                >
                                    提交
                                </el-button>
                            </el-col>

                            <!-- 重置按钮 -->
                            <el-col :span="8">
                                <el-button
                                        @click="resetForm('updateCipForm')"
                                        class="form-button"
                                        style ="width: 100%;"
                                >
                                    重置
                                </el-button>
                            </el-col>
                        </el-row>
                      </el-form-item>
            </el-form>
        </div>
      </el-drawer>
    </el-row>
  </div>
</template>

<script>
export default {
  name: "MyCipDetails",
  props: ['id'],
  data(){
    return{
      cipArr:[
        {
          label:'filename',
          key: '项目名',
          value: '',
        },
        {
          label:'version',
          key: '版本号',
          value: '',
        },
        {
          label:'author',
          key: '创作者',
          value: '',
        },
        {
          label:'cpi',
          key: '代码产权数字摘要',
          value: '',
        },
        {
          label:'createdTime',
          key: '项目初次上传时间',
          value: '',
        },
        {
          label:'regTime',
          key: '代码产权确立时间',
          value: '',
        },
        {
          label:'updatedTime',
          key: '最近更新时间',
          value: '',
        },
        {
          label:'language',
          key: '项目程序设计语言',
          value: '',
        },
        {
          label:'price',
          key: '订阅价格（元/天）',
          value: '',
        },
        {
          label:'status',
          key: '代码产权当前开放程度',
          value: '',
        },
        {
          label:'addr',
          key: '项目文件存储地址',
          value: '',
        },
        {
          label:'description',
          key: '项目描述',
          value: '',
        }
      ],
      cipObject:{},
      howWidth: 20,
      dialog: false,
      updateCipForm: {
        price: '',
        status: ''
      },
      rules: {
        price: [
          { required: true, message: '请设置订阅单价', trigger: 'blur' },
        ],
        status: [
          { required: true, message: '请选择是否开放订阅', trigger: 'change' }
        ]
      }
    }
  },
  created() {
    this.queryDetails();
  },
  methods:{
    async queryDetails(){
      const {data:res}=await this.$http({
        url:'/code/queryDetails',
        method:'POST',
        headers:{
          "Authorization":'Bearer '+localStorage.getItem("userLogin_jwtToken")
        },
        data:{
          id:this.id
        }
      })
      if(res.code===0){
        this.$message.error(res.msg);
        await this.$router.push({name: 'MyCips'});
      } else {
        this.cipObject=res.data; // 接收对象参数并赋值cipObject对象
        this.assignCipValues();
      }
    },
    // 将对象各字段赋值至数组对应位置
    assignCipValues() {
      for (const item of this.cipArr) {
        const key = item.label;
        if (this.cipObject.hasOwnProperty(key)) {
          item.value = this.cipObject[key];
        }
      }
    },
    submitUpdateCipForm(formName) {
      this.$refs[formName].validate((valid) => {
        if (valid) {
          const loading = this.$loading({
            lock: true,
            text: '正在提交，请耐心等待，请勿关闭',
            spinner: 'el-icon-loading',
            background: 'rgba(0, 0, 0, 0.7)'
          });
          this.updateCip().then(()=>{
            loading.close();
            this.$router.push({ name: 'MyCips'});
          })
        } else {
          this.$message({
            message: '请将信息填写完整',
            type: 'warning',
            center:true
          });
          return false;
        }
      });
    },
    async updateCip(){
      const {data:res}=await this.$http({
        url:'/code/updateCip',
        method:'PUT',
        headers:{
          "Authorization":'Bearer '+localStorage.getItem("userLogin_jwtToken")
        },
        data:{
          cpi:this.cipObject['cpi'],
          price:this.updateCipForm.price,
          status:this.updateCipForm.status
        }
      })
      if(res.code===0){
        this.$message.error(res.msg);
      } else {
        this.$message({
          message: '产权信息修改成功',
          type: 'success'
        });
      }
    },
    resetForm(formName) {
      this.$refs[formName].resetFields();
    }
  }
}
</script>

<style scoped>
#detailsTable {
  width: 100%;
  box-sizing: border-box;
  padding: 50px;
}
.box {
  width: 100%;
  height: 50px;
  display: flex;
  margin-left: 120px;
  border-left: 1px solid #e9e9e9;
  border-top: 1px solid #e9e9e9;
}
.content1 {
  width: 40%;
  height: 50px;
  line-height: 50px;
  text-align: center;
  background-color: #fafafa;
  border-right: 1px solid #e9e9e9;
  border-bottom: 1px solid #e9e9e9;
  color: #333;
  font-size: 17px;
}
.content2 {
  width: 100%;
  height: 50px;
  line-height: 50px;
  text-align: center;
  background-color: #fff;
  border-right: 1px solid #e9e9e9;
  border-bottom: 1px solid #e9e9e9;
  color: #7d7d80;
  font-size: 17px;
}
/*#edit{*/
/*  margin-top: 60px;*/
/*  font-size: medium;*/
/*  font-family: 'PingFang SC';*/
/*}*/
.drawer-title {
    font-family: 宋体;
    font-size: x-large;
    text-align: center;
}
.form-button {
    width: 100%;
    height: 40px;
    font-size: medium;
}
/*#drawerItem{*/
/*  margin-left: 20px;*/
/*  margin-top: 75px;*/
/*}*/
/*#drawerItem2 {*/
/*  margin-left: 20px;*/
/*  margin-top: 50px;*/
/*}*/
</style>
