<template>
  <div class="line-chart" :style="{height:height,width:width}" />
</template>

<script>
  import echarts from 'echarts'

  export default {
    name: 'LineChart',
    props: {
      width: {
        type: String,
        default: '100%'
      },
      height: {
        type: String,
        default: '240px'
      },
      chartData: {
        type: Object,
        required: true
      }
    },
    data() {
      return {
        chart:null
      }
    },
    mounted() {
      this.initChart()
      window.addEventListener('resize', this.listenerResize)
    },
    beforeDestroy() {
      window.removeEventListener('resize', this.listenerResize)
    },
    watch: {
        chartData:{
            handler(){
                this.initChart()
            },
            deep:true
        }
    },
    methods: {
      listenerResize() {
        this.chart.resize()
      },
      initChart() {
        this.chart = echarts.init(this.$el)

        let option = {
          color: ['#2D5CF6','#50B142'],
          title: {
            text: this.chartData.title || 'Line chart',
            textStyle:{fontWeight:'normal',fontSize:16, color:'#8C8C8C'}
          },
          tooltip: {
            trigger: 'axis'
          },
          legend: {
            bottom: 0,
          },
          grid: {
            left: '3%',
            right: '4%',
            bottom: '10%',
            containLabel: true
          },
          xAxis: {
            type: 'category',
            boundaryGap: false,
            data: this.chartData.xname, 
            splitLine: {
                show: false
            }                                                     
          },
          yAxis: {
            type: 'value',
            minInterval: undefined,
            name: this.chartData.yname,
            splitLine: {
                lineStyle: {
                    color: "#F0F0F0",
                },
            }     
          },
          series: [],
        };
        if (this.chartData.yminInterval != undefined) {
            option.yAxis.minInterval = this.chartData.yminInterval
        }
        let xdata = this.chartData.xdata
        for (let key in xdata) {
          let a = {
            name: key,
            type: 'line',
            showSymbol: false,
            data: xdata[key]
          };
          option.series.push(a)
        }
        this.chart.setOption(option)
      },
    }    
  } 
</script>
