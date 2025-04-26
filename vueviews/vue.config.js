const { defineConfig } = require('@vue/cli-service')
module.exports = defineConfig({
  devServer: {
    port: 8083
  },
  lintOnSave:false,
  transpileDependencies: true
})
