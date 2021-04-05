<template>
  <v-app id="inspire">
    <v-navigation-drawer id="nav" v-model="inputValue" app>
      <v-list-item link to="/">
        <v-list-item-content>
          <v-list-item-title class="title">
            Application
          </v-list-item-title>
          <v-list-item-subtitle>
            subtext
          </v-list-item-subtitle>
        </v-list-item-content>
      </v-list-item>
      <v-divider></v-divider>

      <v-list
          dense
          nav
      >
        <v-list-item
            class="menu-item"
            v-for="database in databases"
            :key="database.name+database.schema"
            link
            :to="getDBLink(database)"
        >
          <v-list-item-content>
            <v-list-item-title>{{ database.name }}</v-list-item-title>
            <v-list-item-content>{{ database.schema }}</v-list-item-content>
          </v-list-item-content>
        </v-list-item>
      </v-list>
    </v-navigation-drawer>
    <v-app-bar app>
      <v-app-bar-nav-icon @click="toggleDrawer"></v-app-bar-nav-icon>
      <v-toolbar-title>Application</v-toolbar-title>
    </v-app-bar>
    <v-main>
      <router-view/>
    </v-main>
  </v-app>
</template>

<script>
import {mapMutations, mapState} from "vuex";

export default {
  data: () => ({
    drawer: null,
    items: [
      {title: 'Dashboard', icon: 'mdi-view-dashboard', path: '/'},
      {title: 'About', icon: 'mdi-help-box', path: '/about'},
    ],
    right: null,
  }),
  watch: {
    '$route'(to) {
      document.title = to.name
    }
  },
  computed: {
    ...mapState('database', ['databases']),
    inputValue: {
      get() {
        return this.$store.state.app.drawer
      },
      set(val) {
        this.setDrawer(val)
      }
    },
  },
  methods: {
    ...mapMutations('app', ['setDrawer', 'toggleDrawer']),
    getDBLink: function (db) {
      return "/database/" + db.name + "/" + db.schema
    }
  },
  mounted() {
    this.$store.dispatch('database/getDatabases');
  },
}
</script>

<style>
.menu-item .v-list-item__content{
  font-size: 0.8rem;
}
.menu-item .v-list-item__title{
  font-size: 1.5rem !important;
  line-height: 1.5rem !important;
}
.v-navigation-drawer__content::-webkit-scrollbar {
  width: 8px;               /* width of the entire scrollbar */
}

.v-navigation-drawer__content::-webkit-scrollbar-track {
  background: #ffffff;        /* color of the tracking area */
}

.v-navigation-drawer__content::-webkit-scrollbar-thumb {
  background-color: #767676;    /* color of the scroll thumb */
  border-radius: 20px;       /* roundness of the scroll thumb */
  border: 1px solid #ffffff;  /* creates padding around scroll thumb */
}
</style>
