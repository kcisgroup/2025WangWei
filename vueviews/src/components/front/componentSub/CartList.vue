<template>
  <div>
    <el-row id="breadcrumb">
      <el-breadcrumb separator-class="el-icon-arrow-right">
        <el-breadcrumb-item style="font-size: small">订阅中心</el-breadcrumb-item>
        <el-breadcrumb-item>我的购物车</el-breadcrumb-item>
      </el-breadcrumb>
    </el-row>
    <el-row id="main_range">
      <el-col :span="22" :offset="1" style="margin-top: 20px;">
        <el-card class="box-card" shadow="hover">
          <el-table
              :data="carts"
              style="width: 100%">
            <el-table-column label="组件信息" align="center">
              <el-table-column
                  prop="filename"
                  label="名称"
                  width="180"
                  align="center">
              </el-table-column>
              <el-table-columnx
                  prop="version"
                  label="版本号"
                  width="100"
                  align="center">
              </el-table-columnx>
              <el-table-column
                  prop="tag"
                  label="标签"
                  width="100"
                  align="center">
              </el-table-column>
              <el-table-column
                  prop="description"
                  label="组件描述"
                  align="center"
                  width="300">
              </el-table-column>
            </el-table-column>
            <el-table-column
                prop="price"
                label="价格（元/天）"
                width="180"
                align="center">
            </el-table-column>
            <el-table-column
                prop="createdTime"
                label="加入时间"
                width="180"
                align="center">
            </el-table-column>
            <el-table-column
                label="操作"
                align="center"
                width="230">
              <template slot-scope="scope">
                  <div class="button-container">
                    <el-button
                      size="mini"
                      type="primary"
                      @click="handleSub(scope.row.id)">订阅
                    </el-button>
                    <el-button
                      size="mini"
                      type="danger"
                      @click="handleDelete(scope.row.id)">移除
                    </el-button>
                  </div>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script>
export default {
  name: "CartList",
  data(){
    return{
      carts:[]
    }
  },
  created() {
    this.queryCart();
  },
  methods:{
    async queryCart(){
      const {data:res}=await this.$http({
        url:'/subscribe/queryCart',
        method:'GET',
        headers:{
          "Authorization":'Bearer '+localStorage.getItem("userLogin_jwtToken")
        }
      })
      if(res.code===0){
        this.$message.error(res.msg);
      } else {
        this.carts=res.data;
      }
    },
    handleSub(id) {
      this.$router.push({ name: 'SubscribeCip', params: { id } });
    },
    async handleDelete(id){
      const {data:res}=await this.$http({
        url:'/subscribe/removeCart',
        method:'DELETE',
        headers:{
          "Authorization":'Bearer '+localStorage.getItem("userLogin_jwtToken")
        },
        data:{
          id:id
        }
      })
      if(res.code===0){
        this.$message.error(res.msg);
      } else {
          await this.queryCart();
      }
    }
  }
}
</script>

<style scoped>

</style>
