<template>
  <div>
    <el-row id="breadcrumb">
      <el-breadcrumb separator-class="el-icon-arrow-right">
        <el-breadcrumb-item style="font-size: small">成果管理</el-breadcrumb-item>
        <el-breadcrumb-item :to="{ name: 'DraftBox' }">草稿箱</el-breadcrumb-item>
        <el-breadcrumb-item>编辑项目初稿</el-breadcrumb-item>
      </el-breadcrumb>
    </el-row>
    <el-row id="main_range">
      <el-row>
        <el-col :span="14" :offset="4" type="flex" justify="center">
          <el-form :model="draftInfo" :rules="editDraftRules" ref="draftInfo" style="margin-top: 38px;" label-width="110px">
            <el-form-item><span class="draftInfoSpan">项目信息更新</span></el-form-item>
            <el-form-item label="项目名" prop="filename">
              <el-input v-model="draftInfo.filename" placeholder="请输入项目名" clearable></el-input>
            </el-form-item>
            <el-form-item label="项目版本号" prop="version">
              <el-input v-model="draftInfo.version" placeholder="请输入版本号" clearable></el-input>
            </el-form-item>
            <el-form-item label="项目标签" prop="tag">
              <el-input v-model="draftInfo.tag" placeholder="请输入标签" clearable></el-input>
            </el-form-item>
            <el-form-item label="文件存储地址" prop="addr">
              <el-input v-model="draftInfo.addr" :disabled="true"></el-input>
            </el-form-item>
            <el-form-item label="项目上传时间" prop="createdTime">
              <el-input v-model="draftInfo.createdTime" :disabled="true"></el-input>
            </el-form-item>
            <el-form-item label="最近更新时间" prop="updatedTime">
              <el-input v-model="draftInfo.updatedTime" :disabled="true"></el-input>
            </el-form-item>
            <el-form-item label="内部访问端口" prop="innerPort">
              <el-input v-model="draftInfo.innerPort"></el-input>
            </el-form-item>
            <el-form-item label="程序设计语言" prop="language">
              <el-select v-model="draftInfo.language" placeholder="请选择程序设计语言">
                <el-option label="Golang" value="Golang"></el-option>
                <el-option label="Java" value="Java"></el-option>
                <el-option label="Python" value="Python"></el-option>
              </el-select>
            </el-form-item>
            <el-form-item label="项目描述" prop="description">
              <el-input type="textarea" v-model="draftInfo.description" :rows="4" placeholder="请描述项目针对的场景，能实现的功能，以及具备的优势等" clearable></el-input>
            </el-form-item>

            <el-form-item>
              <el-button type="primary" @click="subupdatedMetadata('draftInfo')" class="submitButton">提&nbsp&nbsp交</el-button>

              <el-button type="primary" @click="dialogFormVisible=true">项目文件更新</el-button>
              <el-dialog title="上传新的文件" :visible.sync="dialogFormVisible">
                <el-form ref="updateFileForm" :model="updateFileForm">
                  <el-form-item prop="draftFile">
                    <el-upload
                        ref="newFile"
                        action=""
                        class="upload-demo"
                        :http-request="getFile"
                        :on-exceed="handleExceedUpdate"
                        :multiple="false"
                        :limit="1"
                        :auto-upload="false"
                        :file-list="updateFileForm.draftFile"
                        drag
                        style="text-align: center">
                      <i class="el-icon-upload" style="display: block;"></i>
                      <div class="el-upload__text">将文件拖到此处，或<em>点击上传</em></div>
                      <div class="el-upload__tip" slot="tip">只能上传.zip文件，且不超过500MB</div>
                    </el-upload>
                  </el-form-item>
                  <el-form-item style="display: flex; justify-content: center; gap: 20px; margin-top: 20px;">
                    <el-button type="primary" @click="subUpdatedFile('updateFileForm')">提&nbsp交</el-button>
                    <el-button @click="dialogFormVisible = false">取&nbsp消</el-button>
                  </el-form-item>
                </el-form>
              </el-dialog>
            </el-form-item>
          </el-form>
        </el-col>
      </el-row>
    </el-row>
  </div>
</template>

