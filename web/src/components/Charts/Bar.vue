<template>
  <div :style="{ padding: '0 0 32px 32px' }">
    <h4 :style="{ marginBottom: '20px' }">{{ title }}</h4>
    <v-chart
      height="254"
      :data="data"
      :forceFit="true"
      :padding="['auto', 'auto', '40', '50']">
      <v-tooltip />
      <v-axis />
      <v-bar position="time*warnings"/>
    </v-chart>
  </div>
</template>

<script>
export default {
  name: 'Bar',
  props: {
    title: {
      type: String,
      default: ''
    },
    data: {
      type: Array,
      default: () => {
        return []
      }
    },
    scale: {
      type: Array,
      default: () => {
        return [{
          dataKey: 'time',
          min: 2
        }, {
          dataKey: 'warnings',
          title: 'warning items',
          min: 1,
          max: 22
        }]
      }
    },
    tooltip: {
      type: Array,
      default: () => {
        return [
          'time*warnings',
          (time, warnings) => ({
            name: time,
            value: warnings
          })
        ]
      }
    }
  },
  data () {
    return {
    }
  }
}
</script>
