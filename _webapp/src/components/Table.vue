<template v-on:hover-table="hoverTable">
  <v-card class="db_table" :style="styles()" @mouseenter="tableEntered(tableData.name)"
          @mouseleave="tableLeft(tableData.name)">
    <v-tooltip right>
      <template v-slot:activator="{ on, attrs }">
        <v-card-title
            v-on="on"
            v-bind="attrs"
        >{{ tableData.name }}
        </v-card-title>
      </template>
      <div class="table-info">
        <v-card-title>Table info</v-card-title>
        <div class="unique-constraints-title" v-if="tableData.primary.length>0">Primary key:</div>
        <v-list-item v-if="tableData.primary.length>0">
          <span class="table-info-text">{{ tableData.primary.join(", ") }}</span>
        </v-list-item>
        <div class="unique-constraints-title" v-if="tableData.uniques.length>0">Unique constraints:</div>
        <v-list-item v-for="uq in tableData.uniques" v-bind:key="uq.name">
          <span class="table-info-text">{{ uq.name }} - ({{ uq.columns.join(", ") }})</span>
        </v-list-item>
        <div class="unique-constraints-title" v-if="tableData.foreigns.length>0">Foreign constraints:</div>
        <v-list-item v-for="fk in tableData.foreigns" v-bind:key="fk.name">
          <span class="table-info-text">{{ fk.name }} ({{ fk.columns.join(", ") }}) -> {{fk.targetTable}} ({{fk.targetColumn}})</span>
        </v-list-item>
      </div>
    </v-tooltip>
    <v-list-item v-for="column in tableData.columns" v-bind:key="column.position" class="dbcolumn">
      <div :ref="column.name" class="column" style="width: 100%; display:block;">
        <div v-if="column.primary" style="text-align: left;float:left;padding-right: 5px;width:20px;">PK</div>
        <div v-if="column.nullable" style="text-align: left;float:left;padding-right: 5px;width:20px;">N</div>
        <div v-if="!column.nullable&&!column.primary" style="text-align: left;float:left;width:20px;height: 16px">
        </div>
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
    }],
    hover: false,
  }),
  methods: {
    styles: function () {
      let styles = "top: " + this.tableData.y + "px;left: " + this.tableData.x + "px;width: " + this.tableData.width + "px;cursor:auto;"
      if (this.hover || this.tableData.hover) {
        styles += "box-shadow: 2px 2px 10px 0px;"
      }
      return styles;
    },
    tableEntered: function (tableName) {
      this.hover = true;
      this.$emit("table-entered", tableName);
    },
    tableLeft: function (tableName) {
      this.hover = false;
      this.$emit("table-left", tableName);
    },
  },
  mounted() {
    this.$emit("table-initialized", this.tableData.name);
  }
}
</script>

<style>
.v-card {
  padding: 0px 0px !important;
  z-index: 1000;
}

.v-card__title {
  display: block !important;
  text-align: center !important;
  border-bottom: 1px solid black;
  padding: 5px 5px !important;
  font-size: .8rem !important;
}

.dbcolumn.v-list-item {
  padding-left: 5px !important;
  padding-right: 5px !important;
  min-height: 10px !important;
  font-size: .7rem !important;

}

.db_table {

}

.table-info-text {
  color: white !important;
}

.table-info .v-list-item {
  padding-left: 5px !important;
  padding-right: 5px !important;
  min-height: 10px !important;
  font-size: .8rem !important;
}

.table-info .v-card__title {
  border-bottom: none;
}
</style>