<script>
export default {
  name: "DraftDetails",
  props: ['id'],
  data(){
    return{
      draftInfo:{},
      editDraftRules:{
        filename:[
          { required: true, message: '请输入项目名', trigger: 'blur' },
          {
            validator: (rule, value, callback) => {
              const regExp = /^[a-zA-Z]{3,20}$/; // 使用正则表达式匹配英文字符，字符数3到10
              if (regExp.test(value)) {
                callback(); // 符合规则，验证通过
              } else {
                callback(new Error('请输入3到10个英文字符，不支持其他字符，如 _ * -')); // 不符合规则，显示错误消息
              }
            },
            trigger: 'blur'
          }
        ],
        version: [
          { required: true, message: '请输入版本号', trigger: 'blur' },
          {
            validator: (rule, value, callback) => {
              const regExp = /^\d+(\.\d+)*$/; // 使用正则表达式匹配版本号格式
              if (regExp.test(value)) {
                callback(); // 符合规则，验证通过
              } else {
                callback(new Error('请输入正确格式的版本号，例如：1.0 或 1.0.0')); // 不符合规则，显示错误消息
              }
            },
            trigger: 'blur'
          }
        ],
        tag: [
          { required: true, message: '请输入标签', trigger: 'blur' }
        ],
        description: [
          { required: true, message: '请输入功能', trigger: 'blur' }
        ],
        language: [
          { required: true, message: '请选择程序设计语言', trigger: 'blur' }
        ],
        innerPort: [
          { required: true, message: '请输入组件内部访问端口号', trigger: 'blur' },
          {
            pattern: /^\d{1,5}$/,
            message: '端口号必须是1-5位数字',
            trigger: 'blur'
          }
        ]
      },
      dialogFormVisible: false,
      updateFileForm: {
        draftFile:null
      },
    }
  },
  created(){
    this.draftDetails();
  },
  methods:{
    async draftDetails(){
      const {data:res}=await this.$http({
        url:'/code/queryDetails',
        method:'POST',
        headers:{
          "Authorization":'Bearer '+localStorage.getItem("userLogin_jwtToken")
        },
        data:{
          id:this.id
        }
      })
      if(res.code===0){
        // this.$message.error('您的草稿设计信息在赶来的路上走丢了，请重新操作');
        this.$message.error(res.msg);
      } else {
        // 将原信息回显至表单中
        this.draftInfo=res.data;
      }
    },
    // 提交表单修改之后的项目元数据信息
    subupdatedMetadata(formName) {
      this.$refs[formName].validate((valid) => {
        if (valid) {
          this.$confirm('您确定提交新的项目信息吗', '提示', {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning'
          }).then(() => {
            const loading = this.$loading({
              lock: true,
              text: '正在上传服务器，请耐心等待，请勿关闭',
              spinner: 'el-icon-loading',
              background: 'rgba(0, 0, 0, 0.7)'
            });
            // 发送更新请求
            this.newDraftMetadata().then(()=>{
              loading.close();
            });
          });
        } else {
          return false;
        }
      });
    },
    async newDraftMetadata(){
      const {data:res}=await this.$http({
        url:'/code/updateDraftMetaData',
        method:'PUT',
        headers:{
          "Authorization":'Bearer '+localStorage.getItem("userLogin_jwtToken")
        },
        data:{
          id:this.draftInfo.id,
          filename:this.draftInfo.filename,
          version:this.draftInfo.version,
          language:this.draftInfo.language,
          innerPort:this.draftInfo.innerPort,
          description:this.draftInfo.description
        }
      })
      if(res.code===0){
        this.$message.error(res.msg);
      } else {
        await this.$router.push({name: 'DraftBox'});
        this.$notify({
          title: '成功',
          message: '新设计上传成功',
          type: 'success'
        });
      }
    },
    // 自定义上传方法，使用上传组件的submit()后才会触发以获取文件实体
    getFile(param) {
      this.updateFileForm.draftFile = param.file;
    },
    handleExceedUpdate() {
      this.$alert("只允许上传一个文件", {  // 限制上传文件数量
        confirmButtonText: "确定",
      });
    },
    async updateDraftFile(data){ // 发送请求，提交新项目文件
      const {data:res}=await this.$http({
        url:'/code/updateDraftFile',
        method:'POST',
        headers:{
          "Authorization":'Bearer '+localStorage.getItem("userLogin_jwtToken")
        },
        data:data
      })
      if(res.code===0){
        this.$message.error(res.msg);
      } else {
        await this.$router.push({name:'DraftBox'})
        this.$notify({
          title: '成功',
          message: '新项目文件上传成功',
          type: 'success'
        });
      }
    },
    subUpdatedFile(formName){ // 上传新项目文件
      this.$refs[formName].validate((valid) => {
        if (valid) {
          this.$confirm('您确定提交新的文件吗？这将覆盖原有文件', '提示', {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning'
          }).then(() => {
            const loading = this.$loading({
              lock: true,
              text: '正在上传至服务器，请耐心等待，请勿关闭',
              spinner: 'el-icon-loading',
              background: 'rgba(0, 0, 0, 0.7)'
            });
            // 手动调用上传，这里会调用我们自己定义的
            this.$refs.newFile.submit()
            // 新建form表单
            const fd = new FormData()
            fd.append("id", this.draftInfo.id)
            fd.append("draftFile", this.updateFileForm.draftFile)
            this.updateDraftFile(fd).then(()=>{
              loading.close();
            });
          });
        } else {
          return false;
        }
      });
    }
  }
}
</script>

<style scoped>
.draftInfoSpan{
  margin-top: 20px;
  margin-left: 200px;
  font-size: 25px;
  font-family: "微软雅黑";
}
.submitButton{
  margin-top: 20px;
  margin-left: 143px;
  width: 120px;
}
</style>
