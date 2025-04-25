<template>
  <div>
    <el-row id="breadcrumb">
      <el-breadcrumb separator-class="el-icon-arrow-right">
        <el-breadcrumb-item style="font-size: small">订阅中心</el-breadcrumb-item>
        <el-breadcrumb-item>浏览开放组件</el-breadcrumb-item>
      </el-breadcrumb>
    </el-row>
    <el-row id="main_range">
      <el-row style="margin-top: 40px;">
        <el-col :span="22" :offset="1">
          <el-card shadow="hover" class="box-card">
            <el-table
                :data="cips"
                border
                style="width: 100%">
              <el-table-column
                  prop="filename"
                  label="名称"
                  align="center"
                  width="180">
              </el-table-column>
              <el-table-column
                  prop="version"
                  label="版本号"
                  align="center"
                  width="120">
              </el-table-column>
              <el-table-column
                  prop="author"
                  label="创作者"
                  align="center"
                  width="120">
              </el-table-column>
              <el-table-column
                  prop="regTime"
                  label="确权日期时间"
                  width="180"
                  align="center">
              </el-table-column>
              <el-table-column
                  prop="language"
                  label="程序设计语言"
                  width="120"
                  :filters="[{ text: 'Golang', value: 'Golang' }, { text: 'Java', value: 'Java' }, { text: 'Python', value: 'Python' }]"
                  :filter-method="filterLanguage"
                  filter-placement="bottom-end"
                  align="center">
                <template slot-scope="scope">
                  <el-tag
                      :type="scope.row.language === 'Golang' ? 'warning' : (scope.row.language === 'Java' ? 'danger' : 'primary')"
                      disable-transitions>{{scope.row.language}}</el-tag>
                </template>
              </el-table-column>
              <el-table-column
                  prop="price"
                  label="订阅价格（元/天）"
                  align="center"
                  width="160">
              </el-table-column>
              <el-table-column
                  prop="description"
                  label="组件描述"
                  width="250"
                  align="center">
              </el-table-column>
              <el-table-column
                  label="操作"
                  width="260"
                  align="center">
                <template slot="header" slot-scope="scope">
                  <el-input
                      v-model="requiredFunction"
                      size="mini"
                      placeholder="请输入组件需求搜索"
                      @input="searchFunction"/>
                </template>
                <template slot-scope="cips">
                  <el-button
                      size="mini"
                      type="primary"
                      @click="goToCipRentDetails(cips.row.id)"
                  >
                    查看详情
                  </el-button>
                  <el-button
                      size="mini"
                      type="danger"
                      @click="addCart(cips.row.id)"
                  >
                    加入购物车
                  </el-button>
                </template>
              </el-table-column>
            </el-table>
            <!--分页-->
            <div style="margin-top:10px; display: flex; justify-content: center; align-items: flex-end;">
              <el-pagination
                  v-if="cips.length > pageSize"
                  @current-change="handlePageChange"
                  :current-page="currentPage"
                  :page-sizes="[8]"
                  :page-size="pageSize"
                  layout="total, sizes, prev, pager, next, jumper"
                  :total="cips.length"
              />
            </div>
          </el-card>
        </el-col>
      </el-row>
    </el-row>
  </div>
</template>

<script>
export default {
  name: "CipsOpened",
  data(){
    return{
      targetLanguage:'',
      cips:[],
      requiredFunction:'',
      currentPage: 1, // 当前页码
      pageSize: 8, // 每页显示条数
    }
  },
  created() {
    this.queryAllCips();
  },
  methods:{
    filterLanguage(value, row) {
      return row.language === value;
    },
    async queryAllCips(){
      const {data:res}=await this.$http({
        url:'/code/queryAllCips',
        method:'GET',
        headers:{
          "Authorization":'Bearer '+localStorage.getItem("userLogin_jwtToken")
        }
      })
      if(res.code===1) {
        this.cips=res.data; // 接收对象参数并赋值cips对象
        this.assignCipValues();
      }
    },
    goToCipRentDetails(id){
      this.$router.push({ name: 'CipOpenedDetails', params: { id } });
    },
    async searchFunction(){
      const {data:res}=await this.$http({
        url:'/code/queryByDesc',
        method:'POST',
        headers:{
          "Content-Type": "application/json", // 设置数据传输格式为application/json
          "Authorization":'Bearer '+localStorage.getItem("userLogin_jwtToken")
        },
        data: {
          description: this.requiredFunction
        }
      })
      if(res.code===0){
        this.$message({
          message: res.msg,
          type: 'warning',
          duration: 1000
        });
      } else {
        this.cips=res.data;
      }
    },
    async addCart(id){
      const {data:res}=await this.$http({
        url:'/subscribe/addCart',
        method:'POST',
        headers:{
          "Authorization":'Bearer '+localStorage.getItem("userLogin_jwtToken")
        },
        data: {
          codeId:id
        }
      })
      if(res.code===0){
        this.$message({
          message: res.msg,
          type: 'warning',
          duration: 1000
        });
      } else {
        this.$message({
          message: '成功加入购物车',
          type: 'success',
          duration: 500
        });
      }
    },
    handlePageChange(newPage) {
      this.currentPage = newPage;
    }
  }
}
</script>

<style scoped>

</style>
