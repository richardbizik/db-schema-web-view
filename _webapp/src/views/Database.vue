<template>
  <div>
    <v-overlay :value="overlay">
      <v-progress-circular
          indeterminate
          size="64"
      ></v-progress-circular>
      <span class="overlay-text">Placing tables...</span>
    </v-overlay>
    <canvas ref="canvasElement"
            :height="this.$refs.canvasScreen?this.$refs.canvasScreen.clientHeight:0"
            :width="this.$refs.canvasScreen?this.$refs.canvasScreen.clientWidth:0"
            style="z-index: 2;position: absolute;pointer-events: none;"></canvas>
    <div id="canvasScreen" ref="canvasScreen" @mousedown="startMove" @mouseup="cancelMove" @mousemove="moveCanvas"
         @click.middle="displayGraph"
         @wheel.prevent="handleZoom"
         :style="isMousePressed?'cursor:grabbing;':'cursor:grab'"
    >
      <div class="zoom-level">{{ Math.round((zoom / 100) * 10000) }}%</div>
      <div id="canvas" :style="getCanvasStyle">
        <v-list-item
            class="table-item"
            :style="isMousePressed?'user-select:none;':''"
            v-for="(table) in database.tables"
            :key="table.name"
            ref="list"
        >
          <Table :ref="table.name" class="table" :table-data="table" v-on:table-entered="tableEntered"
                 v-on:table-left="tableLeft"
                 v-on:table-initialized="tableInitialized"/>

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
import Vue from "vue";

