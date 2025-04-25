<template>
  <div>
    <el-row id="breadcrumb">
      <el-breadcrumb separator-class="el-icon-arrow-right">
        <el-breadcrumb-item style="font-size: small">订阅中心</el-breadcrumb-item>
        <el-breadcrumb-item style="font-size: small" :to="{ name: 'MySubs' }">订阅申请记录</el-breadcrumb-item>
        <el-breadcrumb-item>我的订阅详情</el-breadcrumb-item>
      </el-breadcrumb>
    </el-row>
    <el-row id="main_range">
      <br>

      <el-row>
        <el-col :span="18" :offset="3">
          <el-card style="height: 680px;" class="box-card">
            <div slot="header" class="clearfix">
              <span>组件详情</span>
              <el-button style="float: right; padding: 3px 0" type="text" @click="backToMySubs">返回</el-button>
            </div>
            <div>
              <el-col :span="howWidth" v-for="(item, index) in projectInfoArr" :key="index">
                <div class="box">
                  <div class="content1">{{ item.key }}</div>
                  <div class="content2">{{ item.value == "" ? "待完善" : item.value }}</div>
                </div>
              </el-col>
            </div>
          </el-card>
          <br>
        </el-col>
      </el-row>
    </el-row>
  </div>
</template>

<script>
export default {
  name: "MySubCipDetails",
  props: ['id'],
  data(){
    return{
      projectInfoArr:[
        {
          label:'filename',
          key: '项目名',
          value: '',
        },
        {
          label:'version',
          key: '项目版本号',
          value: '',
        },
        {
          label:'tag',
          key: '项目标签',
          value: '',
        },
        {
          label:'regTime',
          key: '项目上线时间',
          value: '',
        },
        {
          label:'updatedTime',
          key: '最新修改时间',
          value: '',
        },
        {
          label:'price',
          key: '意向单价',
          value: '',
        },
        {
          label:'name',
          key: '创作者',
          value: '',
        },
        {
          label:'role',
          key: '创作者单位',
          value: '',
        },
        {
          label:'tel',
          key: '创作者联系方式',
          value: '',
        },
        {
          label:'language',
          key: '编程语言',
          value: '',
        },
        {
          label:'description',
          key: '功能描述',
          value: '',
        }
      ],
      projectInfo:{},
      howWidth: 20
    }
  },
  created() {
    this.getDetails();
  },
  methods: {
    async getDetails() {
      const {data: res} = await this.$http({
        url: '/code/queryCipInfo',
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
        this.projectInfo = res.data;
        this.assignCipValues();
      }
    },
    // 将对象各字段赋值至数组对应位置
    assignCipValues() {
      for (const item of this.projectInfoArr) {
        const key = item.label;
        // eslint-disable-next-line no-prototype-builtins
        if (this.projectInfo.hasOwnProperty(key)) {
          // eslint-disable-next-line no-empty
          item.value = this.projectInfo[key];
        }
      }
    },
    backToMySubs(){
      this.$router.push({ name: 'MySubs'});
    }
  }
}
</script>

<style scoped>
.box {
  width: 100%;
  height: 50px;
  display: flex;
  margin-left: 90px;
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
</style>
