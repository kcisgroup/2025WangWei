<template>
  <div>
    <el-row id="breadcrumb">
      <el-breadcrumb separator-class="el-icon-arrow-right">
        <el-breadcrumb-item style="font-size: small">成果管理</el-breadcrumb-item>
        <el-breadcrumb-item>申请代码产权</el-breadcrumb-item>
      </el-breadcrumb>
    </el-row>
    <el-row id="main_range">
      <br><br>
      <el-row>
        <el-col :span="20" :offset="2">
          <el-card shadow="hover" class="box-card" style="background-color: #f4f7fc;">
            <div slot="header" class="clearfix">
              <el-row>
                <el-col :span="5" :offset="9">
                  <span>代码知识产权申请须知</span>
                </el-col>
              </el-row>
            </div>
            <div class="text item">
              <div id="notice">1. 申请人必须自己上传本人项目，并办理代码知识产权登记。</div>
              <div id="notice">2. 申请人应当将所提交的申请文件留存一份，便于在查验代码知识产权时校验文件内容完整性。</div>
              <div id="notice">3. 办理代码知识产权登记需登录代码产权可信互用平台，全程在线办理。</div>
                <div id="notice">① 代码知识产权登记咨询电话</div>
                  <div id="notice">代码知识产权登记业务咨询：xxx-xxxxxxxx</div>
                  <div id="notice">代码知识产权信息更新业务咨询：xxx-xxxxxxxx</div>
                  <div id="notice">代码知识产权申诉：xxx-xxxxxxxx</div>
                  <div id="notice">其他疑问：xxx-xxxxxxxx</div>
                <div id="notice">② 代码知识产权登记投诉及建议</div>
                  <div id="notice">cipti@xxxxxxxxx.com</div>
              <div id="notice">4. 申请代码知识产权前，必须确保代码项目已上传至本平台，在提交申请后，不能对代码项目进行修改（订阅价格和产权状态除外）。</div>
              <div id="notice">5. 申请代码知识产权前，确保代码项目必须满足以下要求：</div>
                <div id="notice">（1）保证"代码组件-版本"唯一，如"codeTest 1.0.0"具有唯一性。</div>
                <div id="notice">（2）上传的代码项目文件后缀名是.zip，且在项目文件解压后的根目录下，包含项目启动运行的DockerFile文件。</div>
                <div id="notice">（2）必须保证代码组件端口设置正确。</div>
                <div id="notice">（3）项目描述清晰明确，可阐述代码项目面向的领域、具备的功能等，杜绝抄袭、弄虚作假和扭曲事实等违法违规行为。</div>
                <div id="notice">（4）API说明文档内容应当详细全面，确保内容完整、可准确识别，方便使用者查看和使用。</div>
              <div id="notice">6. 本平台技术支持对代码项目全过程保护，不会泄露任何代码知识产权隐私，同时也不制造任何有损代码隐私性的事件。</div>
              <div id="notice">7. 在代码知识产权确立后，如果组件开放订阅，那么可以被平台调用和生成服务，经创作者同意之后可以被其他用户订阅和使用。
                如果将"是否开放订阅"设置为"禁用"，那么平台保证知识产权不参与使用权交易。"是否开放订阅"的设置可随时更改的，但已订阅的用户仍可继续使用。
                  注意：所有在本平台订阅和使用的软件服务均不可查看源代码，即开放订阅≠开源</div>
            </div>
          </el-card>
        </el-col>
      </el-row>
      <br><br>
      <el-row type="flex">
        <el-col :span="4" :offset="10">
          <el-button @click="firstStep" :type="setButtonType" :disabled="isDisabled" :icon="setIcon" style="margin-bottom: 50px;" plain>{{ buttonText }}</el-button>
        </el-col>
      </el-row>
    </el-row>
  </div>
</template>

<script>
export default {
  name: "ApplyForCip",
  data() {
    return {
      countdownTime: 5, // 设置初始倒计时时间（单位为秒）
      isDisabled: true, // 控制按钮状态的变量
      timerId: null, // 定时器ID
    };
  },
  created() {
    this.startCountdown();
  },
  computed: {
    setButtonType(){
      if (this.countdownTime > 0) {
        return "info";
      } else {
        return "success";
      }
    },
    setIcon(){
      if (this.countdownTime <= 0) {
        return "el-icon-circle-check";
      }
    },
    buttonText() {
      if (this.countdownTime > 0) {
        return `我已知晓上述内容(${this.countdownTime}s)`;
      } else {
        return "我已知晓上述内容";
      }
    }
  },
  methods:{
    firstStep(){
      this.$router.push({ name: 'ApplySteps' });
    },
    startCountdown() {
      this.timerId = setInterval(() => {
        this.countdownTime--;
        if (this.countdownTime <= 0) {
          clearInterval(this.timerId); // 在组件销毁前清除定时器，以防止内存泄漏
          this.isDisabled=false;
        }
      }, 1000);
    },
  },
}
</script>

<style scoped>
#notice{
  margin-bottom: 4px;
}
</style>