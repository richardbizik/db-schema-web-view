

function findPos(obj) {
    var curtop = 0;
    if (obj.offsetParent) {
        do {
            curtop += obj.offsetTop;
        } while (obj == obj.offsetParent);
        return [curtop];
    }
}
export default {

    async applicationLoaded(context) {
        context.dispatch('user/applicationLoaded', {}, {root:true});
    },
    scrollToElement(context, data){
        let elmnt = document.getElementById(data.id);
        window.scroll(0, findPos(elmnt)-data.offset);
    },
    addNotification(context, data){
        let notification = {
            color: data.color,
            top: false,
            bottom: true,
            left: false,
            right: true,
            snackbar: true,
            text: data.text,
            timeout: data.timeout
        }
        context.commit("addNotification", notification);
    },
    removeNotification(context, key){
        context.commit("removeNotification", key);
    }
}