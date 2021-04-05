// https://vuex.vuejs.org/en/getters.html

export default {
    getAccessToken: (state) => {
        return state.access_token;
    }
}
