<template>
  <div class="home">
        <HelloWorld msg="Welcome to Your Vue.js App"/>
  </div>
</template>

<script>
// @ is an alias to /src

import {mapState} from "vuex";
import HelloWorld from "@/components/HelloWorld";

export default {
  name: 'Home',
  components: {
    HelloWorld
  },
  computed: {
    ...mapState('database', ['databases', 'databaseModel']),
  },
  methods: {
  },
  data: () => ({tables: [{name: 'someTable'}]}),
}
</script>

<style>

.table {
  padding: 5px 5px;
  border: 1px solid;
  border-radius: 5px;
}

#canvas {
  width: 100%;
  height: calc(100vh - 64px);
}
</style>
