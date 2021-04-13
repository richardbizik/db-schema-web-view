<template>
  <v-app id="inspire">
    <v-navigation-drawer id="nav" v-model="inputValue" app>
      <div class="navigation">
        <div>
        <v-list-item link to="/" style="min-height: 64px">
          <v-list-item-content>
            <v-list-item-title class="title">
              Schema web view
            </v-list-item-title>
          </v-list-item-content>
        </v-list-item>
        <v-divider></v-divider>
        </div>
        <div class="database-list">
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
        </div>
      </div>
      <div class="navigation-spacer"></div>
      <div class="table-list">
        <v-list
            dense
        >
          <v-list-item
              class="tables-item"
              v-for="table in tables"
              :key="table.name"
          >
            <v-list-item-content @click="setSelectedTable(table.name)" style="display: block;">
              <div style="float: left;max-width: min-content;font-size:.9rem">
              {{ table.name }}
              </div>
              <div style="float:right;max-width: min-content;color: #d77700;font-weight: bold;" v-if="!table.hasRelation">
                <v-tooltip right>
                  <template v-slot:activator="{ on, attrs }">
                    <div v-bind="attrs" v-on="on" style="font-size:.9rem;min-width: 10px;text-align: center">!</div>
                  </template>
                  <span>Table is not connected to other tables</span>
                </v-tooltip>
              </div>
            </v-list-item-content>
          </v-list-item>
        </v-list>
      </div>
    </v-navigation-drawer>
    <v-app-bar app style="z-index: 9999999">
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
    ...mapState('database', ['databases', 'databaseModel']),
    tables: {
      get() {
        if(this.databaseModel===undefined||this.databaseModel.tables===undefined){
          return [];
        }
        let tables = [];
        let tablesWithRelation = [];
        for(let table of this.databaseModel.tables){
          for(let column of table.columns){
            for(let fk of column.fkConstraints){
              if(fk.rTable!==table.name) {
                if (tablesWithRelation.indexOf(table.name) < 0) {
                  tablesWithRelation.push(table.name);
                }
                if (tablesWithRelation.indexOf(fk.rTable) < 0) {
                  tablesWithRelation.push(fk.rTable);
                }
              }
            }
          }
        }

        for(let table of this.databaseModel.tables){
          tables.push({name:table.name,hasRelation: tablesWithRelation.indexOf(table.name) > -1});
        }
        tables.sort(function (a, b) {
          return (a.name.toUpperCase() < b.name.toUpperCase()) ? -1 : 1
        })
        return tables;
      }
    },
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
    },
    setSelectedTable: function(table){
      this.$eventHub.$emit("select-table", table)
    }
  },
  mounted() {
    this.$store.dispatch('database/getDatabases');
  },
}
</script>

<style>
.menu-item .v-list-item__content {
  font-size: 0.6rem;
  padding: 2px 0 2px 0 !important;
}

.menu-item .v-list-item__title {
  font-size: 1rem !important;
  line-height: 1rem !important;
}

.database-list {
  overflow-y: auto !important;
  max-height: calc(100% - 47px);
}

.navigation-spacer {
  min-height: 5%;
}

.table-list {
  background-color: #ececec;
  height: 35%;
  overflow-y:scroll;
}
.tables-item .v-list-item__content {
  line-height: 0.8rem;
  padding: 1px 5px 1px 5px !important;
}
.tables-item:hover{
  background-color: #a1a1a1;
  border-radius: 5px;
  cursor: pointer;
}

.navigation {
  overflow-y: hidden !important;
  height: 60%;
  padding-bottom: 10%;
}

.database-list::-webkit-scrollbar {
  width: 8px; /* width of the entire scrollbar */
}

.database-list::-webkit-scrollbar-track {
  background: #ffffff; /* color of the tracking area */
}
.v-navigation-drawer__content{
  overflow-y: hidden !important;
}
.database-list::-webkit-scrollbar-thumb {
  background-color: #767676; /* color of the scroll thumb */
  border-radius: 20px; /* roundness of the scroll thumb */
  border: 1px solid #ffffff; /* creates padding around scroll thumb */
}
.tables-item{
  padding-left: 5px !important;
  padding-right: 5px !important;
  min-height: 10px !important;
  font-size: .7rem !important;
}
.table-list::-webkit-scrollbar {
  width: 8px; /* width of the entire scrollbar */
}

.table-list::-webkit-scrollbar-track {
  background: #ffffff; /* color of the tracking area */
}

.table-list::-webkit-scrollbar-thumb {
  background-color: #767676; /* color of the scroll thumb */
  border-radius: 20px; /* roundness of the scroll thumb */
  border: 1px solid #ffffff; /* creates padding around scroll thumb */
}
</style>
