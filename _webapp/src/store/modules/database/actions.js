import axios from '@/utils/axios';
import * as notificationUtils from '@/utils/notifications';
import router from '../../../router';

export default {
    async getDatabases(context) {
        axios().get("/database").then(result => {
            context.commit('setDatabases', result.data);
        }).catch((error) => {
            console.error("error while getting database list", error);
        });
    },
    getDatabaseModel(context, payload) {
        axios().get("/database/"+payload.dbName+"/"+payload.schema).then(result => {
            context.commit('setDatabaseModel', result.data);
        }).catch((error) => {
            console.error("error while getting database model", error);
        });
    },
    setModel(context, payload){
        context.commit('setModel', payload);
    },
}
