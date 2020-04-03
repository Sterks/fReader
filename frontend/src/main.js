import Vue from 'vue'
// import App from './App.vue'
// Vue.prototype.$http = axios
// window.axios = axios
import 'materialize-css/dist/css/materialize.css'
import 'materialize-css/dist/js/materialize.js'
import VueRouter from 'vue-router'
import Tasks from  './components/Tasks'
Vue.use(VueRouter);

var router = new VueRouter({
  router: [
    {
      path: '/tasks',
      name: 'tasks',
      component: Tasks
    }
  ]
});

new Vue({
  // el: '#app',
  render: h => h(Tasks),
  router: router
}).$mount('#app');

// new Vue({
//   el: '#app',
//   router: router
// });