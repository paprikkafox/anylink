<template>
  <el-container style="height: 100%;">
    <el-aside :width="is_active?'200':'64'">
      <LayoutAside :is_active="is_active" :route_path="route_path"/>
    </el-aside>

    <el-container>
      <el-header>
        <LayoutHeader :is_active.sync="is_active" :route_name="route_name"/>
      </el-header>
      <el-main style="background-color: #fbfbfb">
        <router-view :route_path.sync="route_path" :route_name.sync="route_name"></router-view>
      </el-main>
    </el-container>
  </el-container>
</template>

<script>
import LayoutAside from "@/layout/LayoutAside";
import LayoutHeader from "@/layout/LayoutHeader";

export default {
  name: "Layout",
  components: {LayoutHeader, LayoutAside},
  data() {
    return {
      is_active: true,
      route_path: '/index',
      route_name: ['Front page'],
    }
  },
  methods: {
    goUrl(url) {
      window.open(url, "_blank")
    },
  },
  watch: {
    route_path: function (val) {
      window.console.log('is_active', val)
    },
  },
  created() {
    window.console.log('layout-route', this.$route)
  },
}
</script>

<style>
.el-header {
  background-color: #fff;
  color: #333;
  line-height: 60px;
  border-bottom: 1px solid #d8dce5;
  box-shadow: 0 1px 3px 0 rgba(0, 0, 0, .12), 0 0 3px 0 rgba(0, 0, 0, .04);
}

.el-footer {
  display: flex;
  align-items: center;
  justify-content: center;
  text-align: center;

  font-size: 12px;
  line-height: 12px;
  margin: 0 12px;
  color: rgb(134, 144, 156);
}

</style>
