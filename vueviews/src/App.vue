<template>
  <div id="app">
    <router-view></router-view>
  </div>
</template>

<script>
export default {
  name: 'App'
}
</script>

<!--高度铺满屏幕，其他组件搭配height:100%使用-->
<style>
html, body, #app{
  height: 100%;
  margin: 0;
  padding: 0;
}
</style>
