<template>
  <div>
    <el-card>
      <el-form :inline="true">
        <el-form-item>
          <el-button
              size="small"
              type="primary"
              icon="el-icon-plus"
              @click="handleEdit('')">Add
          </el-button>
        </el-form-item>
        <el-form-item>
          <el-dropdown size="small" placement="bottom">
            <el-upload
                class="uploaduser"
                action="uploaduser"
                accept=".xlsx, .xls"
                :http-request="upLoadUser"
                :limit="1"
                :show-file-list="false">
              <el-button size="small" icon="el-icon-upload2" type="primary">Batch add</el-button>
            </el-upload>
            <el-dropdown-menu slot="dropdown">
              <el-dropdown-item>
                <el-link style="font-size:12px;" type="success" href="user_templates.xlsx"><i
                    class="el-icon-download"></i>Download template
                </el-link>
              </el-dropdown-item>
            </el-dropdown-menu>
          </el-dropdown>
        </el-form-item>
        <el-form-item label="Username or email:">
          <el-input size="small" v-model="searchData" placeholder="Please enter"
                    @keydown.enter.native="searchEnterFun"></el-input>
        </el-form-item>

        <el-form-item>
          <el-button
              size="small"
              type="primary"
              icon="el-icon-search"
              @click="handleSearch()">Search
          </el-button>
          <el-button
              size="small"
              icon="el-icon-refresh"
              @click="reset">Reset
          </el-button>
        </el-form-item>
      </el-form>

      <el-table
          ref="multipleTable"
          :data="tableData"
          border>

        <el-table-column
            sortable="true"
            prop="id"
            label="ID"
            width="60">
        </el-table-column>

        <el-table-column
            prop="username"
            label="Username"
            width="150">
        </el-table-column>

        <el-table-column
            prop="nickname"
            label="Name"
            width="100">
        </el-table-column>

        <el-table-column
            prop="email"
            label="Email">
        </el-table-column>
        <el-table-column
            prop="otp_secret"
            label="OTP"
            width="110">
          <template slot-scope="scope">
            <el-button
                v-if="!scope.row.disable_otp"
                type="text"
                icon="el-icon-view"
                @click="getOtpImg(scope.row)">
              {{ scope.row.otp_secret.substring(0, 6) }}
            </el-button>
          </template>
        </el-table-column>

        <el-table-column
            prop="groups"
            label="Groups">
          <template slot-scope="scope">
            <el-row v-for="item in scope.row.groups" :key="item">{{ item }}</el-row>
          </template>
        </el-table-column>

        <el-table-column
            prop="status"
            label="Status"
            width="70">
          <template slot-scope="scope">
            <el-tag v-if="scope.row.status === 1" type="success">Available</el-tag>
            <el-tag v-if="scope.row.status === 0" type="danger">Inactive</el-tag>
            <el-tag v-if="scope.row.status === 2">Expired</el-tag>
          </template>
        </el-table-column>

        <el-table-column
            prop="updated_at"
            label="Updated at"
            :formatter="tableDateFormat">
        </el-table-column>

        <el-table-column
            label="Action"
            width="210">
          <template slot-scope="scope">
            <el-button
                size="mini"
                type="primary"
                @click="handleEdit(scope.row)">Edit
            </el-button>
            <el-popconfirm
                class="m-left-10"
                @confirm="handleDel(scope.row)"
                title="Are you sure you want to delete the user?">
              <el-button
                  slot="reference"
                  size="mini"
                  type="danger">Delete
              </el-button>
            </el-popconfirm>

          </template>
        </el-table-column>
      </el-table>

      <div class="sh-20"></div>

      <el-pagination
          background
          layout="prev, pager, next"
          :pager-count="11"
          @current-change="pageChange"
          :current-page="page"
          :total="count">
      </el-pagination>

    </el-card>

    <el-dialog
        title="OTP"
        :visible.sync="otpImgData.visible"
        width="350px"
        center>
      <div style="text-align: center">{{ otpImgData.username }} : {{ otpImgData.nickname }}</div>
      <img :src="otpImgData.base64Img" alt="otp-img"/>
    </el-dialog>

    <el-dialog
        :close-on-click-modal="false"
        title="User"
        :visible="user_edit_dialog"
        @close="disVisible"
        width="650px"
        center>

      <el-form :model="ruleForm" :rules="rules" ref="ruleForm" label-width="100px" class="ruleForm">
        <el-form-item label="ID" prop="id">
          <el-input v-model="ruleForm.id" disabled></el-input>
        </el-form-item>
        <el-form-item label="Username" prop="username">
          <el-input v-model="ruleForm.username" :disabled="ruleForm.id > 0"></el-input>
        </el-form-item>
        <el-form-item label="Name" prop="nickname">
          <el-input v-model="ruleForm.nickname"></el-input>
        </el-form-item>
        <el-form-item label="Email" prop="email">
          <el-input v-model="ruleForm.email"></el-input>
        </el-form-item>

        <el-form-item label="PIN" prop="pin_code">
          <el-input v-model="ruleForm.pin_code" placeholder="If left blank, the system will automatically generate it"></el-input>
        </el-form-item>

        <el-form-item label="Expiration" prop="limittime">
          <el-date-picker
              v-model="ruleForm.limittime"
              type="date"
              size="small"
              align="center"
              style="width:130px"
              :picker-options="pickerOptions"
              placeholder="Select date">
          </el-date-picker>
        </el-form-item>

        <el-form-item label="OTP" prop="disable_otp">
          <el-switch
              v-model="ruleForm.disable_otp"
              active-text="After turning on OTP, the user password is [PIN code + OTP dynamic code] (there is no + sign in the middle)">
          </el-switch>
        </el-form-item>

        <el-form-item label="OTP" prop="otp_secret" v-if="!ruleForm.disable_otp">
          <el-input v-model="ruleForm.otp_secret" placeholder="If left blank, the system will automatically generate it"></el-input>
        </el-form-item>

        <el-form-item label="Group" prop="groups">
          <el-checkbox-group v-model="ruleForm.groups">
            <el-checkbox v-for="(item) in grouNames" :key="item" :label="item" :name="item"></el-checkbox>
          </el-checkbox-group>
        </el-form-item>

        <el-form-item label="Send email" prop="send_email">
          <el-switch
              v-model="ruleForm.send_email">
          </el-switch>
        </el-form-item>

        <el-form-item label="Status" prop="status">
          <el-radio-group v-model="ruleForm.status">
            <el-radio :label="1" border>Enabled</el-radio>
            <el-radio :label="0" border>Disabled</el-radio>
            <el-radio :label="2" border>Expired</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="submitForm('ruleForm')">Save</el-button>
          <el-button @click="disVisible">Cancel</el-button>
        </el-form-item>
      </el-form>

    </el-dialog>

  </div>
