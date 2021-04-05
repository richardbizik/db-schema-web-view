<template>
  <v-card class="db_table" :style="styles()">
    <v-card-title>{{ tableData.name }}</v-card-title>
    <v-list-item v-for="column in tableData.columns" v-bind:key="column.position" >
      <div class="column" style="width: 100%; display:block;">
        <div v-if="column.primary" style="text-align: left;float:left;padding-right: 5px;width:20px;">PK</div>
        <div v-if="column.nullable" style="text-align: left;float:left;padding-right: 5px;width:20px;">N</div>
        <div v-if="!column.nullable&&!column.primary" style="text-align: left;float:left;width:20px;height: 16px">&nbsp;</div>
        <div style="float:left">{{ column.name }}</div>
        <div style="text-align: right;float:right">{{ column.type }}</div>
      </div>
    </v-list-item>
  </v-card>
</template>

<script>
export default {
  name: 'Table',
  props: ['tableData', 'position'],
  data: () => ({
    name: "",
    columns: [{
      default: "",
      name: "",
      nullable: false,
      position: 0,
      type: "",
      typeDetail: "",
    }]
  }),
  methods: {
    styles: function(){
      return "top: "+this.tableData.y+"px;left: "+this.tableData.x+"px;width: "+this.tableData.width+"px;";
    }
  }
}
</script>

<style>
.v-card{
  padding: 0px 0px !important;
  z-index: 1000;
}
.v-card__title{
  display: block !important;
  text-align: center !important;
  border-bottom: 1px solid black;
  padding: 5px 5px !important;
  font-size: .8rem !important;
}
.v-list-item {
  padding-left: 5px !important;
  padding-right: 5px !important;
  min-height: 10px !important;
  font-size: .7rem !important;
}
.db_table {

}
</style>
