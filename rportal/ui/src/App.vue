<template>
  <div id="app">
    <h1>Reseller Portal <sub>Project Spartan</sub></h1>
    <table border="1">
      <thead>
        <tr>
          <th>ID</th>
          <th>Reseller Name</th>
          <th>Contact Person</th>
          <th>Phone Number</th>
          <th>Business Address</th>
          <th>Action</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="item in items" :key="item.id">
          <td>{{ item.id }}</td>
          <td>{{ item.name }}</td>
          <td>{{ item.person }}</td>
          <td>{{ item.phone }}</td>
          <td>{{ item.address }}</td>
          <td><a href="#" @click="deleteItem(item)">Delete</a></td>
        </tr>
      </tbody>
    </table>

    <br><br>
    <table border="1">
      <tbody>
        <tr>
          <th>Reseller Name:</th>
          <td><input type="text" v-model="newItem.name"></td>
        </tr>
        <tr>
          <th>Contact Person:</th>
          <td><input type="text" v-model="newItem.person"></td>
        </tr>
        <tr>
          <th>Phone Number:</th>
          <td><input type="text" v-model="newItem.phone"></td>
        </tr>
        <tr>
          <th>Business Address:</th>
          <td><input type="text" v-model="newItem.address"></td>
        </tr>
      </tbody>
    </table>
    <button @click="createItem">Add New Reseller</button>
  </div>
</template>

<script>
export default {
  name: 'App',

  data: function () {
    return {
      items: [],
      newItem: {}
    }
  },

  created: function () {
    this.$http.get('/api/items').then(function (response) {
      this.items = response.data ? response.data : []
    })
  },

  methods: {
    createItem: function () {
      this.$http.post('/api/items', this.newItem).then(function (response) {
        this.newItem.id = response.data.id
        this.items.push(this.newItem)
        this.newItem = {}
      }, function () {
        alert('Something went wrong on the server!')
      })
    },

    deleteItem: function (item) {
      this.$http.delete('/api/items/' + item.id).then(function (response) {
        location.reload()
      }, function () {
        alert('Something went wrong on the server!')
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
