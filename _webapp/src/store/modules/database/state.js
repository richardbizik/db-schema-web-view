export default {
    databases: [],
    databaseModel: {name:"", columns:[]},

    canCreateNew: false,
    currentBusinesses: [],
    businessLoaded: false,
    selectedBusiness: null,
    businessDetail: null,
    registration: {
        newBusinessForm: {
            name: "",
            description: "",
            businessHours: [],
            tags: [],
            location: null
        },
        locationData: null
    }
}
