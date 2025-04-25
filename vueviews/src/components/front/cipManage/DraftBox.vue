<template>
  <div>
    <el-row id="breadcrumb">
      <el-breadcrumb separator-class="el-icon-arrow-right">
        <el-breadcrumb-item style="font-size: small">成果管理</el-breadcrumb-item>
        <el-breadcrumb-item>草稿箱</el-breadcrumb-item>
      </el-breadcrumb>
    </el-row>
    <el-row id="main_range">
      <el-row style="margin-top: 40px;">
        <el-col :span="22" :offset="1">
          <el-card shadow="hover" class="box-card">
            <el-table
                :data="drafts"
                border
                style="margin-top: 10px;width: 100%;">
              <el-table-column
                  prop="filename"
                  label="名称"
                  align="center"
                  width="150">
              </el-table-column>
              <el-table-column
                  prop="version"
                  label="版本号"
                  align="center">
              </el-table-column>
              <el-table-column
                  prop="tag"
                  label="标签"
                  align="center">
                <template slot-scope="scope">
                  <el-tag type="success" effect="plain" disable-transitions>{{scope.row.tag}}</el-tag>
                </template>
              </el-table-column>
              <el-table-column
                  prop="author"
                  label="创作者"
                  align="center"
                  width="100">
              </el-table-column>
              <el-table-column
                  prop="language"
                  label="程序设计语言"
                  width="150"
                  align="center"
                  :filters="[{ text: 'Golang', value: 'Golang' }, { text: 'Java', value: 'Java' }, { text: 'Python', value: 'Python' }]"
                  :filter-method="filterLanguage"
                  filter-placement="bottom-end">
                <template slot-scope="scope">
                  <el-tag
                      :type="scope.row.language === 'Golang' ? 'warning' : (scope.row.language === 'Java' ? 'danger' : 'primary')"
                      disable-transitions>{{scope.row.language}}</el-tag>
                </template>
              </el-table-column>
              <el-table-column
                  prop="description"
                  label="组件描述"
                  align="center"
                  width="250">
              </el-table-column>
              <el-table-column
                  label="操作"
                  align="center"
                  width="200">
                <template slot-scope="scope">
                  <el-button
                      size="mini"
                      type="primary"
                      @click="handleEdit(scope.row.id)">编辑详情
                  </el-button>
                  <el-button
                      size="mini"
                      type="danger"
                      @click="handleDelete(scope.row.id)">删除
                  </el-button>
                </template>
              </el-table-column>
            </el-table>
            <!--分页-->
            <div style="margin-top:10px; display: flex; justify-content: center; align-items: flex-end;">
              <el-pagination
                  v-if="drafts.length > pageSize"
                  @current-change="handlePageChange"
                  :current-page="currentPage"
                  :page-sizes="[8]"
                  :page-size="pageSize"
                  layout="total, sizes, prev, pager, next, jumper"
                  :total="drafts.length"
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
  name: "DraftBox",
  data() {
    return {
      drafts: [],
      currentPage: 1, // 当前页码
      pageSize: 10, // 每页显示条数
    }
  },
  created() {
    this.queryDraftBox();
  },
  methods:{
    filterLanguage(value, row) {
      return row.language === value;
    },
    async queryDraftBox(){
      const {data:res}=await this.$http({
        url:'/code/draftBox',
        method:'GET',
        headers:{
          "Authorization":'Bearer '+localStorage.getItem("userLogin_jwtToken")
        }
      })
      if(res.code===0){
        // this.$message.error("您的草稿箱无内容");
      } else {
        this.drafts=res.data;
      }
    },

    handleEdit(id) {
        this.gotoEdit(id);
    },
    gotoEdit(id) {
      this.$router.push({name: 'DraftDetails', params: {id}})
    },
    handleDelete(id){
      this.$confirm('此操作将永久删除该组件代码，无法恢复，是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        this.deleteDraft(id);
      });
    },
    async deleteDraft(id){
      const {data:res}=await this.$http({
        url:'/code/deleteDraft',
        method:'DELETE',
        headers:{
          "Authorization":'Bearer '+localStorage.getItem("userLogin_jwtToken")
        },
        data: {
          id:id
        }
      })
      if(res.code===0){
        this.$notify.error({
          title: '错误',
          message: res.msg
        });
      } else {
        this.$notify({
          title: '成功',
          message: '删除成功',
          type: 'success'
        });
        await this.queryDraftBox();
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
