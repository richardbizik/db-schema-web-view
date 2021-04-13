import Vue from 'vue'
import VueRouter from 'vue-router'
import Home from '../views/Home.vue'
import Database from "@/views/Database";

Vue.use(VueRouter)

const routes = [
  {
    path: '/',
    name: 'Home',
    component: Home
  },
  {
    path: '/database/:dbName/:schema',
    name: 'Model',
    component: Database
  }
]

const router = new VueRouter({
  routes
})

export default router
