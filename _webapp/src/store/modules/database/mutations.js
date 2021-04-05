import {set} from '@/utils/vuex';


export default {
    setDatabases: (state, payload) => {
        state.databases = payload;
    },
    setDatabaseModel: (state, payload) => {
        for(let i = 0;i<payload.tables.length;i++){
            // payload.tables[i].position = {x:Math.round(Math.random() * 1000), y:Math.round(Math.random()*1000)};
            // payload.tables[i].position = {x:Math.round(i * 10), y:Math.round(10)};
            payload.tables[i].x = Math.round(i * 20);
            payload.tables[i].y = Math.round(i * 20);
            payload.tables[i].position = {};
            payload.tables[i].position.y = Math.round(i * 20);
            payload.tables[i].position.x = Math.round(i * 20);
            payload.tables[i].height = payload.tables[i].columns.length*16+48;
            let minWidth = payload.tables[i].name.length*8;
            payload.tables[i].width = minWidth>200?minWidth:200;
        }
        // let tables = payload.tables.slice(10,20);
        // payload.tables = tables;
        state.databaseModel = payload;

    },
    setModel: (state, payload) =>{
        state.databaseModel = payload;
    },


    setCanCreateNew: set('canCreateNew'),
    setSelectedBusiness: set('selectedBusiness'),
    setBusinessDetail: (state, payload) => {
        state.businessDetail = payload;
    },
    setCurrentUserBusinesses: (state, payload) => {
        state.businessLoaded = true;
        state.currentBusinesses = payload;
    },
    setLocationDetails: (state, payload) => {
        if (payload != null) {
            state.registration.locationData = payload;
        }
    }
}