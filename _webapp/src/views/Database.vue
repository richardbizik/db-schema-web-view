<template>
  <div>
    <v-overlay :value="overlay">
      <v-progress-circular
          indeterminate
          size="64"
      ></v-progress-circular>
    </v-overlay>
    <canvas ref="canvasElement"
            :height="this.$refs.canvasScreen?this.$refs.canvasScreen.clientHeight:0"
            :width="this.$refs.canvasScreen?this.$refs.canvasScreen.clientWidth:0"
            style="z-index: 2;position: absolute;pointer-events: none;"></canvas>
    <div id="canvasScreen" ref="canvasScreen" @mousedown="startMove" @mouseup="cancelMove" @mousemove="moveCanvas"
         @click.middle="displayGraph"
         @wheel.prevent="handleZoom">
      <div class="zoom-level">{{ Math.round((zoom / 100) * 10000) }}%</div>
      <div id="canvas" :style="getCanvasStyle">
        <v-list-item
            class="table-item"
            :style="isMousePressed?'user-select:none;':''"
            v-for="(table) in database.tables"
            :key="table.name"
            ref="list"
        >
          <Table class="table" :table-data="table"/>

        </v-list-item>
      </div>
    </div>
  </div>
</template>

<script>
// @ is an alias to /src

import {mapState} from "vuex";
import Table from "@/components/Table";
import * as d3 from "d3";

