<template>

  <div class="login">
    <el-card style="width: 550px;">

      <div class="issuer">AnyLink SSL VPN Server</div>

      <el-form :model="ruleForm" status-icon :rules="rules" ref="ruleForm" label-width="100px" class="ruleForm">
        <el-form-item label="Username" prop="admin_user">
          <el-input v-model="ruleForm.admin_user"></el-input>
        </el-form-item>
        <el-form-item label="Password" prop="admin_pass">
          <el-input type="password" v-model="ruleForm.admin_pass" autocomplete="off"></el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="isLoading" @click="submitForm('ruleForm')">Log in</el-button>
          <el-button @click="resetForm('ruleForm')">Reset</el-button>
        </el-form-item>
      </el-form>

    </el-card>
  </div>

</template>

<script>
import axios from "axios";
import qs from "qs";
import {setToken, setUser} from "@/plugins/token";

export default {
  name: "Login",
  mounted() {
    console.log("Login created")
    window.addEventListener('keydown', this.keyDown);
  },
  destroyed(){
    window.removeEventListener('keydown',this.keyDown,false);
  },
  data() {
    return {
      ruleForm: {},
      rules: {
        admin_user: [
          {required: true, message: 'Please enter username', trigger: 'blur'},
          {max: 50, message: 'Username must be less than 50 characters long', trigger: 'blur'}
        ],
        admin_pass: [
          {required: true, message: 'Please enter password', trigger: 'blur'},
          {min: 6, message: 'Password must be longer than 6 characters', trigger: 'blur'}
        ],
      },
    }
  },
  methods: {
    keyDown(e) {
      if (e.keyCode === 13) {
        this.submitForm('ruleForm');
      }
    },
    submitForm(formName) {
      this.$refs[formName].validate((valid) => {
        if (!valid) {
          console.log('error submit!!');
          return false;
        }
        this.isLoading = true
        axios.post('/base/login', qs.stringify(this.ruleForm)).then(resp => {
          var rdata = resp.data
          if (rdata.code === 0) {
            this.$message.success(rdata.msg);
            setToken(rdata.data.token)
            setUser(rdata.data.admin_user)
            this.$router.push("/home");
          } else {
            this.$message.error(rdata.msg);
          }
          console.log(rdata);
        }).catch(error => {
          this.$message.error('Request error');
          console.log(error);
        }).finally(() => {
              this.isLoading = false
            }
        );

      });
    },
    resetForm(formName) {
      this.$refs[formName].resetFields();
    },
  },
}
</script>

<style scoped>
.login {
  height: 100%;
  text-align: center;

  display: flex;
  justify-content: center;
  align-items: center;
}

.issuer {
  font-size: 26px;
  font-weight: bold;
  margin-bottom: 50px;
}
</style>