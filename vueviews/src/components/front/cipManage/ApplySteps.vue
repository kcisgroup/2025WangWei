<template>
  <div>
    <el-row id="breadcrumb">
      <el-breadcrumb separator-class="el-icon-arrow-right">
        <el-breadcrumb-item style="font-size: small">成果管理</el-breadcrumb-item>
        <el-breadcrumb-item>申请代码产权</el-breadcrumb-item>
      </el-breadcrumb>
    </el-row>
    <el-row id="main_range">
      <br>
      <!--进度条-->
      <el-row>
        <el-col :span="22" :offset="1" class="steps">
          <el-steps :active="active" finish-status="success" simple>
            <el-step title="选择项目" icon="el-icon-s-promotion"></el-step>
            <el-step title="完善信息" icon="el-icon-edit"></el-step>
            <el-step title="确认申请信息" icon="el-icon-upload"></el-step>
            <el-step title="查看产权登记结果" icon="el-icon-document-checked"></el-step>
          </el-steps>
        </el-col>
      </el-row>
      <!--各阶段-->
      <el-row style="margin-bottom: 50px;">
        <br>
        <!--选择项目-->
        <el-col :span="22" :offset="1" v-if="active===0">
          <el-table
              ref="singleTable"
              :data="drafts"
              highlight-current-row
              border
              @current-change="handleCurrentChange">
            <el-table-column
                prop="filename"
                label="项目名"
                align="center"
                width="150">
            </el-table-column>
            <el-table-column
                prop="version"
                label="版本号"
                align="center"
                width="100">
            </el-table-column>
            <el-table-column
                prop="tag"
                label="标签"
                align="center"
                width="100">
              <template slot-scope="scope">
                <el-tag
                    type="warning"
                    disable-transitions
                    effect="plain">{{scope.row.tag}}</el-tag>
              </template>
            </el-table-column>
            <el-table-column
                prop="language"
                label="程序设计语言"
                align="center"
                width="130">
              <template slot-scope="filteredDrafts">
                <el-tag v-if="filteredDrafts.row.language=='Golang'">Golang</el-tag>
                <el-tag v-if="filteredDrafts.row.language=='Java'">Java</el-tag>
                <el-tag v-if="filteredDrafts.row.language=='Python'">Python</el-tag>
              </template>
            </el-table-column>
            <el-table-column
                prop="description"
                label="项目描述"
                align="center"
                width="235">
            </el-table-column>
            <el-table-column
                prop="addr"
                label="项目文件地址"
                align="center"
                width="238">
            </el-table-column>
            <el-table-column
                prop="createdTime"
                label="上传时间"
                align="center"
                width="180">
            </el-table-column>
            <el-table-column
                prop="updatedTime"
                label="最近更新时间"
                align="center"
                width="180">
            </el-table-column>
          </el-table>
          <el-button type="text" disabled>单击表格中的项目，即可完成选择</el-button>
        </el-col>
        <!--完善信息-->
        <el-col :span="14" :offset="4" v-if="active===1">
          <br>
          <el-form :model="newDraftInfo" :rules="rules" ref="newDraftInfo" label-width="119px">
            <el-form-item label="项目名" prop="filename" class="subFormItem">
              <el-input v-model="projectDetails.filename" :disabled="true"></el-input>
            </el-form-item>
            <el-form-item label="版本号" prop="version" class="subFormItem">
              <el-input v-model="projectDetails.version" :disabled="true"></el-input>
            </el-form-item>
            <el-form-item label="标签" prop="tag" class="subFormItem">
              <el-input v-model="projectDetails.tag" :disabled="true"></el-input>
            </el-form-item>
            <el-form-item label="程序设计语言" prop="language" class="subFormItem">
              <el-input v-model="projectDetails.language" :disabled="true"></el-input>
            </el-form-item>
            <el-form-item label="内部访问端口" prop="innerPort" class="subFormItem">
              <el-input v-model="projectDetails.innerPort" :disabled="true"></el-input>
            </el-form-item>
            <el-form-item label="项目描述" prop="description" class="subFormItem">
              <el-input type="textarea" v-model="projectDetails.description" :rows="3" placeholder="请描述项目针对的场景，能实现的功能，以及具备的优势等" :disabled="true"></el-input>
            </el-form-item>
            <el-form-item label="是否开放订阅" prop="status" class="subFormItem">
              <el-select
                v-model="newDraftInfo.status"
                placeholder="请选择是否开放"
                style="display: block; width: 100%">
                <el-option label="禁用" value="禁用"></el-option>
                <el-option label="开放" value="开放"></el-option>
              </el-select>
            </el-form-item>
            <el-form-item label="预设订阅价格" prop="price" class="subFormItem">
              <el-input v-model="newDraftInfo.price" placeholder="请预设订阅价格，若禁用项目，订阅价格请填写0，反之请填写实际定价">
                <template slot="append">元/天</template>
              </el-input>
            </el-form-item>
            <br>
            <el-row>
              <el-col :offset="9">
                <el-button type="info" style="width: 100px;" @click="backToSelect" plain>上一步</el-button>
                <el-button type="primary" style="width: 100px;" @click="handleEdit('newDraftInfo')">下一步</el-button>
              </el-col>
            </el-row>
          </el-form>
        </el-col>
        <!--确认申请信息-->
        <el-col :span="18" :offset="3" v-if="active===2">
          <br>
          <el-row>
            <el-descriptions title="代码知识产权申请表" direction="vertical" border>
              <el-descriptions-item label="项目名">{{this.projectDetails.filename}}</el-descriptions-item>
              <el-descriptions-item label="项目版本号">{{this.projectDetails.version}}</el-descriptions-item>
              <el-descriptions-item label="标签"><el-tag size="small">{{this.projectDetails.tag}}</el-tag></el-descriptions-item>
              <el-descriptions-item label="申请人用户名">{{this.projectDetails.author}}</el-descriptions-item>
              <el-descriptions-item label="程序设计语言">{{this.projectDetails.language}}</el-descriptions-item>
              <el-descriptions-item label="是否开放交易">{{this.newDraftInfo.status}}</el-descriptions-item>
              <el-descriptions-item label="订阅价格（元/天）">{{this.newDraftInfo.price}}</el-descriptions-item>
              <el-descriptions-item label="项目上传时间">{{this.projectDetails.createdTime}}</el-descriptions-item>
              <el-descriptions-item label="最近修改时间">{{this.projectDetails.updatedTime}}</el-descriptions-item>
              <el-descriptions-item label="项目文件存储地址" :span="3">{{this.projectDetails.addr}}</el-descriptions-item>
              <el-descriptions-item label="项目描述" :span="3">{{this.projectDetails.description}}</el-descriptions-item>
            </el-descriptions>
          </el-row>
          <br><br>
          <el-row>
            <el-col :span="8" :offset="9">
              <el-button type="info" @click="backToEditInfo" plain>上一步</el-button>
              <el-button type="primary" @click="handleApply">提&nbsp&nbsp交</el-button>
            </el-col>
          </el-row>
        </el-col>
        <!--查看产权登记结果-->
        <el-col v-if="active===3">
          <!--申请失败-->
          <div v-if="this.isFail===true">
            <el-col :span="10" :offset="1"> <!--知识产权所有人信息-->
              <el-row>
                <el-card shadow="hover" class="box-card">
                  <div slot="header" class="clearfix">
                    <span>申请人信息</span>
                  </div>
                  <el-descriptions :column="2" border>