export default {
  name: 'Database',
  components: {Table},
  data: () => ({
    viewPosition: {x: 0, y: 0},
    zoom: 1,
    dragStart: {x: 0, y: 0, originalViewPosition: {x: 0, y: 0}},
    isMousePressed: false,
    simulation: null,
    database: [],
    links: [],
    overlay: false,
  }),
  computed: {
    ...mapState('database', ['databases', 'databaseModel']),
    getCanvasStyle() {
      return 'transform-origin: 0px 0px 0px; transform: matrix(' + this.zoom + ', 0, 0, ' + this.zoom + ', ' + this.viewPosition.x + ', ' + this.viewPosition.y + ';';
    }
  },
  methods: {
    getDatabases: function () {
      this.$store.dispatch('database/getDatabases');
    },
    getDatabaseModel: function (dbName, schema) {
      this.$store.dispatch('database/getDatabaseModel', {dbName, schema});
    },
    startMove: function (event) {
      if (event.path[0].id === "canvasScreen" || event.path[0].id === "canvas") {
        this.isMousePressed = true;
        if (this.isMousePressed) {
          this.dragStart.x = event.x
          this.dragStart.y = event.y
          this.dragStart.originalViewPosition.x = this.viewPosition.x
          this.dragStart.originalViewPosition.y = this.viewPosition.y
        }
      }
    },
    cancelMove: function () {
      this.isMousePressed = false;
    },
    moveCanvas: function (event) {
      if (this.isMousePressed) {
        // this.viewPosition.x = event.
        this.viewPosition.x = this.dragStart.originalViewPosition.x - (this.dragStart.x - event.x)
        this.viewPosition.y = this.dragStart.originalViewPosition.y - (this.dragStart.y - event.y)
        this.drawLines();
      }
    },
    clearLines: function () {
      let width = this.$refs.canvasScreen.clientWidth;
      let height = this.$refs.canvasScreen.clientHeight;
      const context = this.$refs.canvasElement.getContext("2d");
      context.clearRect(0, 0, width, height);
    },
    drawLines: function () {
      let width = this.$refs.canvasScreen.clientWidth;
      let height = this.$refs.canvasScreen.clientHeight;
      let tableThing = 15;
      let arrowSize = 5;
      let circleSize = 5;
      const context = this.$refs.canvasElement.getContext("2d");
      context.clearRect(0, 0, width, height);
      context.save();
      context.beginPath();
      for (const d of this.links) {
        if (Math.round(d.sourceX) < Math.round(d.targetX)) {
          let baseX = ((d.source.x + d.sourceWidth) * (this.zoom) + this.viewPosition.x);
          let baseY = (d.source.y + (d.sourceIndex * 16 + 5) + 48) * this.zoom + this.viewPosition.y;

          //arrow
          context.moveTo(baseX, baseY-arrowSize);
          context.lineTo(baseX+arrowSize, baseY);
          context.lineTo(baseX, baseY+arrowSize);

          context.moveTo(baseX, baseY);
          context.lineTo(baseX+tableThing, baseY);
          let targetX = (d.target.x * (this.zoom) + this.viewPosition.x);
              let targetY = (d.target.y + (d.targetIndex * 16 + 5) + 48) * this.zoom + this.viewPosition.y;
          context.lineTo(targetX-tableThing, targetY);
          context.lineTo(targetX, targetY);
          //tick
          context.moveTo(targetX-5, targetY-5);
          context.lineTo(targetX-5, targetY+5);
        } else {

          let baseX = (d.source.x * (this.zoom) + this.viewPosition.x);
          let baseY = (d.source.y + (d.sourceIndex * 16 + 5) + 48) * this.zoom + this.viewPosition.y;

          //arrow
          context.moveTo(baseX, baseY-arrowSize);
          context.lineTo(baseX-arrowSize, baseY);
          context.lineTo(baseX, baseY+arrowSize);

          context.moveTo(baseX, baseY);
          context.lineTo(baseX-tableThing, baseY);
          let targetX = ((d.target.x + d.targetWidth) * (this.zoom) + this.viewPosition.x);
          let targetY = (d.target.y + (d.targetIndex * 16 + 5) + 48) * this.zoom + this.viewPosition.y;
          context.lineTo(targetX+tableThing, targetY);
          context.lineTo(targetX, targetY);
          //tick
          context.moveTo(targetX+5, targetY-5);
          context.lineTo(targetX+5, targetY+5);
        }
      }
      context.strokeStyle = "#aaa";
      context.stroke();
      context.beginPath();
      context.fill();
      context.strokeStyle = "#fff";
      context.stroke();
      context.restore();
    },
    handleZoom: function (event) {
      if (event.deltaY > 0) {
        if (this.zoom < 0.05) {
          return;
        }
        this.zoom = this.zoom - 0.05;
        this.viewPosition.x = this.viewPosition.x - (event.clientX - this.viewPosition.x) * (-0.05);
        this.viewPosition.y = this.viewPosition.y - (event.clientY - this.viewPosition.y) * (-0.05);
      } else {
        if (this.zoom == 1) {
          return;
        }
        this.zoom = this.zoom + 0.05;
        this.viewPosition.x = this.viewPosition.x + (event.clientX - this.viewPosition.x) * (-0.05);
        this.viewPosition.y = this.viewPosition.y + (event.clientY - this.viewPosition.y) * (-0.05);
      }
      this.drawLines()
    },
    displayGraph: function () {
      this.overlay = true;
      let model = JSON.parse(JSON.stringify(this.database));
      let links = [];

      function getSourceTableData(tableName, columnName) {
        for (let i = 0; i < model.tables.length; i++) {
          if (tableName === model.tables[i].name) {
            for (let c = 0; c < model.tables[i].columns.length; c++) {
              if (model.tables[i].columns[c].name === columnName) {
                return [i, c, model.tables[i].x, model.tables[i].width, model.tables[i].columns[c].nullable];
              }
            }
          }
        }
        return [-1, -1, -1, -1, -1];
      }
      function getTargetTableData(tableName, columnName) {
        for (let i = 0; i < model.tables.length; i++) {
          if (tableName === model.tables[i].name) {
            for (let c = 0; c < model.tables[i].columns.length; c++) {
              if (model.tables[i].columns[c].name === columnName) {
                return [i, c, model.tables[i].x, model.tables[i].width, model.tables[i].columns[c].primary];
              }
            }
          }
        }
        return [-1, -1, -1, -1, -1];
      }

      function recalculateLinksX() {
        let index = 0;
        for (let i = 0; i < model.tables.length; i++) {
          for (let c = 0; c < model.tables[i].columns.length; c++) {
            for(let f = 0;f< model.tables[i].columns[c].fkConstraints.length;f++) {
              let column = model.tables[i].columns[c];
              let targetTable = getTargetTableData(column.fkConstraints[f].rTable, column.fkConstraints[f].rColumn);
              let sourceTable = getSourceTableData(model.tables[i].name, column.name);
              links[index].sourceX = sourceTable[2];
              links[index].targetX = targetTable[2];
              index++;
            }
          }
        }
      }

      function fillLinks() {
        for (let i = 0; i < model.tables.length; i++) {
          for (let c = 0; c < model.tables[i].columns.length; c++) {
            for (let f = 0; f < model.tables[i].columns[c].fkConstraints.length; f++) {
              let column = model.tables[i].columns[c];
              let targetTable = getTargetTableData(column.fkConstraints[f].rTable, column.fkConstraints[f].rColumn);
              let sourceTable = getSourceTableData(model.tables[i].name, column.name);
              links.push({
                source: i,
                sourceIndex: sourceTable[1],
                sourceX: sourceTable[2],
                sourceWidth: sourceTable[3],
                sourceOptional: sourceTable[4],
                target: targetTable[0],
                targetIndex: targetTable[1],
                targetX: targetTable[2],
                targetWidth: targetTable[3],
                targetIsId: targetTable[4],
              })
            }
          }
        }
      }

      fillLinks();

      function forceCollide() {
        let padding = 100;

        function force(alpha) {
          const quad = d3.quadtree(model.tables, d => d.x, d => d.y);
          for (let d of model.tables) {
            quad.visit((q, x1, y1, x2, y2) => {
              let updated = false;
              if (q.data && q.data !== d) {
                let x = d.x - q.data.x,
                    y = d.y - q.data.y,
                    xSpacing = padding + (q.data.width + d.width) / 2,
                    ySpacing = padding + (q.data.height + d.height) / 2,
                    absX = Math.abs(x),
                    absY = Math.abs(y),
                    l,
                    lx,
                    ly;

                if (absX < xSpacing && absY < ySpacing) {
                  l = Math.sqrt(x * x + y * y);

                  lx = (absX - xSpacing) / l;
                  ly = (absY - ySpacing) / l;

                  // the one that's barely within the bounds probably triggered the collision
                  if (Math.abs(lx) > Math.abs(ly)) {
                    lx = 0;
                  } else {
                    ly = 0;
                  }
                  d.x -= x *= lx;
                  d.y -= y *= ly;
                  q.data.x += x;
                  q.data.y += y;

                  updated = true;
                }
              }
              return updated;
            });
          }
        }

        force.initialize = _ => model.tables = _;

        return force;
      }

      this.simulation = d3.forceSimulation(model.tables)
          .force("link", d3.forceLink(links).distance(200).strength(.1).iterations(100))
          .force("x", d3.forceX().strength(.05))
          .force("y", d3.forceY().strength(.05))
          .force("center", d3.forceCenter(this.$refs.canvasScreen.clientWidth / 2, this.$refs.canvasScreen.clientHeight / 2).strength(.1))
          .force("collide", forceCollide());
      // this.simulation.on('tick', () => {
      //   recalculateLinksX();
      //   this.links = links;
      //   this.database = model;
      //   this.drawLines();
      //   this.overlay = false;
      // });
      this.simulation.on('end', ()=>{
        for (let table of model.tables) {
          table.x = table.x - table.width / 2;
          table.y = table.y - table.height / 2;
        }
        recalculateLinksX();
        this.links = links;
        this.database = model;
        this.drawLines();
        this.overlay = false;
      });
    },
  },
  mounted() {
    this.getDatabaseModel(this.$route.params.dbName, this.$route.params.schema);
  },
  beforeRouteLeave: function (to, from, next) {
    if (this.simulation) {
      this.simulation.stop();
    }
    next();
  },
  watch: {
    databaseModel: function () {
      this.database = this.databaseModel
      this.displayGraph();
    },
    $route() {
      if (this.simulation) {
        this.simulation.stop();
      }
      this.clearLines();
      this.getDatabaseModel(this.$route.params.dbName, this.$route.params.schema);
    }
  }
}
</script>

<style>

.table {
  padding: 5px 5px;
  border: 1px solid #454545 !important;
  border-radius: 5px;
  position: fixed;
}

.table-item {
  width: max-content;
}

#canvasScreen {
  z-index: 1000;
  overflow: hidden;
}

#canvas {
  width: 100%;
  height: calc(100vh - 64px);
  z-index: 0;
}

.zoom-level {
  position: absolute;
  right: -5px;
  top: -5px;
  padding: 15px 15px 10px 10px;
  background-color: #c7c7c7;
  border-radius: 5px;
  border: 1px solid #aaaaaa;
  box-shadow: 0px 3px 1px -2px rgb(0 0 0 / 20%), 0px 2px 2px 0px rgb(0 0 0 / 14%), 0px 1px 5px 0px rgb(0 0 0 / 13%);
}
</style>
