<template>
<div class="search" @mousedown="startDrag" @touchstart="startDrag" @mousemove="onDrag" @touchmove="onDrag" @mouseup="stopDrag" @touchend="stopDrag" @mouseleave="stopDrag">
  <svg class="bg" width="640" height="560">
      <path :d="headerPath" fill="#F5F5F5"></path>
    </svg>
  <div class="header">
    <slot name="header"></slot>
  </div>
  <div class="content" :style="contentPosition">
    <slot name="content"></slot>
  </div>
</div>
</template>

<script>
export default {
  name: 'search',
  data: function () {
    return {
      dragging: false,
      // quadratic bezier control point
      c: {
        x: 160,
        y: 160
      },
      // record drag start point
      start: {
        x: 0,
        y: 0
      }
    }
  },
  computed: {
    headerPath: function () {
      return 'M0,0 L640,0 640,160' +
        'Q' + this.c.x + ',' + this.c.y +
        ' 0,160'
    },
    contentPosition: function () {
      var dy = this.c.y - 160
      var dampen = dy > 0 ? 2 : 4
      return {
        transform: 'translate3d(0,' + dy / dampen + 'px,0)'
      }
    }
  },
  methods: {
    startDrag: function (e) {
      e = e.changedTouches ? e.changedTouches[0] : e
      this.dragging = true
      this.start.x = e.pageX
      this.start.y = e.pageY
    },
    onDrag: function (e) {
      e = e.changedTouches ? e.changedTouches[0] : e
      if (this.dragging) {
        this.c.x = 320 + (e.pageX - this.start.x)
          // dampen vertical drag by a factor
        var dy = e.pageY - this.start.y
        var dampen = dy > 0 ? 1.5 : 4
        this.c.y = 160 + dy / dampen
      }
    },
    stopDrag: function () {
      var dynamics = require('dynamics.js')
      if (this.dragging) {
        this.dragging = false
        dynamics.animate(this.c, {
          x: 160,
          y: 160
        }, {
          type: dynamics.spring,
          duration: 700,
          friction: 280
        })
      }
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
h1 {
  font-weight: 300;
  font-size: 1.8em;
  margin-top: 0;
  background: #F5F5F5;
}

.search {
  background: #F5F5F5;
  box-shadow: 0 4px 16px rgba(0, 0, 0, .15);
  width: 640px;
  height: 560px;
  overflow: hidden;
  margin: 30px auto;
  position: relative;
  font-family: 'Roboto', Helvetica, Arial, sans-serif;
  color: #fff;
}

.search .bg {
  position: absolute;
  top: 0;
  left: 0;
  z-index: 0;
}

.search .header {
  position: relative;
  padding: 30px;
}
.search .content {
  padding-left: 20px;
  padding-right: 20px;
}
.search .header {
  height: 160px;
}
</style>
