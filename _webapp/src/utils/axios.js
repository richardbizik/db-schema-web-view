import axios from 'axios'
import store from '../store'


export default () => {

    function getHeaders() {
        var headers = {};

        var token = store.getters['user/getAccessToken'];
        if (token) {
            headers['Authorization'] = "bearer " + token;
        }
        var tokenHeader = document.cookie.match(new RegExp(`XSRF-TOKEN=([^;]+)`));
        if (tokenHeader) {
            headers['X-XSRF-TOKEN'] = tokenHeader[1];
        }
        return headers;
    }

    return axios.create({
        baseURL: process.env.VUE_APP_GATEWAY_ENDPOINT,
        // withCredentials: true,
        headers: getHeaders()
    });

}