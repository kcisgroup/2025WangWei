<template>
  <div>
    <el-row id="breadcrumb">
      <el-breadcrumb separator-class="el-icon-arrow-right">
        <el-breadcrumb-item style="font-size: small">首页</el-breadcrumb-item>
      </el-breadcrumb>
    </el-row>
    <el-row id="main_range">
      <div>
        <el-row type="flex" :gutter="20">
          <el-col :span="5" :offset="2">
            <div class="nk">
              <el-statistic group-separator="," :value="nums.totalUsers">
                <template slot="title">
                  <span id="elEtatistic">平台用户数量</span>
                </template>
              </el-statistic>
            </div>
          </el-col>

          <el-col :span="5">
            <div class="nk">
              <el-statistic group-separator="," :value="nums.totalDrafts">
                <template slot="title">
                  <span id="elEtatistic">上传项目总数</span>
                </template>
              </el-statistic>
            </div>
          </el-col>

          <el-col :span="5">
            <div class="nk">
              <el-statistic group-separator="," :value="nums.totalCips">
                <template slot="title">
                  <span id="elEtatistic">授权项目总数</span>
                </template>
              </el-statistic>
            </div>
          </el-col>

          <el-col :span="5">
            <div class="nk">
              <el-statistic group-separator="," :value="nums.totalDockers">
                <template slot="title">
                  <span id="elEtatistic">数字产品使用量</span>
                </template>
              </el-statistic>
            </div>
          </el-col>
        </el-row>

        <el-row type="flex" :gutter="30">
          <el-col :span="6" :offset="2">
            <div>
              <el-card shadow="hover" style="margin-top: 100px;width: 100%;">
                <div slot="header" class="clearfix" style="text-align: center">
                  <span style="font-size: large;">授权组件占比</span>
                </div>
                <div style="margin-top: 35px;text-align: center">
                  <el-progress type="circle" :color="customColor" :percentage="cipPer" width="180"></el-progress>
                </div>
                <div style="margin-top: 35px;"></div>
              </el-card>
            </div>
          </el-col>

          <el-col :span="6" :offset="1">
            <div>
              <el-card shadow="hover" style="margin-top: 100px;width: 100%;">
                <div slot="header" class="clearfix" style="text-align: center">
                  <span style="font-size: large">开放产权占比</span>
                </div>

                <div style="margin-top: 35px;text-align: center">
                  <el-progress type="circle" :color="customColor2" :percentage="cipPer2" width="180"></el-progress>
                </div>
                <div style="margin-top: 35px;"></div>
              </el-card>
            </div>
          </el-col>

          <el-col :span="6" :offset="1">
            <div>
              <el-card shadow="hover" style="margin-top: 100px;width: 100%;">
                <div slot="header" class="clearfix" style="text-align: center">
                  <span style="font-size: large;">数字产品转化率</span>
                </div>
                <div style="margin-top: 35px;text-align: center">
                  <el-progress type="circle" :color="customColor3" :percentage="cipPer3" width="180"></el-progress>
                </div>
                <div style="margin-top: 35px;"></div>
              </el-card>
            </div>
          </el-col>
        </el-row>
      </div>
    </el-row>
  </div>
</template>

<script>
export default {
  name: "UserHome",
  data() {
    return {
      like: true,
      nums:{},
      customColor:"#67C239",
      customColor2:"#409EFE",
      customColor3:"#F66C6C"
    };
  },
    computed: {
    cipPer() {
        return (this.nums.totalCips / this.nums.totalDrafts * 100).toFixed(1);
    },
    cipPer2() {
        return (this.nums.totalOpenCips / this.nums.totalCips * 100).toFixed(1);
    },
    cipPer3() {
        return (this.nums.totalDockers / this.nums.totalCips * 100).toFixed(1);
    },
  },
  created() {
    this.getCount();
  },
  methods: {
    async getCount() {
      const {data:res}=await this.$http({
        url:'/code/count',
        method:'GET',
        headers:{
            "Authorization":'Bearer '+localStorage.getItem("userLogin_jwtToken")
        }
      })
      this.nums=res.data;
        console.log(res.data);
    }
  }
}
</script>

<style scoped>
#elEtatistic{
  display: block;
  margin-bottom: 20px;
  font-size: large;
  font-weight: bolder;
}
.nk{
    margin-top: 50px;
}
</style>
