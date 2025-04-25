<template>
  <div>
    <el-row id="breadcrumb">
      <el-breadcrumb separator-class="el-icon-arrow-right">
        <el-breadcrumb-item style="font-size: small">成果管理</el-breadcrumb-item>
        <el-breadcrumb-item>代码产权查验</el-breadcrumb-item>
      </el-breadcrumb>
    </el-row>
    <el-row id="main_range">
      <el-row style="margin-top: 40px;">
        <el-col :span="22" :offset="1">
          <el-card shadow="hover" class="box-card">
            <el-table
                :data="cipRecords"
                :default-sort = "{prop: 'date', order: 'descending'}"
                border
                style="width: 100%">
              <el-table-column
                  prop="filename"
                  label="产权名"
                  width="150"
                  align="center">
              </el-table-column>
              <el-table-column
                  prop="version"
                  label="版本号"
                  width="80"
                  align="center">
              </el-table-column>
              <el-table-column
                  prop="tag"
                  label="标签"
                  width="80"
                  align="center">
                <template slot-scope="scope">
                  <el-tag
                      type="success"
                      disable-transitions
                      effect="plain">{{scope.row.tag}}</el-tag>
                </template>
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
                  prop="regTime"
                  label="代码产权确立时间"
                  sortable
                  width="180"
                  align="center">
              </el-table-column>
              <el-table-column
                  prop="status"
                  label="是否开放订阅"
                  width="120"
                  :filters="[{ text: '禁用', value: '禁用' }, { text: '开放', value: '开放' }]"
                  :filter-method="filterStatus"
                  filter-placement="bottom-end"
                  align="center">
                <template slot-scope="scope">{{scope.row.status}}</template>
              </el-table-column>
              <el-table-column
                  prop="description"
                  label="项目描述"
                  width="250"
                  align="center">
              </el-table-column>
              <el-table-column
                  align="center"
                  width="225"
                  :resizable="false">
                <template slot="header" slot-scope="scope">
                  <el-input
                      v-model="filterFilename"
                      size="mini"
                      placeholder="请输入产权名搜索"
                      @input="search"/>
                </template>
                <template slot-scope="cipRecords">
                    <div class="button-container">
                      <el-button
                        size="mini"
                        type="primary"
                        @click="goToDetails(cipRecords.row.id)"
                    >
                      查看详情
                    </el-button>
                    <el-button
                        size="mini"
                        type="danger"
                        @click="subCheck(cipRecords.row.id)"
                    >
                      产权安全检查
                    </el-button>
                    </div>
                </template>
              </el-table-column>
            </el-table>
            <!--分页-->
            <div style="margin-top:10px; display: flex; justify-content: center; align-items: flex-end;">
              <el-pagination
                  v-if="cipRecords.length > pageSize"
                  @current-change="handlePageChange"
                  :current-page="currentPage"
                  :page-sizes="[10]"
                  :page-size="pageSize"
                  layout="total, sizes, prev, pager, next, jumper"
                  :total="cipRecords.length"
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
  name: "MyCips",
  data(){
    return{
      // filterLanguage: '',
      cipRecords: [],
      filterFilename:'',
      currentPage: 1, // 当前页码
      pageSize: 10, // 每页显示条数
    };
  },
  created() {
    this.getAllCip();
  },
  methods:{
    filterLanguage(value, row) {
      return row.language === value;
    },
    filterStatus(value, row) {
      return row.status === value;
    },
    async getAllCip(){ // 发送请求
      const {data:res}=await this.$http({
        url:'/code/queryByAuthor',
        method:'GET',
        headers:{
          "Authorization":'Bearer '+localStorage.getItem("userLogin_jwtToken")
        }
      })
      if(res.code===1){
        this.cipRecords=res.data;
      }
    },
    // 跳转到详情页面
    goToDetails(id) {
      // 使用 Vue Router 的编程式导航进行页面跳转
      this.$router.push({ name: 'MyCipDetails', params: { id } });
    },
    async checkCip(id){
      const {data:res}=await this.$http({
        url:'/code/check',
        method:'POST',
        headers:{
          "Authorization":'Bearer '+localStorage.getItem("userLogin_jwtToken")
        },
        data:{
          id:id
        }
      })
      if(res.code===0){
        this.$notify.error({
          title: '错误',
          message: '该产权存在问题'+res.msg,
          duration: 0
        });
      } else {
        this.$notify({
          title: '成功',
          message: res.msg,
          type: 'success'
        });
      }
    },
    subCheck(id){
      this.$confirm('此操作将检查代码产权信息有效性及产权文件完整性，是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        const loading = this.$loading({
          lock: true,
          text: '正在检查您的代码产权，请耐心等待，请勿关闭',
          spinner: 'el-icon-loading',
          background: 'rgba(0, 0, 0, 0.7)'
        });
        this.checkCip(id).then(() => {
          loading.close();
        });
      });
    },
    async search(){
      const {data:res}=await this.$http({
        url:'/code/queryCipByFilename',
        method:'POST',
        headers:{
          "Authorization":'Bearer '+localStorage.getItem("userLogin_jwtToken")
        },
        data:{
          filename:this.filterFilename
        }
      })
      if(res.code===0){
        this.$message({
          message: res.msg,
          type: 'warning',
          duration: 1000
        });
      } else {
        this.cipRecords=res.data;
      }
    },
    handlePageChange(newPage) {
      this.currentPage = newPage;
    },
  }
}
</script>

<style scoped>
.button-container {
    display: flex;
    overflow-x: auto;
    white-space: nowrap;
    gap: 8px; /* 按钮间距 */
    max-width: 280px; /* 略小于列宽 */
}

/* 隐藏滚动条（可选） */
.button-container::-webkit-scrollbar {
    display: none;
}
</style>