<!--                    <el-descriptions-item label="姓名">{{this.result.name}}</el-descriptions-item>-->
                    <el-descriptions-item label="性别">{{this.result.gender}}</el-descriptions-item>
                    <el-descriptions-item label="联系电话">{{this.result.tel}}</el-descriptions-item>
                    <el-descriptions-item label="平台用户名">{{this.result.username}}</el-descriptions-item>
                    <el-descriptions-item label="就职公司（组织）">{{this.result.role}}</el-descriptions-item>
                    <el-descriptions-item label="地址">{{this.result.addr}}</el-descriptions-item>
                    <el-descriptions-item label="身份证号" :span="2">{{this.result.idnum}}</el-descriptions-item>
                  </el-descriptions>
                </el-card>
              </el-row>
              <br><br>
              <el-row>
                <el-card shadow="hover" style="background-color: #f6f8fc;" class="box-card">
                  <el-result icon="error" title="申请失败" :subTitle="this.cause">
                    <template slot="extra">
                      <el-button type="primary" size="medium" @click="reApply">再次申请</el-button>
                      <el-button type="primary" size="medium" @click="goToDraftBox">前往草稿箱</el-button>
                    </template>
                  </el-result>
                </el-card>
              </el-row>
            </el-col>
            <el-col :span="11" :offset="1"> <!--申请的项目信息-->
              <el-card shadow="hover" class="box-card">
                <div slot="header" class="clearfix">
                  <span>申请产权项目信息</span>
                </div>
                <el-descriptions class="margin-top" direction="vertical" :column="2" border>
                  <el-descriptions-item label="名称">{{this.projectDetails.filename}}</el-descriptions-item>
                  <el-descriptions-item label="版本号">{{this.projectDetails.version}}</el-descriptions-item>
                  <el-descriptions-item label="项目标签"><el-tag effect="light" type="warning">{{this.projectDetails.tag}}</el-tag></el-descriptions-item>
                  <el-descriptions-item label="程序设计语言">{{this.projectDetails.language}}</el-descriptions-item>
                  <el-descriptions-item label="初次上传时间">{{this.projectDetails.createdTime}}</el-descriptions-item>
                  <el-descriptions-item label="近期更新时间">{{this.projectDetails.updatedTime}}</el-descriptions-item>
                  <el-descriptions-item label="项目文件存储地址" :span="2"><el-tag effect="light">{{this.projectDetails.addr}}</el-tag></el-descriptions-item>
                  <el-descriptions-item label="项目描述" :span="2">{{this.projectDetails.description}}</el-descriptions-item>
                </el-descriptions>
              </el-card>
            </el-col>
          </div>
          <!--申请成功-->
          <div v-else-if="this.isFail===false">
            <el-col :span="10" :offset="1">
              <el-row :span="10" :offset="1"> <!--知识产权基本信息-->
                <el-card shadow="hover" class="box-card">
                  <div slot="header" class="clearfix">
                    <span>知识产权基本信息</span>
                  </div>
                  <el-descriptions :column="4" direction="vertical" border>
                    <el-descriptions-item label="名称">{{this.result.filename}}</el-descriptions-item>
                    <el-descriptions-item label="版本号">{{this.result.version}}</el-descriptions-item>
                    <el-descriptions-item label="是否开放订阅">{{this.result.status}}</el-descriptions-item>
                    <el-descriptions-item v-if="this.result.status==='开放'" label="初拟订阅价格" :span="4">{{this.result.price}}</el-descriptions-item>
                    <el-descriptions-item v-if="this.result.status!=='开放'" label="初拟订阅价格" :span="4">—</el-descriptions-item>
                    <el-descriptions-item label="代码知识产权唯一标识码" :span="4"><el-tag effect="plain" type="danger">{{this.result.cpi}}</el-tag></el-descriptions-item>
                    <el-descriptions-item label="确权日期时间">{{this.result.regTime}}</el-descriptions-item>
                  </el-descriptions>
                </el-card>
              </el-row>
              <br>
              <el-row> <!--知识产权所有人信息-->
                <el-card shadow="hover" class="box-card">
                  <div slot="header" class="clearfix">
                    <span>知识产权所有人信息</span>
                  </div>
                  <el-descriptions :column="4" direction="vertical" border>
                    <el-descriptions-item label="姓名">{{this.result.name}}</el-descriptions-item>
                    <el-descriptions-item label="性别">{{this.result.gender}}</el-descriptions-item>
                    <el-descriptions-item label="联系电话">{{this.result.tel}}</el-descriptions-item>
                    <el-descriptions-item label="平台用户名">{{this.result.username}}</el-descriptions-item>
                    <el-descriptions-item label="身份证号" :span="2">{{this.result.idnum}}</el-descriptions-item>
                    <el-descriptions-item label="就职公司（组织）">{{this.result.role}}</el-descriptions-item>
                    <el-descriptions-item label="地址">{{this.result.authorAddr}}</el-descriptions-item>
                  </el-descriptions>
                </el-card>
              </el-row>
            </el-col>

            <el-col :span="11" :offset="1">
              <el-row> <!--项目信息-->
                <el-card shadow="hover" class="box-card">
                  <div slot="header" class="clearfix">
                    <span>项目信息</span>
                  </div>
                  <el-descriptions :column="4" size="medium" direction="vertical" border>
                    <el-descriptions-item label="项目标签"><el-tag effect="light" type="warning">{{this.result.tag}}</el-tag></el-descriptions-item>
                    <el-descriptions-item label="程序设计语言">{{this.result.language}}</el-descriptions-item>
                    <el-descriptions-item label="初次上传时间">{{this.result.createdTime}}</el-descriptions-item>
                    <el-descriptions-item label="近期更新时间">{{this.result.updatedTime}}</el-descriptions-item>
                    <el-descriptions-item label="项目文件存储地址" :span="4"><el-tag effect="light">{{this.result.fileAddr}}</el-tag></el-descriptions-item>
                    <el-descriptions-item label="项目描述" :span="4">{{this.result.description}}</el-descriptions-item>
                  </el-descriptions>
                </el-card>
              </el-row>
              <br><br><br>
              <el-row> <!--成功确权-->
                <div>
                  <el-result icon="success" title="成功确权" style="height: 193px;">
                    <template slot="extra">
                      <el-button type="primary" size="medium" @click="acceptResult">确认收到</el-button>
                    </template>
                  </el-result>
                </div>
              </el-row>
            </el-col>
          </div>
        </el-col>
      </el-row>
    </el-row>
  </div>
