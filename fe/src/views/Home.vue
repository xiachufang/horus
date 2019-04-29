<template>
  <div class="home">
    <div class="tools-bar">
      <div style="flex: 1">
        Filter:
        <input
        class="search"
        v-model="keyword"
        placeholder="输入过滤内容"
        type="search" />
      </div>
      <div>
        展示最近：
        <input
          class="input"
          type="number"
          step="1"
          min="1"
          v-model="maxCount" > &nbsp;条
      </div>
    </div>
    <table class="table">
      <thead>
        <tr>
          <th
            v-for="(column, index) in columns"
            :key="index"
            @click="toggleSort(column)">
            {{column.label}}
            <span v-if="sortColumn === column">
              {{ orders === 'asc' ? '▲' : '▼' }}
            </span>
          </th>
          <th>source</th>
        </tr>
      </thead>
      <transition-group tag="tbody" name="list">
        <tr
          v-for="(msg) in (filteredList.length ? filteredList : sortedList)"
          :key="msg['@timestamp']">
          <td
            v-for="(column, index) in columns"
            :key="index">
            {{pathValue(column.path, msg)}}
          </td>
          <td>
            <div
              class="prop-item"
              v-for="(value, key) in msg.event.properties"
              :key="key"><label>{{key}}</label>: {{value}}</div>
          </td>
        </tr>
      </transition-group>
    </table>
  </div>
</template>

<script>
import _ from 'lodash'

export default {
  name: 'home',
  data () {
    return {
      messages: [],
      keyword: '',
      sortColumn: null,
      orders: 'desc',
      maxCount: 2000,
      columns: window['columns'] || [
        {
          'label': 'Time',
          'path': 'event.time'
        },
        {
          'label': 'distinct_id',
          'path': 'event.distinct_id'
        },
        {
          'label': 'type',
          'path': 'event.type'
        },
        {
          'label': 'event',
          'path': 'event.event'
        }
      ]
    }
  },
  created () {
    const url = process.env.NODE_ENV === 'development'
      ? 'wss://horus-x.xiachufang.com/topic'
      : ((window.location.protocol === 'https:') ? 'wss://' : 'ws://') + window.location.host + '/topic'

    const socket = new WebSocket(url)
    socket.addEventListener('open', (event) => {
      // console.log('Hello Server!')
    })
    socket.addEventListener('error', (event) => {
      alert('socket 连接失败')
    })

    socket.addEventListener('message', this.onMessage)
  },
  computed: {
    sortedList () {
      if (!this.sortColumn) return this.messages
      return _.orderBy(this.messages, this.sortColumn.path, this.orders)
    },
    filteredList () {
      if (!this.keyword) return []
      return this.messages.filter((msg) => {
        return JSON.stringify(msg).indexOf(this.keyword) > -1
      })
    }
  },
  methods: {
    onMessage (event) {
      const result = JSON.parse(event.data)
      result.event = JSON.parse(result.event)
      if (result.event && result.event.properties) {
        result.event.properties = this.sortByKeys(result.event.properties)
      }
      this.messages.unshift(result)
      this.messages = this.messages.slice(0, this.maxCount)
    },
    sortByKeys (object) {
      const keys = Object.keys(object)
      const sortedKeys = _.sortBy(keys)
      return _.fromPairs(_.map(sortedKeys, key => [key, object[key]]))
    },
    toggleSort (column) {
      this.sortColumn = column
      if (this.orders === 'desc') {
        this.orders = 'asc'
      } else {
        this.orders = 'desc'
      }
    },
    pathValue (path, val) {
      return _.get(val, path, '')
    }
  }
}
</script>

<style lang="stylus">
  .prop-item
    display inline-block
    margin 0 10px 0 0
    label
      font-weight bold

  .table
    border-collapse: collapse;
    border none
    width 100%;
    margin-bottom: 1rem
    color #212529
    font-size 12px

  .table td,.table th
    padding .75rem
    vertical-align top
    border-top 1px solid #dee2e6
    min-width 100px

  .table thead th
    vertical-align bottom
    border-bottom 2px solid #dee2e6
    text-align left

  .list-item
    background-color transparent
  .list-enter-active, .list-leave-active
    transition all 1s
  .list-enter, .list-leave-to
    background-color pink
  .tools-bar
    display flex
    flex-wrap wrap
    flex-direction row
    align-items center
    position sticky
    top 0
    padding 20px
    background-color #ffffff
  .tools-bar > div
    display flex
    flex-wrap wrap
    flex-direction row
    align-items center
  .search, .input
    display block
    height 30px
    line-height 20px
    padding 0 10px
    -webkit-appearance none
    border 1px solid #888
    border-radius 5px
    font-size 14px
  .search
    width 50%
    min-width 200px
    margin 0 20px
</style>
