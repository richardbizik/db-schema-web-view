import {set} from '@/utils/vuex';


export default {
    setDatabases: (state, payload) => {
        state.databases = payload;
    },
    setDatabaseModel: (state, payload) => {
        for (let i = 0; i < payload.tables.length; i++) {
            payload.tables[i].x = Math.round(-9999 + i * -99);
            payload.tables[i].y = Math.round(-9999 + i * -99);
            payload.tables[i].position = {};
            payload.tables[i].hover = false;
            payload.tables[i].position.y = Math.round(i * 20);
            payload.tables[i].position.x = Math.round(i * 20);
            payload.tables[i].height = payload.tables[i].columns.length * 16 + 48;
            let minWidth = payload.tables[i].name.length * 8;
            payload.tables[i].width = minWidth > 200 ? minWidth : 200;
            //fill table.uniques
            payload.tables[i].uniques = [];
            payload.tables[i].foreigns = [];
            payload.tables[i].primary = [];
            for (let column of payload.tables[i].columns) {
                if(column.primary) {
                    payload.tables[i].primary.push(column.name);
                }
                uniqueConstraints:
                    for (let unique of column.uniqueConstraints) {
                        //check if unique name is not present
                        for (let u of payload.tables[i].uniques) {
                            if (u.name === unique) {
                                u.columns.push(column.name);
                                continue uniqueConstraints;
                            }
                        }
                        payload.tables[i].uniques.push({name: unique, columns: [column.name]});
                    }
                foreignConstraints:
                    for (let foreignConstraint of column.fkConstraints) {
                        //check if unique name is not present
                        for (let f of payload.tables[i].foreigns) {
                            if (f.name === foreignConstraint.name) {
                                f.columns.push(column.name);
                                continue foreignConstraints;
                            }
                        }
                        payload.tables[i].foreigns.push({
                            name: foreignConstraint.name,
                            targetTable: foreignConstraint.rTable,
                            targetColumn: foreignConstraint.rColumn,
                            columns: [column.name]
                        });
                    }
            }

        }
        // let tables = payload.tables.slice(10,20);
        // payload.tables = tables;
        state.databaseModel = payload;

    },
    setModel: (state, payload) => {
        state.databaseModel = payload;
    },
}