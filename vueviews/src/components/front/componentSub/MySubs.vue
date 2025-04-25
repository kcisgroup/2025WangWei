<template>
  <div>
    <el-row id="breadcrumb">
      <el-breadcrumb separator-class="el-icon-arrow-right">
        <el-breadcrumb-item style="font-size: small">订阅中心</el-breadcrumb-item>
        <el-breadcrumb-item>订阅申请记录</el-breadcrumb-item>
      </el-breadcrumb>
    </el-row>
    <el-row id="main_range">
      <el-row style="margin-top: 40px;">
        <el-col :span="22" :offset="1">
          <el-card shadow="hover" class="box-card">
            <el-table
                :data="subs"
                default-sort = "{prop: 'date', order: 'descending'}"
                border
                style="width: 100%">
              <el-table-column
                  prop="filename"
                  label="组件名"
                  width="150"
                  align="center">
              </el-table-column>
              <el-table-column
                  prop="version"
                  label="组件版本号"
                  width="120"
                  align="center">
              </el-table-column>
              <el-table-column
                  prop="partyB"
                  label="乙方联系人"
                  width="110"
                  align="center">
              </el-table-column>
              <el-table-column
                  prop="startTime"
                  label="订阅生效时间"
                  width="180"
                  sortable
                  align="center">
              </el-table-column>
              <el-table-column
                  prop="realPrice"
                  label="意向价格（元/天）"
                  width="180"
                  align="center">
              </el-table-column>
              <el-table-column
                  prop="note"
                  label="备注"
                  width="200"
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
                  width="221"
                  align="center">
                <template slot-scope="scope">
                  <el-tag
                      :type="selectStatus(scope.row.status)"
                      disable-transitions
                      effect="plain">{{scope.row.status}}</el-tag>
                </template>
              </el-table-column>

              <el-table-column
                  label="操作"
                  align="center"
                  width="260">
                <template slot-scope="scope">
                  <div class="button-container">
                    <el-button
                        size="mini"
                        type="primary"
                        v-if="scope.row.status==='甲方申请订阅，待乙方审核'"
                        @click="goToDetails(scope.row.codeId)"
                    >
                      组件详情
                    </el-button>
                    <el-button
                        size="mini"
                        type="danger"
                        v-if="scope.row.status==='甲方申请订阅，待乙方审核'"
                        @click="cancelSub(scope.row.id)"
                    >
                      取消订阅
                    </el-button>
                    <el-button
                        size="mini"
                        type="info"
                        v-if="scope.row.status==='甲方取消订阅申请' || scope.row.status === '乙方已拒绝' || scope.row.status === '甲方拒绝订阅条款'"
                        @click="cancelSub(scope.row.id)"
                        disabled
                    >
                      交易关闭
                    </el-button>
                    <el-button
                        size="mini"
                        type="primary"
                        v-if="scope.row.status==='乙方同意，待甲方确认订阅条款'"
                        @click="handleConfirm(scope.row.id)"
                    >
                      前往确认
                    </el-button>
                    <el-button
                        size="mini"
                        type="primary"
                        v-if="scope.row.status==='双方已签订合同'"
                        @click="handleUse(scope.row.dockerId)"
                    >
                      前往使用
                    </el-button>
                      <!-- Todo 暂不支付支付功能 -->
                    <el-button
                        size="mini"
                        type="success"
                        v-if="scope.row.status==='费用已清算但未支付'"
                        @click="handlePay(scope.row.dockerId)"
                        disabled
                    >
                      支付
                    </el-button>
                    <el-button
                        size="mini"
                        type="success"
                        v-if="scope.row.status==='费用已支付，本次交易完成'"
                        disabled
                    >
                        已支付
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
  name: "MySubs",
  data(){
    return{
      subs: [],
      currentPage: 1, // 当前页码
      pageSize: 8, // 每页显示条数
    };
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
        url:'/subscribe/queryByPartyA',
        method:'GET',
        headers:{
          "Authorization":'Bearer '+localStorage.getItem("userLogin_jwtToken")
        }
      })
      if(res.code===1){
        this.subs=res.data;
      }
    },
    // 跳转到详情页面
    goToDetails(id) {
      this.$router.push({ name: 'MySubCipDetails', params: { id } });
    },
    handlePageChange(newPage) {
      this.currentPage = newPage;
    },
    async cancelSub(id){
      const {data:res}=await this.$http({
        url:'/subscribe/cancel',
        method:'PUT',
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
          await this.getAllSubs();
      }
    },
    handleConfirm(id){
      this.$router.push({ name: 'ConfirmSub', params: { id } });
    },
    handleUse(dockerId){
      if (dockerId===0) {
        this.$router.push({ name: 'ServiceCreate' });
      } else {
        this.$router.push({ name: 'ServiceList' });
      }
    }
  }
}
</script>

<style scoped>

</style>
