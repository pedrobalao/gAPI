import Vue from "vue";
import Vuex from "vuex";
var auth = require('./auth')
Vue.use(Vuex);


import users from './store/modules/users/index';
import serviceGroups from './store/modules/services-groups'
import fullscreen from './store/modules/full-screen'

export default new Vuex.Store({
  modules:{
    users,
    serviceGroups,
    fullscreen
  },
  state: {
    isLoggedIn : false,
    loggedInUser: null
  },
  mutations: {
    loggedIn (state) {
      state.isLoggedIn = true
    },
    loggedOut (state) {
      state.isLoggedIn = false
    },
    loggedInUserUpdate (state,user) {
      state.loggedInUser = user
      localStorage.setItem('user', user)
    }
  },
  getters: {
    isLoggedIn: state => {
      return state.isLoggedIn
    },
    loggedInUser: state => {
      return state.loggedInUser
    },
    isAdmin: state => {
      if (! state.loggedInUser) return false
      return state.loggedInUser.IsAdmin
    }
  },
  actions: {
    loggedInUserUpdate: ({commit}, user) => {
      commit("loggedInUserUpdate", user)
    },
    login : ({ commit }) => {
      commit("loggedIn")
    },
    logout : ({ commit }) => {
      commit("loggedOut")
    }
  }
});
