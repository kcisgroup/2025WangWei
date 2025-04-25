<template>
  <div>
    <el-row id="breadcrumb">
      <el-breadcrumb separator-class="el-icon-arrow-right">
        <el-breadcrumb-item style="font-size: small">订阅中心</el-breadcrumb-item>
        <el-breadcrumb-item>我处理的订单</el-breadcrumb-item>
      </el-breadcrumb>
    </el-row>
    <el-row id="main_range">
      <br>
      <br>

      <el-row>
        <el-col :span="22" :offset="1">
          <el-card shadow="hover" class="box-card">
            <el-table
                :data="subs"
                default-sort = "{prop: 'date', order: 'descending'}"
                border
                style="width: 100%">
              <el-table-column
                  prop="partyA"
                  label="申请人"
                  width="150"
                  align="center">
              </el-table-column>
              <el-table-column
                  prop="cpi"
                  label="组件的代码产权数字摘要"
                  width="280"
                  align="center">
              </el-table-column>
              <el-table-column
                  prop="createdTime"
                  label="申请时间"
                  width="180"
                  sortable
                  align="center">
              </el-table-column>
              <el-table-column
                  prop="status"
                  label="当前状态"
                  :filters="[{ text: '甲方申请订阅，待乙方审核', value: '甲方申请订阅，待乙方审核' },{ text: '甲方取消订阅申请', value: '甲方取消订阅申请' },
                  { text: '乙方同意，待甲方确认订阅条款', value: '乙方同意，待甲方确认订阅条款' },{ text: '乙方已拒绝', value: '乙方已拒绝' },
                  { text: '双方已签订合同', value: '双方已签订合同' },{ text: '甲方拒绝订阅条款', value: '甲方拒绝订阅条款' },
                  { text: '费用已清算但未支付', value: '费用已清算但未支付' },{ text: '交易完成', value: '费用已支付，本次交易完成' }]"
                  :filter-method="filterStatus"
                  filter-placement="bottom-end"
                  width="230"
                  align="center">
                <template slot-scope="scope">
                  <el-tag
                      :type="selectStatus(scope.row.status)"
                      disable-transitions
                      effect="plain">{{scope.row.status}}</el-tag>
                </template>
              </el-table-column>
              <el-table-column
                  prop="note"
                  label="备注"
                  width="200"
                  align="center">
              </el-table-column>

              <el-table-column
                  label="操作"
                  width="260"
                  align="center">
                <template slot-scope="scope">
                  <div class="button-container">
                    <el-button
                        size="mini"
                        type="primary"
                        v-if="scope.row.status==='甲方申请订阅，待乙方审核'"
                        @click="goToSolve(scope.row.id,scope.row.codeId)"
                    >
                      前往处理
                    </el-button>
                    <el-button
                        size="mini"
                        type="info"
                        v-if="scope.row.status==='甲方取消订阅申请' || scope.row.status === '甲方拒绝订阅条款'"
                        disabled
                    >
                      交易关闭
                    </el-button>
                    <el-button
                        size="mini"
                        type="success"
                        v-if="scope.row.status==='双方已签订合同'"
                        disabled
                    >
                      已处理
                    </el-button>
                    <el-button
                        size="mini"
                        type="info"
                        v-if="scope.row.status==='乙方同意，待甲方确认订阅条款'"
                        disabled
                    >
                      待对方确认
                    </el-button>
                      <el-button
                              size="mini"
                              type="info"
                              v-if="scope.row.status==='费用已清算但未支付'"
                              disabled
                      >
                          等待对方支付
                      </el-button>
                      <el-button
                              size="mini"
                              type="success"
                              v-if="scope.row.status==='费用已支付，本次交易完成'"
                              disabled
                      >
                          交易完成
                      </el-button>
                  </div>
                </template>
              </el-table-column>
            </el-table>
            <!--分页-->
            <div style="margin-top:10px; display: flex; justify-content: center; align-items: flex-end;">
              <el-pagination
                  v-if="subs.length > pageSize"
                  @current-change="handlePageChange"
                  :current-page="currentPage"
                  :page-sizes="[8]"
                  :page-size="pageSize"
                  layout="total, sizes, prev, pager, next, jumper"
                  :total="subs.length"
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
  name: "SolveSub",
  data(){
    return{
      subs: [],
      currentPage: 1, // 当前页码
      pageSize: 8, // 每页显示条数
    }
  },
  created() {
    this.getAllSubs();
  },
  methods:{
    filterStatus(value, row) {
      return row.status === value;
    },
    selectStatus(status){
      if (status === '甲方取消订阅申请' || status === '乙方已拒绝'||status === '甲方拒绝订阅条款'||status === '费用已清算但未支付') {
        return 'danger'
      } else if (status === '双方已签订合同' || status === '费用已支付，本次交易完成') {
        return 'success'
      } else {
        return 'warning'
      }
    },
    async getAllSubs(){
      const {data:res}=await this.$http({
        url:'/subscribe/queryByPartyB',
        method:'GET',
        headers:{
          "Authorization":'Bearer '+localStorage.getItem("userLogin_jwtToken")
        }
      })
      if(res.code===1){
        this.subs=res.data;
      } else {
        this.$message.error(res.msg);
      }
    },
    // 跳转到详情页面
    goToSolve(id) {
      this.$router.push({ name: 'Solve', params: { id } });
    },
    handlePageChange(newPage) {
      this.currentPage = newPage;
    },
  }
}
</script>

<style scoped>

</style>
