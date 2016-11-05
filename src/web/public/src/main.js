import Vue from 'vue'
import App from './App'
import Resource from 'vue-resource'

/* eslint-disable no-new */
new Vue({
  el: '#app',
  template: '<App/>',
  components: { App }
})

Vue.use(Resource)
