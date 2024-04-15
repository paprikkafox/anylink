<template>

  <!--background-color="#304156"-->
  <!--text-color="#bfcbd9"-->
  <!--active-text-color="#409EFF"-->
  <!--:unique-opened="false"-->

  <!--<div class="layout-aside" :style="aside_style">-->

  <el-menu :collapse="!is_active"
           :default-active="route_path"
           :style="is_active?'width:200px':''"
           router
           class="layout-menu"
           :collapse-transition="false"

           background-color="#545c64"
           text-color="#fff"
           active-text-color="#ffd04b"
  >
    <el-menu-item index="/admin/home">
      <i class="el-icon-s-home"></i>
      <span slot="title">Front Page</span>
    </el-menu-item>

    <el-submenu index="1">
      <template slot="title">
        <i class="el-icon-menu"></i>
        <span slot="title">Basic info</span>
      </template>

      <el-menu-item index="/admin/set/system">System log</el-menu-item>
      <el-menu-item index="/admin/set/soft">Server config</el-menu-item>
      <el-menu-item index="/admin/set/other">Other settings</el-menu-item>
      <el-menu-item index="/admin/set/audit">Audit log</el-menu-item>
    </el-submenu>

    <el-submenu index="2">
      <template slot="title">
        <i class="el-icon-s-custom"></i>
        <span slot="title">Users</span>
      </template>

      <el-menu-item index="/admin/user/list">User list</el-menu-item>
      <el-menu-item index="/admin/user/policy">Policies</el-menu-item>
      <el-menu-item index="/admin/user/online">Online</el-menu-item>
      <el-menu-item index="/admin/user/ip_map">IP mapping</el-menu-item>
    </el-submenu>

    <el-submenu index="3">
      <template slot="title">
        <i class="el-icon-s-order"></i>
        <span slot="title">Groups</span>
      </template>

      <el-menu-item index="/admin/group/list">Groups list</el-menu-item>
    </el-submenu>

    <el-submenu index="4">
      <template slot="title">
        <i class="el-icon-s-order"></i>
        <span slot="title">Debug info</span>
      </template>

      <el-menu-item>
        <a href="/debug/pprof/" target="_blank">pprof</a>
      </el-menu-item>
      <el-menu-item>
        <a href="/debug/statsviz/" target="_blank">statsviz</a>
      </el-menu-item>

    </el-submenu>

  </el-menu>

</template>

<script>
export default {
  name: "LayoutAside",
  data() {
    return {}
  },
  props: ['is_active', 'route_path'],
  mounted() {
  }
}
</script>

<style scoped>
.layout-menu {
  height: 100%;
}

.el-menu-item a {
  display: block;
  color: #fff;
}
</style>
