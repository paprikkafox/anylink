<template>
  <div>
<el-form :model="searchForm" :rules="rules" ref="searchForm" :inline="true" class="search-form">
        <el-form-item label="Username:" prop="username">
          <el-input size="mini" v-model="searchForm.username" clearable style="width: 130px" @keydown.enter.native="searchEnterFun"></el-input>
        </el-form-item>
        <el-form-item label="Source IP address:" prop="src">
          <el-input size="mini" v-model="searchForm.src" clearable style="width: 130px" @keydown.enter.native="searchEnterFun"></el-input>
        </el-form-item>    
        <el-form-item label="Destination IP address:" prop="dst">
          <el-input size="mini" v-model="searchForm.dst" clearable style="width: 130px" @keydown.enter.native="searchEnterFun"></el-input>
        </el-form-item> 
        <el-form-item label="Destination port:" prop="dst_port">
          <el-input size="mini" v-model="searchForm.dst_port" clearable style="width: 80px" @keydown.enter.native="searchEnterFun"></el-input>
        </el-form-item> 
        <el-form-item label="Access protocol：">
            <el-select size="mini" v-model="searchForm.access_proto" clearable placeholder="Please choose" style="width: 100px">
                    <el-option v-for="(item,index) in access_proto" :key="index" :label="item.text" :value="item.value">
                    </el-option>
            </el-select>           
        </el-form-item>  
        <el-form-item label="Data range：">
            <el-date-picker
                v-model="searchForm.date"
                type="datetimerange"
                value-format="yyyy-MM-dd HH:mm:ss"
                size="mini"
                align="left"
                start-placeholder="Start date"
                end-placeholder="End date"
                :default-time="['00:00:00', '23:59:59']">
            </el-date-picker>
        </el-form-item>
        <el-form-item label="Details:">
          <el-input size="mini" v-model="searchForm.info" placeholder="Please enter key words" clearable style="width: 200px" @keydown.enter.native="searchEnterFun"></el-input>
        </el-form-item>         
        <el-form-item>
          <el-button
              size="mini"
              type="primary"
              icon="el-icon-search"
              @click="handleSearch">Search
          </el-button>
          <el-button
              size="mini"
              icon="el-icon-refresh"
              @click="rest">Reset
          </el-button>
          <el-button
              size="mini"
              icon="el-icon-download"
              @click="handleExport">Export
          </el-button>          
        </el-form-item>
      </el-form>

      <el-table
          ref="multipleTable"
          :data="tableData"
          v-loading="loading"
          element-loading-text="Loading"
          element-loading-spinner="el-icon-loading"
          :default-sort="{ prop: 'id', order: 'descending' }"
          @sort-change="sortChange" 
          :header-cell-style="{backgroundColor:'#fcfcfc'}"
          border>

        <el-table-column
            prop="id"
            label="ID"
            sortable="custom"
            width="100">
        </el-table-column>

        <el-table-column
            prop="username"
            label="Username"
            width="140">
        </el-table-column>

        <el-table-column
            prop="src"
            label="Source IP address"
            width="140">
        </el-table-column>

        <el-table-column
            prop="dst"
            label="Destination IP address"
            width="140">
        </el-table-column>

        <el-table-column
            prop="dst_port"
            label="Destination port"
            width="85">
        </el-table-column>

        <el-table-column
            prop="access_proto"
            label="Access protocol"
            width="80"
            :formatter="protoFormat">
        </el-table-column>

        <el-table-column
            prop="info"
            label="Details">
        </el-table-column>        

        <el-table-column
            prop="created_at"
            label="Created at"
            width="160"
            :formatter="tableDateFormat">
        </el-table-column>
      </el-table>

      <div class="sh-20"></div>

      <el-pagination
          background
          layout="prev, pager, next"
          :pager-count="11"
          :current-page.sync="currentPage"
          @current-change="pageChange"
          :total="count">
      </el-pagination>
  </div>
</template>

<script>
import axios from "axios";

