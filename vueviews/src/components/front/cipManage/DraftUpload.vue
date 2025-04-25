<template>
  <div>
    <el-row id="breadcrumb">
      <el-breadcrumb separator-class="el-icon-arrow-right">
        <el-breadcrumb-item style="font-size: small">成果管理</el-breadcrumb-item>
        <el-breadcrumb-item>成果上传</el-breadcrumb-item>
      </el-breadcrumb>
    </el-row>
    <el-row id="main_range">
      <el-row>
        <el-col :span="15" offset="3">
          <div style="margin-top: 18px; margin-bottom: 50px;">
            <el-form
                ref="draftForm"
                :model="draftForm"
                :rules="rules"
                hide-required-asterisk="true"
                label-width="200px">
              <el-form-item label="名称" prop="filename">
                <el-input v-model="draftForm.filename" placeholder="请输入名称" clearable></el-input>
              </el-form-item>
              <el-form-item label="版本号" prop="version" id="elFormMargin">
                <el-input v-model="draftForm.version" placeholder="请输入版本号" clearable></el-input>
              </el-form-item>
              <el-form-item label="标签" prop="tag" id="elFormMargin">
                <el-input v-model="draftForm.tag" placeholder="请输入标签" clearable></el-input>
              </el-form-item>
<!--              <el-form-item label="程序设计语言" prop="language" id="elFormMargin">-->
<!--                <el-select v-model="draftForm.language" placeholder="请选择程序设计语言">-->
<!--                  <el-option label="Golang" value="Golang"></el-option>-->
<!--                  <el-option label="Java" value="Java"></el-option>-->
<!--                  <el-option label="Python" value="Python"></el-option>-->
<!--                </el-select>-->
<!--              </el-form-item>-->
              <el-form-item label="程序设计语言" prop="language" id="elFormMargin">
                <el-select
                  v-model="draftForm.language"
                  placeholder="请选择程序设计语言"
                  style="width: 100%">
                  <el-option label="Golang" value="Golang"></el-option>
                  <el-option label="Java" value="Java"></el-option>
                  <el-option label="Python" value="Python"></el-option>
                </el-select>
              </el-form-item>
              <el-form-item label="组件内部端口" prop="innerPort" id="elFormMargin">
                <el-input v-model="draftForm.innerPort" placeholder="请输入组件内部访问端口号" clearable></el-input>
              </el-form-item>
              <el-form-item label="组件描述" prop="description" id="elFormMargin">
                <el-input type="textarea" v-model="draftForm.description" :rows="3" placeholder="请描述组件针对的场景，实现的功能，以及具备的优势等，但字数不超过250字" clearable></el-input>
              </el-form-item>
              <el-form-item label="文件" prop="draftFile" id="elFormMargin">
                <el-upload
                    ref="upload"
                    action=""
                    class="upload-demo"
                    :http-request="uploadFile"
                    :on-exceed="handleExceed"
                    :multiple="false"
                    :limit="1"
                    :auto-upload="false"
                    :file-list="draftForm.draftFile"
                    drag>
                  <i class="el-icon-upload"></i>
                  <div class="el-upload__text">将文件拖到此处，或<em>点击上传</em></div>
                  <div class="el-upload__tip" slot="tip">只能上传.zip文件，且不超过500MB</div>
                </el-upload>
              </el-form-item>
              <el-form-item>
                <el-row type="flex">
                  <el-col :span="6" offset="6">
                    <el-button type="primary" @click="submitDraft('draftForm')" style="width: 100%">上传至草稿箱</el-button>
                  </el-col>
                </el-row>
              </el-form-item>
            </el-form>
          </div>
        </el-col>
      </el-row>
    </el-row>
  </div>
</template>

<script>
export default {
  name: "UploadDraft",
  data() {
    return {
      draftForm: {
        filename:'',
        version:'',
        tag:'',
        description:'',
        language:'',
        innerPort:'',
        draftFile:null
      },
      rules:{
        filename:[
          { required: true, message: '请输入名称', trigger: 'blur' },
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
        tag:[
          { required: true, message: '请输入标签', trigger: 'blur' }
        ],
        description: [
          { required: true, message: '请描述组件', trigger: 'blur' },
          { max: 250, message: '组件描述不能超过250字', trigger: 'blur' }
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
      loading: false
    }
  },
  methods: {
    // 自定义上传方法，使用上传组件的submit()后才会触发以获取文件实体
    uploadFile(param) {
      this.draftForm.draftFile = param.file;
    },
    handleExceed() {
      this.$alert("只允许上传一个文件", {  // 限制上传文件数量
        confirmButtonText: "确定",
      });
    },
    submitDraft(formName) {
      this.$refs[formName].validate((valid) => {
        if (valid) {
          this.$confirm('您确定要上传此组件吗', '提示', {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning'
          }).then(() => {
            const loading = this.$loading({
              lock: true,
              text: '正在上传，请耐心等待，请勿关闭',
              spinner: 'el-icon-loading',
              background: 'rgba(0, 0, 0, 0.7)'
            });
            // 手动调用上传，这里会调用我们自己定义的
            this.$refs.upload.submit()
            // 新建form表单
            const fd = new FormData()
            // 将draftForm转换为form-data，文件也被放进去
            Object.keys(this.draftForm).forEach(key => {
              if (this.draftForm[key] instanceof Array) {
                // 如果是数组就循环加入表单，key保持相同即可，这就是表达单的数组
                this.draftForm[key].forEach(item => {
                  fd.append(key, item)
                })
              } else {
                // 如果不是数组就直接追加进去
                fd.append(key, this.draftForm[key])
              }
            });
            this.saveDraft(fd).then(()=>{
              loading.close();
            });
          });
        } else {
          return false;
        }
      });
    },
    async saveDraft(data){
      const {data:res}=await this.$http({
        url:'/code/saveDraft',
        method:'POST',
        headers:{
          "Authorization":'Bearer '+localStorage.getItem("userLogin_jwtToken")
        },
        data:data
      })
      if(res.code===0){
        this.$message.error(res.msg);
      } else {
        this.$refs.draftForm.resetFields(); // 清空表单
        this.draftForm.draftFile = null; // 清空文件
        this.$refs.upload.clearFiles(); // 清空上传文件列表
        this.$notify({
          title: '成功',
          message: '您的组件已成功上传',
          type: 'success'
        });
      }
    }
  }
}
</script>

<style scoped>
#elFormMargin {
  /*margin-top: 25px;*/
}
:deep(.upload-demo .el-upload-dragger) {
    height: 150px;  /* 原高度为180px */
    padding: 20px 0;  /* 调整内边距 */
}

/* 调整图标和文字大小 */
:deep(.upload-demo .el-icon-upload) {
    font-size: 40px;
    margin: 10px 0;
}
:deep(.upload-demo .el-upload__text) {
    font-size: 12px;
}
</style>
