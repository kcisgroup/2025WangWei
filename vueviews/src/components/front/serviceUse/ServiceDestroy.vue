<template>
    <div>
        <el-row id="breadcrumb">
            <el-breadcrumb separator-class="el-icon-arrow-right">
                <el-breadcrumb-item style="font-size: small">组件服务</el-breadcrumb-item>
                <el-breadcrumb-item>销毁服务</el-breadcrumb-item>
            </el-breadcrumb>
        </el-row>
        <el-row id="main_range">
            <el-col :span="22" :offset="1" style="margin-top: 20px;">
                <el-card class="box-card" shadow="hover">
                    <el-table
                            :data="services"
                            style="width: 100%">
                        <el-table-column
                                prop="dockerName"
                                label="服务名称"
                                width="280"
                                align="center">
                        </el-table-column>
                        <el-table-column
                                prop="cost"
                                label="费用"
                                align="center">
                        </el-table-column>
                        <el-table-column
                                prop="cpi"
                                label="产权数字摘要"
                                width="280"
                                align="center">
                        </el-table-column>
                        <el-table-column
                                prop="term"
                                label="使用时长（小时）"
                                align="center">
                        </el-table-column>
                        <el-table-column
                                label="操作"
                                align="center">
                            <template slot-scope="scope">
                                <div class="button-container">
                                    <el-button
                                            size="mini"
                                            type="danger"
                                            @click="handleR(scope.row.id)">销毁
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
    name: "ServiceDestroy",
    data() {
        return {
            services: []
        }
    },
    created() {
        this.queryDestroyList();
    },
    methods: {
        async queryDestroyList() {
            const {data: res} = await this.$http({
                url: '/product/queryDestroyList',
                method: 'GET',
                headers: {
                    "Authorization": 'Bearer ' + localStorage.getItem("userLogin_jwtToken")
                }
            })
            this.services = res.data;
        },
        async handleR(id) {
            this.$confirm('此操作将清空服务和数据，且不可恢复，是否继续?', '提示', {
                confirmButtonText: '确定',
                cancelButtonText: '取消',
                type: 'warning'
            }).then(async () => {
                const loading = this.$loading({
                    lock: true,
                    text: '正在清空，请耐心等待',
                    spinner: 'el-icon-loading',
                    background: 'rgba(0, 0, 0, 0.7)'
                });

            try {
            const { data: res } = await this.$http({
                url: '/product/destroy',
                method: 'DELETE',
                headers: {
                    "Authorization": 'Bearer ' + localStorage.getItem("userLogin_jwtToken")
                },
                data: {
                    id: id
                }
            });

            if (res.code === 0) {
                this.$message.error(res.msg);
            } else {
                this.$message.success(res.msg);
                await this.queryDestroyList();
            }
        } catch (error) {
                 this.$message.error('请求失败，请稍后重试');
                } finally {
                    loading.close();
                }
            });
        }
    }
}
</script>

<style scoped>
.button-container::-webkit-scrollbar {
    display: none;
}
</style>