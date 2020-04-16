<template>
  <div id="app">
    <div class="row">
      <router-link to="/tasks">Перейти к Foo</router-link>
      <router-view></router-view>
      <div class="col s12 m3">
        <div class="card blue-grey darken-1">
          <div class="card-content white-text">
            <span class="card-title">Информация</span>
            <p>Всего файлов обработано: {{lastID}}</p>
            <p>Извещений: {{notification}}</p>
            <p>Протоколов: {{protocols}}</p>
            <p>Итого: {{total = notification + protocols}} </p>
          </div>
        </div>
      
      </div>
    </div>
    <div class="row">
      <div class="card-body col s12 m3">
        <form @submit="tasks">
          <div class="input-field col s12">
            <input type="text" id="name" v-model="name">
            <label for="name">Имя</label>
          </div>
          <div class="input-field col s12">
            <input type="text" id="family" v-model="family">
            <label for="family">Фамилия</label>
          </div>
          <button class="btn btn-success">Отправить</button>
        </form>
        <strong>Output:</strong>
        <pre>
           {{output}}
        </pre>
      </div>
    </div>
  </div>
</template>

<script>
  import Axios from 'axios'

  export default {
    data() {
      return {
        lastID: {},
        notification: {},
        protocols: {}
      }
    },
    mounted() {
      this.updateLastID();
      setInterval(this.updateLastID, 10000);
      this.tasks();
      // this.showBlock();
    },
    methods: {
      updateLastID () {
        Axios.get(`http://127.0.0.1:8000/GetID`).then(response => {
          this.lastID = response.data.LastID;
          this.notification = response.data.Notification;
          this.protocols = response.data.Protocols;
        })
      },
      tasks () {
        Axios.post(`http://localhost`)
      }
      // showBlock () {
      //   // this.$el.className("card").style.transform = "rotate("+window.pageYOffset+"deg)"
      // }
    }
  }

</script>
