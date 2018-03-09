<template>
  <div id="app">
    <h1>Control Panel <sub>Project Spartan</sub></h1>
    <h3>HAProxy</h3>
    <table border="1">
      <thead>
        <tr>
          <th>Name</th>
          <th>Address</th>
          <th>Status</th>
          <!-- <th>Action</th> -->
        </tr>
      </thead>
      <tbody>
        <tr v-for="item in haproxies" :key="item.name">
          <td>{{ item.name }}</td>
          <td>{{ item.address }}</td>
          <td>{{ item.status }}</td>
        </tr>
      </tbody>
    </table>

    <h3>Application</h3>
    <table border="1">
      <thead>
        <tr>
          <th>Name</th>
          <th>Address</th>
          <th>Status</th>
          <!-- <th>Action</th> -->
        </tr>
      </thead>
      <tbody>
        <tr v-for="item in servers" :key="item.name">
          <td>{{ item.name }}</td>
          <td>{{ item.address }}</td>
          <td>{{ item.status }}</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script>
export default {
  name: 'App',

  data: function () {
    return {
      haproxies: [],
      servers: []
    }
  },

  created: function () {
    this.loadData()

    setInterval(function () {
      this.loadData()
    }.bind(this), 3000)
  },

  methods: {
    loadData: function () {
      this.$http.get('/api/dashboard').then(function (response) {
        this.haproxies = response.data ? response.data.haproxies : []
        this.servers = response.data ? response.data.servers : []
      })
    }
  }
}
</script>

<style>
#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
table {
  width: 80%;
  margin: auto;
}
table th {
  padding: 3px;
}
table input {
  width: 80%;
}
</style>