export default {
  name: 'Database',
  components: {Table},
  data: () => ({
    viewPosition: {x: 0, y: 0},
    zoom: 1,
    dragStart: {x: 0, y: 0, originalViewPosition: {x: 0, y: 0}},
    isMousePressed: false,
    simulation: null,
    database: {tables: []},
    links: [],
    lines: [],
    overlay: false,
    context: null,
    initializedTables: [],
    simulationRunning: false,
  }),
  computed: {
    ...mapState('database', ['databases', 'databaseModel']),
    getCanvasStyle() {
      return 'transform-origin: 0px 0px 0px; transform: matrix(' + this.zoom + ', 0, 0, ' + this.zoom + ', ' + this.viewPosition.x + ', ' + this.viewPosition.y + ';';
    }
  },
  methods: {
    init: function () {
      this.context = this.$refs.canvasElement.getContext("2d");
      this.getDatabaseModel(this.$route.params.dbName, this.$route.params.schema);
    },
    getDatabases: function () {
      this.$store.dispatch('database/getDatabases');
    },
    getDatabaseModel: function (dbName, schema) {
      this.$store.dispatch('database/getDatabaseModel', {dbName, schema});
    },
    startMove: function (event) {
      if (event.composedPath()[0].id === "canvasScreen" || event.composedPath()[0].id === "canvas") {
        this.isMousePressed = true;
        if (this.isMousePressed) {
          this.dragStart.x = event.x
          this.dragStart.y = event.y
          this.dragStart.originalViewPosition.x = this.viewPosition.x
          this.dragStart.originalViewPosition.y = this.viewPosition.y
        }
      }
    },
    tableEntered: function (tableName) {
      this.highlightLines(tableName);
    },
    tableLeft: function (tableName) {
      this.removeHighlightedLines(tableName);
    },
    tableInitialized: function (tableName) {
      for (let i = 0; i < this.initializedTables.length; i++) {
        if (this.initializedTables[i].table === tableName) {
          Vue.set(this.initializedTables, i, {table: tableName, ready: true});
          return;
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
      } else if (event.composedPath()[0].id === "canvasScreen" || event.composedPath()[0].id === "canvas") {
        let isHovering = false;
        for (let lineObj of this.lines) {
          for (let linePart of lineObj.lineParts) {
            this.context.restore();
            this.context.save();
            this.context.beginPath();
            this.context.moveTo(linePart.sx, linePart.sy);
            this.context.lineTo(linePart.tx, linePart.ty);
            this.context.lineWidth = 10 * this.zoom;
            if (this.context.isPointInStroke(event.layerX, event.layerY)) {
              this.context.lineWidth = 1;
              this.drawLines();
              this.context.beginPath();
              this.context.lineWidth = 1;
              this.highlightLine(lineObj);
              this.context.strokeStyle = "#aaaaaa";
              this.context.stroke();
              this.context.beginPath();
              this.context.fill();
              this.context.strokeStyle = "#fff";
              this.context.stroke();
              isHovering = true;
              for (let i = 0; i < this.database.tables.length; i++) {
                if (this.database.tables[i].name === lineObj.base || this.database.tables[i].name === lineObj.target) {
                  this.database.tables[i].hover = true;
                }
              }
            }
          }
        }
        this.context.restore();
        this.context.lineWidth = 1;
        if (isHovering == false) {
          this.drawLines();
          for (let table of this.database.tables) {
            table.hover = false;
          }
        }
      }
    },
    clearLines: function () {
      let width = this.$refs.canvasScreen.clientWidth;
      let height = this.$refs.canvasScreen.clientHeight;
      this.context.clearRect(0, 0, width, height);
    },
    highlightLine: function (line) {
      for (let linePart of line.lineParts) {
        this.context.moveTo(linePart.sx, linePart.sy);
        this.context.lineTo(linePart.tx, linePart.ty);
        this.context.shadowOffsetX = 2 * this.zoom;
        this.context.shadowOffsetY = 2 * this.zoom;
        this.context.shadowBlur = 5 * this.zoom;
        this.context.shadowColor = "#000000";
      }
    },
    highlightLines: function (tableName) {
      this.context.save();
      this.context.beginPath();
      for (let lineObj of this.lines) {
        if (lineObj.base === tableName || lineObj.target === tableName) {
          this.highlightLine(lineObj);
        }
      }
      this.context.strokeStyle = "#aaaaaa";
      this.context.stroke();
      this.context.beginPath();
      this.context.fill();
      this.context.strokeStyle = "#fff";
      this.context.stroke();
      this.context.restore();
    },
    removeHighlightedLines: function (tableName) {
      this.drawLines();
    },
    drawLines: function () {
      if (!this.$refs.canvasScreen) {
        return;
      }
      let width = this.$refs.canvasScreen.clientWidth;
      let height = this.$refs.canvasScreen.clientHeight;
      let tickSize = 15 * this.zoom;
      let arrowSize = 5 * this.zoom;
      let circleSize = 3 * this.zoom;
      let lines = [];
      this.context.clearRect(0, 0, width, height);
      this.context.save();
      this.context.beginPath();
      for (const d of this.links) {
        if (Math.round(d.sourceX) < Math.round(d.targetX)) {
          let baseX = ((d.source.x + d.sourceWidth) * (this.zoom) + this.viewPosition.x);
          let baseY = (d.source.y + (d.sourceIndex * 16 + 5) + 48) * this.zoom + this.viewPosition.y;

          //arrow
          this.context.moveTo(baseX, baseY - arrowSize);
          this.context.lineTo(baseX + arrowSize, baseY);
          this.context.lineTo(baseX, baseY + arrowSize);
          //circle
          if (d.sourceOptional) {
            this.context.moveTo(baseX + arrowSize + circleSize, baseY);
            this.context.arc(baseX + arrowSize + circleSize, baseY, circleSize, 0, 2 * Math.PI);
          }

          this.context.moveTo(baseX, baseY);
          this.context.lineTo(baseX + tickSize, baseY);
          let targetX = (d.target.x * (this.zoom) + this.viewPosition.x);
          let targetY = (d.target.y + (d.targetIndex * 16 + 5) + 48) * this.zoom + this.viewPosition.y;
          this.context.lineTo(targetX - tickSize, targetY);
          this.context.lineTo(targetX, targetY);
          //tick
          this.context.moveTo(targetX - tickSize / 3, targetY - tickSize / 3);
          this.context.lineTo(targetX - tickSize / 3, targetY + tickSize / 3);
          lines.push({
            lineParts: [
              {sx: baseX, sy: baseY, tx: baseX + tickSize, ty: baseY},
              {sx: baseX + tickSize, sy: baseY, tx: targetX - tickSize, ty: targetY},
              {sx: targetX - tickSize, sy: targetY, tx: targetX, ty: targetY},
            ], base: d.sourceTable, target: d.targetTable
          });
        } else {

          let baseX = (d.source.x * (this.zoom) + this.viewPosition.x);
          let baseY = (d.source.y + (d.sourceIndex * 16 + 5) + 48) * this.zoom + this.viewPosition.y;

          //arrow
          this.context.moveTo(baseX, baseY - arrowSize);
          this.context.lineTo(baseX - arrowSize, baseY);
          this.context.lineTo(baseX, baseY + arrowSize);
          //circle
          if (d.sourceOptional) {
            this.context.moveTo(baseX - (arrowSize + circleSize), baseY);
            this.context.arc(baseX - (arrowSize + circleSize), baseY, circleSize, 0, 2 * Math.PI);
          }

          this.context.moveTo(baseX, baseY);
          this.context.lineTo(baseX - tickSize, baseY);
          let targetX = ((d.target.x + d.targetWidth) * (this.zoom) + this.viewPosition.x);
          let targetY = (d.target.y + (d.targetIndex * 16 + 5) + 48) * this.zoom + this.viewPosition.y;
          this.context.lineTo(targetX + tickSize, targetY);
          this.context.lineTo(targetX, targetY);
          //tick
          this.context.moveTo(targetX + tickSize / 3, targetY - tickSize / 3);
          this.context.lineTo(targetX + tickSize / 3, targetY + tickSize / 3);
          lines.push({
            lineParts: [
              {sx: baseX, sy: baseY, tx: baseX - tickSize, ty: baseY},
              {sx: baseX - tickSize, sy: baseY, tx: targetX + tickSize, ty: targetY},
              {sx: targetX + tickSize, sy: targetY, tx: targetX, ty: targetY},
            ], base: d.sourceTable, target: d.targetTable
          });
        }
      }
      this.context.strokeStyle = "#aaa";
      this.context.stroke();
      this.context.beginPath();
      this.context.fill();
      this.context.strokeStyle = "#fff";
      this.context.stroke();
      this.context.restore();
      this.lines = lines;
    },
    handleZoom: function (event) {
      let boundingClientRect = this.$refs.canvasScreen.getBoundingClientRect();
      if (event.deltaY > 0) {
        if (this.zoom < 0.05) {
          return;
        }
        this.zoom = this.zoom - 0.05;
        this.viewPosition.x = this.viewPosition.x - (event.clientX - boundingClientRect.x - this.viewPosition.x) * -0.05;
        this.viewPosition.y = this.viewPosition.y - (event.clientY - boundingClientRect.y - this.viewPosition.y) * -0.05;
      } else {
        if (this.zoom > 1.96) {
          return;
        }
        this.zoom = this.zoom + 0.05;
        this.viewPosition.x = this.viewPosition.x + (event.clientX - boundingClientRect.x - this.viewPosition.x) * -0.05;
        this.viewPosition.y = this.viewPosition.y + (event.clientY - boundingClientRect.y - this.viewPosition.y) * -0.05;
      }
      this.drawLines();
    },
    displayGraph: function () {
      this.simulationRunning = true;
      this.overlay = true;
      let model = JSON.parse(JSON.stringify(this.database));
      let links = [];
      let refs = this.$refs;

      function getSourceTableData(tableName, columnName) {
        for (let i = 0; i < model.tables.length; i++) {
          if (tableName === model.tables[i].name) {
            let columnOffset = 0;
            for (let c = 0; c < model.tables[i].columns.length; c++) {
              let tableRef = refs[model.tables[i].name];
              if (tableRef !== undefined) {
                let col = tableRef[0].$refs[model.tables[i].columns[c].name];
                if (col !== undefined) {
                  columnOffset = columnOffset + ((col[0].clientHeight / 16) - 1);
                }
              }
              if (model.tables[i].columns[c].name === columnName) {
                return [i, (c + columnOffset), model.tables[i].x, model.tables[i].width, model.tables[i].columns[c].nullable];
              }
            }
          }
        }
        return [-1, -1, -1, -1, -1];
      }

      function getTargetTableData(tableName, columnName) {
        for (let i = 0; i < model.tables.length; i++) {
          if (tableName === model.tables[i].name) {
            let columnOffset = 0;
            for (let c = 0; c < model.tables[i].columns.length; c++) {
              let tableRef = refs[model.tables[i].name];
              if (tableRef !== undefined) {
                let col = tableRef[0].$refs[model.tables[i].columns[c].name];
                if (col !== undefined) {
                  columnOffset = columnOffset + ((col[0].clientHeight / 16) - 1);
                }
              }
              if (model.tables[i].columns[c].name === columnName) {
                return [i, c + columnOffset, model.tables[i].x, model.tables[i].width, model.tables[i].columns[c].nullable];
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
            for (let f = 0; f < model.tables[i].columns[c].fkConstraints.length; f++) {
              let column = model.tables[i].columns[c];
              let targetTable = getTargetTableData(column.fkConstraints[f].rTable, column.fkConstraints[f].rColumn);
              let sourceTable = getSourceTableData(model.tables[i].name, column.name);
              links[index].sourceIndex = sourceTable[1];
              links[index].sourceX = sourceTable[2];
              links[index].targetIndex = targetTable[1];
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
                sourceTable: model.tables[i].name,
                sourceX: sourceTable[2],
                sourceWidth: sourceTable[3],
                sourceOptional: sourceTable[4],
                target: targetTable[0],
                targetTable: column.fkConstraints[f].rTable,
                targetX: targetTable[2],
                targetWidth: targetTable[3],
                targetOptional: targetTable[4],
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
      this.simulation.on('end', () => {
        this.simulationRunning = false;
        if (model.name !== this.database.name) {
          return;
        }
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
    selectTable: function (tableName) {
      if (!this.$refs.canvasScreen) {
        return;
      }
      let screenWidth = this.$refs.canvasScreen.clientWidth;
      let screenHeight = this.$refs.canvasScreen.clientHeight;
      let table;
      for (let t of this.database.tables) {
        if (t.name === tableName) {
          table = t;
        }
      }
      if (table !== undefined) {
        this.viewPosition.x = screenWidth / 2 - (table.x + table.width / 2) * this.zoom;
        this.viewPosition.y = screenHeight / 2 - (table.y + table.height / 2) * this.zoom;
        this.drawLines();
      }
    },
    resetDatabase(){
      this.overlay = false;
      this.database.tables = [];
      this.zoom = 1;
      this.viewPosition.x = 0;
      this.viewPosition.y = 0;
      this.clearLines();
      this.getDatabaseModel(this.$route.params.dbName, this.$route.params.schema);
    }
  },
  mounted() {
    this.init();
  },
  created() {
    this.$eventHub.$on("select-table", this.selectTable);
  },
  beforeDestroy() {
    this.$eventHub.$off('select-table');
  },
  beforeRouteLeave: function (to, from, next) {
    if (this.simulation) {
      this.simulation.stop();
    }
    next();
  },
  watch: {
    databaseModel: function (newValue) {
      for (let table of newValue.tables) {
        this.initializedTables.push({table: table.name, ready: false});
      }
      this.database = newValue;
    },
    initializedTables: function (newValue) {
      if (newValue.length !== 0) {
        for (let table of newValue) {
          if (!table.ready) {
            return;
          }
        }
        this.displayGraph();
      }
    },
    $route() {
      if (this.database.name != this.$route.params.dbName || this.database.schema != this.$route.params.schema) {
        this.initializedTables = [];
        if (this.simulation) {
          this.simulation.stop();
        }
        this.resetDatabase();
      }
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
  z-index: 999;
  padding: 15px 15px 10px 10px;
  background-color: #363636;
  border-radius: 5px;
  border: 1px solid #1e1e1e;
  box-shadow: 0px 3px 1px -2px rgb(0 0 0 / 20%), 0px 2px 2px 0px rgb(0 0 0 / 14%), 0px 1px 5px 0px rgb(0 0 0 / 13%);
}
.overlay-text{
  font-size: 2rem !important;
  padding-left:20px;
}
</style>