</template>

<script>
export default {
  name: "ApplySteps",
  data(){
    return{
      active: 0,
      drafts:[],
      projectDetails:{},
      newDraftInfo:{
        price: null,
        status: null
      },
      rules:{
        status: [
          { required: true, message: '请选择是否开放', trigger: 'blur' }
        ],
        price: [
          { required: true, message: '请预设订阅价格，若禁用，请填写0', trigger: 'blur' },
          {
            validator: (rule, value, callback) => {
              if (!/^\d+$/.test(value)) {
                callback(new Error('请输入数字，不能包含小数'));
              } else {
                callback();
              }
            },
            trigger: 'blur'
          }
        ]
      },
      result: {},
      cause: null,
      isFail: false
    }
  },
  created() {
    this.queryDraftBox();
  },
  methods:{
    // 下一步
    nextStep(){
      this.active++;
    },

    // 上一步
    backStep(){
      this.active--;
    },

    // 查看草稿箱
    async queryDraftBox(){
      const {data:res}=await this.$http({
        url:'/code/draftBox',
        method:'GET',
        headers:{
          "Authorization":'Bearer '+localStorage.getItem("userLogin_jwtToken")
        }
      })
      if(res.code===1){
        this.drafts=res.data;
      }
    },

    // 选择项目确认按钮
    handleCurrentChange(val) {
      this.projectDetails = val;
      this.$confirm('您确定选择此项目吗', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '重新选择',
        type: 'warning'
      }).then(() => {
        this.nextStep();
        this.projectDetails.status=null;
        this.projectDetails.price=null;
      });
    },

    // 返回选择项目步骤
    backToSelect(){
      this.backStep();
      this.projectDetails=null;
      this.queryDraftBox();
    },

    // 处理完善信息步骤提交按钮
    handleEdit(formName){
      this.$refs[formName].validate((valid) => {
        if (valid) {
          this.nextStep();
        } else {
          return false;
        }
      });
    },

    // 返回完善信息步骤
    backToEditInfo(){
      this.backStep();
    },

    // 处理申请
    handleApply(){
      const loading = this.$loading({
        lock: true,
        text: '正在处理您的申请，请耐心等待，请勿关闭',
        spinner: 'el-icon-loading',
        background: 'rgba(0, 0, 0, 0.7)'
      });
      this.applyForCip().then(() => {
        loading.close();
      });
    },

    // 发送代码知识产权申请请求
    async applyForCip() {
      const {data: res} = await this.$http({
        url: '/code/register',
        method: 'POST',
        headers: {
          "Authorization": 'Bearer ' + localStorage.getItem("userLogin_jwtToken")
        },
        data: {
          id: this.projectDetails.id,
          price: parseInt(this.newDraftInfo.price, 10),
          status: this.newDraftInfo.status
        },
      });
      if (res.code===0) { // 申请失败
        this.result = res.data;
        this.cause = "驳回原因："+res.msg;
        this.isFail=true;
        if (this.projectDetails.description.length > 250) {
          this.projectDetails.description = this.projectDetails.description.substring(0, 250) + '······（略）';
        }
        this.nextStep();
      } else { // 申请成功
        this.result = res.data;
        this.cause = null;
        if (this.result.description.length > 198) {
          this.result.description = this.result.description.substring(0, 198) + '······（略）';
        }
        this.nextStep();
      };
      console.log("active===========",this.active);
      console.log("cause===========",this.cause);
    },
    acceptResult(){ // 确认申请成功
      this.$router.push({ name:'MyCips' });
    },
    reApply(){ // 重新申请
      this.$router.push({ name:'ApplyForCip' });
    },
    goToDraftBox(){ // 前往草稿箱
      this.$router.push({ name:'DraftBox' });
    }
  }
}
</script>

<style scoped>
.subFormItem{
  margin-bottom: 30px;
}
</style>