</template>

<script>
import axios from "axios";

export default {
  name: "List",
  components: {},
  mixins: [],
  created() {
    this.$emit('update:route_path', this.$route.path)
    this.$emit('update:route_name', ['User info', 'User list'])
  },
  mounted() {
    this.getGroups();
    this.getData(1)
  },

  data() {
    return {
      page: 1,
      grouNames: [],
      tableData: [],
      count: 10,
      pickerOptions: {
        disabledDate(time) {
          return time.getTime() < Date.now();
        }
      },
      searchData: '',
      otpImgData: {visible: false, username: '', nickname: '', base64Img: ''},
      ruleForm: {
        send_email: true,
        status: 1,
        groups: [],
      },
      rules: {
        username: [
          {required: true, message: 'Please enter username', trigger: 'blur'},
          {max: 50, message: 'Username must be less than 50 characters long=', trigger: 'blur'}
        ],
        nickname: [
          {required: true, message: 'Please enter name', trigger: 'blur'}
        ],
        email: [
          {required: true, message: 'Please enter user email', trigger: 'blur'},
          {type: 'email', message: 'Please input the correct email address', trigger: ['blur', 'change']}
        ],
        password: [
          {min: 6, message: 'Password must be longer than 6 characters', trigger: 'blur'}
        ],
        pin_code: [
          {min: 6, message: 'PIN must be longer than 6 characters', trigger: 'blur'}
        ],
        date1: [
          {type: 'date', required: true, message: 'Please select a date', trigger: 'change'}
        ],
        groups: [
          {type: 'array', required: true, message: 'Please select at least one group', trigger: 'change'}
        ],
        status: [
          {required: true}
        ],
      },
    }
  },

  methods: {
    upLoadUser(item) {
      const formData = new FormData();
      formData.append("file", item.file);
      axios.post('/user/uploaduser', formData, {
        headers: {
          'Content-Type': 'multipart/form-data'
        }
      }).then(resp => {
        if (resp.data.code === 0) {
          this.$message.success(resp.data.data);
          this.getData(1);
        } else {
          this.$message.error(resp.data.msg);
          this.getData(1);
        }
        console.log(resp.data);
      })
    },
    getOtpImg(row) {
      this.otpImgData.visible = true
      axios.get('/user/otp_qr', {
        params: {
          id: row.id,
          b64: '1',
        }
      }).then(resp => {
        var rdata = resp.data;
        this.otpImgData.username = row.username;
        this.otpImgData.nickname = row.nickname;
        this.otpImgData.base64Img = 'data:image/png;base64,' + rdata
      }).catch(error => {
        this.$message.error('Request error');
        console.log(error);
      });
    },
    handleDel(row) {
      axios.post('/user/del?id=' + row.id).then(resp => {
        var rdata = resp.data
        if (rdata.code === 0) {
          this.$message.success(rdata.msg);
          this.getData(1);
        } else {
          this.$message.error(rdata.msg);
        }
        console.log(rdata);
      }).catch(error => {
        this.$message.error('Request error');
        console.log(error);
      });
    },
    handleEdit(row) {
      !this.$refs['ruleForm'] || this.$refs['ruleForm'].resetFields();
      console.log(row)
      this.user_edit_dialog = true
      if (!row) {
        return;
      }

      axios.get('/user/detail', {
        params: {
          id: row.id,
        }
      }).then(resp => {
        var data = resp.data.data
        data.send_email = false
        this.ruleForm = data
      }).catch(error => {
        this.$message.error('Request error');
        console.log(error);
      });
    },
    handleSearch() {
      console.log(this.searchData)
      this.getData(1, this.searchData)
    },
    pageChange(p) {
      this.getData(p)
    },
    getData(page, prefix) {
      this.page = page
      axios.get('/user/list', {
        params: {
          page: page,
          prefix: prefix || '',
        }
      }).then(resp => {
        var data = resp.data.data
        console.log(data);
        this.tableData = data.datas;
        this.count = data.count
      }).catch(error => {
        this.$message.error('Request error');
        console.log(error);
      });
    },
    getGroups() {
      axios.get('/group/names', {}).then(resp => {
        var data = resp.data.data
        console.log(data.datas);
        this.grouNames = data.datas;
      }).catch(error => {
        this.$message.error('Request error');
        console.log(error);
      });
    },
    submitForm(formName) {
      this.$refs[formName].validate((valid) => {
        if (!valid) {
          console.log('error submit!!');
          return false;
        }
        axios.post('/user/set', this.ruleForm).then(resp => {
          var data = resp.data
          if (data.code === 0) {
            this.$message.success(data.msg);
            this.getData(1);
            this.user_edit_dialog = false
          } else {
            this.$message.error(data.msg);
          }
          console.log(data);
        }).catch(error => {
          this.$message.error('Request error');
          console.log(error);
        });
      });
    },
    resetForm(formName) {
      this.$refs[formName].resetFields();
    },
    searchEnterFun(e) {
      var keyCode = window.event ? e.keyCode : e.which;
      if (keyCode == 13) {
        this.handleSearch()
      }
    },
    reset() {
      this.searchData = "";
      this.handleSearch();
    },
  },
}
</script>

<style scoped>

</style>
