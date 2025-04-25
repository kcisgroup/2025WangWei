<template>
  <div>
    <el-row id="breadcrumb">
      <el-breadcrumb separator-class="el-icon-arrow-right">
        <el-breadcrumb-item style="font-size: small">组件服务</el-breadcrumb-item>
        <el-breadcrumb-item>创建服务</el-breadcrumb-item>
      </el-breadcrumb>
    </el-row>
    <el-row id="main_range">
      <el-col :span="22" :offset="1" style="margin-top: 20px;">
        <el-card class="box-card" shadow="hover">
          <el-table
              :data="waitCons"
              style="width: 100%">
            <el-table-column
                prop="filename"
                label="名称"
                align="center">
            </el-table-column>
            <el-table-columnx
                prop="version"
                label="版本号"
                align="center">
            </el-table-columnx>
            <el-table-column
                prop="tag"
                label="标签"
                align="center">
            </el-table-column>
            <el-table-column
                prop="startTime"
                label="订阅生效时间"
                align="center">
            </el-table-column>
            <el-table-column
                prop="language"
                label="编程语言"
                align="center">
            </el-table-column>
            <el-table-column
                prop="cur_status"
                label="状态"
                align="center"
                :formatter="statusFormatter">
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
                      @click="handleStart(scope.row.cpi)">创建
                  </el-button>
                </div>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>
    <!-- 创建的Docker信息抽屉 --->
    <el-drawer
        title="服务创建成功"
        :visible.sync="drawer"
        :direction="direction"
        :before-close="handleClose">
        <p style="margin-left: 20px;">IP: {{ docker.ip }}</p>
        <p style="margin-left: 20px;">Port: {{ docker.port }}</p>
    </el-drawer>
  </div>
</template>

<script>
export default {
  name: "ServiceCreate",
  data() {
    return {
      waitCons: [],
      docker:{},
        direction: 'rtl',
      drawer: false
    }
  },
  created() {
    this.queryWaitCons();
  },
  methods: {
    async queryWaitCons() {
      const {data: res} = await this.$http({
        url: '/product/queryRaw',
        method: 'GET',
        headers: {
          "Authorization": 'Bearer ' + localStorage.getItem("userLogin_jwtToken")
        }
      })
      if (res.code === 1) {
        this.waitCons = res.data;
      }
    },
    statusFormatter(row, column, cellValue) {
      return this.$createElement('el-tag', {
        props: {type: 'danger', effect: 'plain'}
      }, '未激活');
    },
    handleStart(cpi) {
      this.$confirm('您确定创建此组件服务吗', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        const loading = this.$loading({
          lock: true,
          text: '创建过程耗时较长，请耐心等待，切勿退出',
          spinner: 'el-icon-loading',
          background: 'rgba(0, 0, 0, 0.7)'
        });
        this.createService(cpi).then(() => {
          loading.close();
        }).catch(() => {
          loading.close(); // 确保在任何情况下都关闭 loading
        });
      });
    },
    createService(cpi) {
      return this.$http({
        url: '/product/create',
        method: 'POST',
        headers: {
          "Authorization": 'Bearer ' + localStorage.getItem("userLogin_jwtToken")
        },
        data: {
          cpi: cpi
        }
      }).then(response => {
        const res = response.data;
        if (res.code !== 0) {
          this.docker=res.data;
          this.drawer = true;  // 打开Drawer
        } else {
          this.$message.error(res.msg);
        }
      }).catch(error => {
        console.error('请求失败:', error);
        this.$message.error('请求失败，请重试');
      });
    },
    // 处理Drawer的关闭逻辑
    handleClose(done) {
      this.$confirm('请保存服务访问地址')
          .then(_ => {
              // 关闭Drawer
            this.drawer = false;
            done();
            this.$router.push({ name:'ServiceList' });
          })
          .catch(_ => {});
    },
  }
}
</script>

<style scoped>

</style>
