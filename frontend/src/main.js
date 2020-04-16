import Vue from 'vue'
import App from './App.vue'
import HelloWorld from './components/HelloWorld'
// Vue.prototype.$http = axios
// window.axios = axios
import 'materialize-css/dist/css/materialize.css'
import 'materialize-css/dist/js/materialize.js'
import VueRouter from 'vue-router'
import Tasks from  './components/Tasks'
import NewFile from './components/Tasks'

Vue.use(VueRouter);

var router = new VueRouter({
  mode: 'history',
  router: [
    {
      path: '/',
      name: 'tasks',
      component: Tasks
    },
    {
      path: '/tasks',
      name: 'tasks',
      component: Tasks
    },
    {
      path: '/helloworld',
      name: 'helloworld',
      component: HelloWorld
    },
    {
      path: '/newfile',
      name: 'newfile',
      component: NewFile
    }
  ]
});

new Vue({
  // el: '#app',
  render: h => h(App),
  router: router
}).$mount('#app');

// new Vue({
//   el: '#app',
//   router: router
// });