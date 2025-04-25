<template>
  <div>
    <el-row id="breadcrumb">
      <el-breadcrumb separator-class="el-icon-arrow-right">
        <el-breadcrumb-item style="font-size: small">组件服务</el-breadcrumb-item>
        <el-breadcrumb-item>服务列表</el-breadcrumb-item>
      </el-breadcrumb>
    </el-row>
    <el-row id="main_range">
      <el-col :span="22" :offset="1" style="margin-top: 20px;">
        <el-card class="box-card" shadow="hover">
          <el-table
              :data="services"
              style="width: 100%">
            <el-table-column
                prop="filename"
                label="名称"
                width="150"
                align="center">
            </el-table-column>
            <el-table-column
                prop="version"
                label="版本"
                width="100"
                align="center">
            </el-table-column>
            <el-table-column
                prop="createdTime"
                label="创建时间"
                width="180"
                align="center">
            </el-table-column>
            <el-table-column
              prop="status"
              label="状态"
              width="180"
              align="center"
              :filters="[{ text: '运行中', value: '服务正在运行中' }, { text: '已停止', value: '服务已暂停' }, { text: '待消毁', value: '服务终止待消毁' }, { text: '已销毁', value: '服务已销毁' }]"
              :filter-method="filterStatus"
              filter-placement="bottom-end">
              <template slot-scope="scope">
                  <el-tag
                          :type="getStatusType(scope.row.status)"
                          disable-transitions>{{ getStatusText(scope.row.status) }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column
                prop="ip"
                label="IP地址"
                width="180"
                align="center">
                <template slot-scope="scope">
                    {{ scope.row.status === '服务已销毁' ? '-' : scope.row.ip }}
                </template>
            </el-table-column>
            <el-table-column
                prop="port"
                label="端口号"
                width="100"
                align="center">
                <template slot-scope="scope">
                    {{ scope.row.status === '服务已销毁' ? '-' : scope.row.port }}
                </template>
            </el-table-column>

            <el-table-column
                label="操作"
                width="260"
                align="center">
              <template slot-scope="scope">
                <div class="button-container">
                  <el-button
                      v-if="scope.row.status === '服务已暂停'"
                      size="mini"
                      type="primary"
                      @click="handleOpen(scope.row.id)">启动
                  </el-button>
                  <el-button
                          v-if="scope.row.status === '服务正在运行中'"
                          size="mini"
                          type="primary"
                          @click="handleStop(scope.row.id)">停止
                  </el-button>
                  <el-button
                          v-if="scope.row.status !== '服务已销毁'"
                          size="mini"
                          type="primary"
                          @click="queryDetails(scope.row.id)">详情
                  </el-button>
                  <el-button
                      v-if="scope.row.status === '服务正在运行中'"
                      size="mini"
                      type="danger"
                      @click="handleX(scope.row.id,scope.row.status)">终止订阅
                  </el-button>
                   <el-button
                      v-if="scope.row.status === '服务已暂停'"
                      size="mini"
                      type="danger"
                      @click="handleX(scope.row.id,scope.row.status)">终止订阅
                  </el-button>
                   <el-button
                      v-if="scope.row.status === '服务已销毁'"
                      size="mini"
                      type="info"
                      disabled>已销毁
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
  name: "ServiceList",
  data(){
    return{
      services:[],
      loading: false
    }
  },
  created() {
    this.queryService();
  },
  methods:{
    filterStatus(value, row) {
      return row.status === value;
    },
    getStatusType(status) {
        switch (status) {
            case '服务正在运行中':
                return 'success';
            case '服务已暂停':
                return 'danger';
            default:
                return 'info';
        }
    },
    getStatusText(status) {
        switch (status) {
            case '服务正在运行中':
                return '运行中';
            case '服务已暂停':
                return '已停止';
            case '服务终止待消毁':
                return '服务终止';
            default:
                return '已销毁';
        }
    },
    async queryService(){
      const {data:res}=await this.$http({
        url:'/product/queryByOwner',
        method:'GET',
        headers:{
          "Authorization":'Bearer '+localStorage.getItem("userLogin_jwtToken")
        }
      })
      if(res.code===1) {
          this.services = res.data;
      }
    },
    async handleStop(id) {
      const loading = this.$loading({
          lock: true,
          text: '正在停止，请耐心等待',
          spinner: 'el-icon-loading',
          background: 'rgba(0, 0, 0, 0.7)'
      });
      const {data: res} = await this.$http({
        url: '/product/stop',
        method: 'PUT',
        headers: {
          "Authorization": 'Bearer ' + localStorage.getItem("userLogin_jwtToken")
        },
        data: {
          id: id
        }
      })
      if (res.code === 0) {
        this.$message.error(res.msg);
      } else {
        this.$message.success(res.msg);
        await this.queryService();
      }
      loading.close();
    },
    async handleOpen(id) {
      const loading = this.$loading({
          lock: true,
          text: '正在启动，请耐心等待',
          spinner: 'el-icon-loading',
          background: 'rgba(0, 0, 0, 0.7)'
      });
      const {data: res} = await this.$http({
        url: '/product/start',
        method: 'PUT',
        headers: {
          "Authorization": 'Bearer ' + localStorage.getItem("userLogin_jwtToken")
        },
        data: {
          id: id
        }
      })
      if (res.code === 0) {
        this.$message.error(res.msg);
      } else {
        this.$message.success(res.msg);
        await this.queryService();
      }
      loading.close();
    },
    queryDetails(id) {
        this.$router.push({ name: 'ServiceDetails', params: { id } });
    },
    async handleX(id, status) {
      try {
        await this.$confirm('此操作将终止订阅，是否继续?', '提示', {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning'
        });

        const loading = this.$loading({
            lock: true,
            text: '正在结算，请耐心等待',
            spinner: 'el-icon-loading',
            background: 'rgba(0, 0, 0, 0.7)'
        });
        // 如果服务正在运行中，先停止服务
        if (status === '服务正在运行中') {
            const stopRes = await this.$http({
                url: '/product/stop',
                method: 'PUT',
                headers: {
                    "Authorization": 'Bearer ' + localStorage.getItem("userLogin_jwtToken")
                },
                data: {
                    id: id
                }
            });
            if (stopRes.data.code === 0) {
                loading.close();
                this.$message.error(stopRes.data.msg);
                return;
            }
        }
        // 终止订阅，结算费用
        const terminateRes = await this.$http({
            url: '/subscribe/terminate',
            method: 'PUT',
            headers: {
                "Authorization": 'Bearer ' + localStorage.getItem("userLogin_jwtToken")
            },
            data: {
                dockerId: id
            }
        });
        if (terminateRes.data.code === 0) {
            this.$message.error(terminateRes.data.msg);
        } else {
            this.$message.success(terminateRes.data.msg);
        }
        loading.close();
        await this.queryService();
    } catch (error) {
        // 捕获用户取消确认框或其他错误
        this.$message.error('操作已取消或发生错误');
        loading.close();
      }
    }
  }
}
</script>

<style scoped>

</style>