export default {    
  name: "auditAccess",
  mixins: [],
  data() {
    return {
      tableData: [],
      count: 10,
      currentPage: 1,
      idSort: 1,
      activeName: "first",
      accessProtoArr:["", "UDP", "TCP", "HTTPS", "HTTP"], 
      defSearchForm: {username:'', src:'', dst:'', dst_port:'', access_proto:'', info:'', date:["",""]},
      searchForm: {},
      access_proto: [
            { text: 'UDP', value: '1' },
            { text: 'TCP', value: '2' },
            { text: 'HTTPS', value: '3' },
            { text: 'HTTP', value: '4' },
      ],
      maxExportNum: 1000000,
      loading: false,
      rules: {
        username: [
          {max: 30, message: 'Username must be less than 30 characters long', trigger: 'blur'}
        ],
        src: [
          {  message: 'Please enter the correct IP address', validator: this.validateIP, trigger: 'blur' },
        ],        
        dst: [
          { message: 'Please enter the correct IP address', validator: this.validateIP, trigger: 'blur' },
        ],               
      },            
    }
  },
  watch: {
    idSort: {
        handler(newValue, oldValue) {
            if (newValue != oldValue) {
                this.getData(1);
            }
        },
    },
  },  
  methods: {
    setSearchData() {
        this.searchForm = JSON.parse(JSON.stringify(this.defSearchForm));
    },    
    handleSearch() {
      this.$refs["searchForm"].validate((valid) => {
        if (!valid) {
          console.log('error submit!!');
          return false;
        }
        this.getData(1)
      })          
    },
    searchEnterFun(e) {
        var keyCode = window.event ? e.keyCode : e.which;
        if (keyCode == 13) {
            this.handleSearch()
        }
    },        
    getData(p) {
      this.loading = true
      if (! this.searchForm.date) {
        this.searchForm.date = ["", ""];
      }
      this.searchForm.sort = this.idSort     
      axios.get('/set/audit/list', {
        params: {
          page: p,
          search: this.searchForm,
        }
      }).then(resp => {
        var data = resp.data.data
        console.log(data);        
        this.tableData = data.datas;
        this.count = data.count
        this.loading = false
        this.currentPage = p;
      }).catch(error => {
        this.$message.error('Request error');
        console.log(error);
      });
    },
    pageChange(p) {     
        this.getData(p)
    },
    handleExport() {
      if (this.count > this.maxExportNum) {
        var formatNum = (this.maxExportNum + "").replace(/\d{1,3}(?=(\d{3})+$)/g,function(s){
           return s+','
        })
        this.$message.error("The amount of data you exported exceeds " + formatNum + " items. Please adjust the search conditions and export again.");
        return ;
      }
      if (! this.searchForm.date) {
        this.searchForm.date = ["", ""];
      }
      const exporting = this.$loading({
            lock: true,
            text: 'Exporting, please wait a moment...',
            spinner: 'el-icon-loading',
            background: 'rgba(0, 0, 0, 0.7)'
      });
      axios.get('/set/audit/export', {
        params: {
          search: this.searchForm,
        }
      }).then(resp => {
        var rdata = resp.data
        if (rdata.code && rdata.code != 0) {
            exporting.close();
            this.$message.error(rdata.msg);
            return ;
        }
        exporting.close();
        this.$message.success("Successfully exported CSV file")
        let csvData = 'data:text/csv;charset=utf-8,\uFEFF' + rdata
        this.createDownLoadClick(csvData, `anylink_audit_log_` + Date.parse(new Date()) + `.csv`)
      }).catch(error => {
        exporting.close();
        this.$message.error('Request error');
        console.log(error);
      });
    },
    createDownLoadClick(content, fileName) {
        const link = document.createElement('a')
        link.href = encodeURI(content)
        link.download = fileName
        document.body.appendChild(link)
        link.click()
        document.body.removeChild(link)
    },    
    protoFormat(row) {
        var access_proto = row.access_proto
        if (row.access_proto == 0) {
            switch (row.protocol) {
                case 6: access_proto = 2; break;
                case 17: access_proto = 1; break;
            }
        }
        return this.accessProtoArr[access_proto]
    },
    rest() {
        console.log("rest");
        this.setSearchData();
        this.handleSearch();
    }, 
    validateIP(rule, value, callback) {
        if (value === '' || typeof value === 'undefined' || value == null) {
            callback()
        } else {
            const reg = /^(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])$/
            if ((!reg.test(value)) && value !== '') {
            callback(new Error('Please enter the correct IP address'))
            } else {
            callback()
            }
        }
    }, 
    sortChange(column) {
        let { order } = column;
        if (order === 'ascending') {
            this.idSort = 2;
        } else {
            this.idSort = 1;
        }
    },            
  },
}
</script>

<style scoped>
.el-form-item {
    margin-bottom: 5px;
}
.el-table {
    font-size: 12px;
}
.search-form >>> .el-form-item__label {
  font-size: 12px;
}
/* /deep/  */
.el-table th {
    padding: 5px 0;
}
</style>