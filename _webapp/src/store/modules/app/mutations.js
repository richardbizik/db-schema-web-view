import {set, toggle} from '@/utils/vuex'

export default {
    setDrawer: set('drawer'),
    toggleDrawer: toggle('drawer'),
    addNotification: (state, payload) => {
        payload.key = Math.random();
        state.notifications.push(payload);
    },
    removeNotification: (state, payload) => {
        state.notifications = state.notifications.map((notification) => {
            if(notification.key==payload){
                notification.snackbar = false;
                setTimeout(()=>{
                    state.notifications = state.notifications.filter(value => value.key !== payload);
                }, 200)
            }
            return notification;
        });
        console.log(state.notifications)
    }
}
