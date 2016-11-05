<template>
  <input class="searchInput"
    autofocus autocomplete="off"
    placeholder="What needs to be done?"
    v-model="searchStr"
    @keyup.enter="search">
</template>

<script>
export default {
  name: 'searchInput',
  data: function () {
    return {
      searchStr: ''
    }
  },

  methods: {
    search: function () {
      this.$http.get(
        '/search',
        {
          searchStr: this.searchStr
        }
      ).then((response) => {
        this.alert('a')
        this.alert(response.data)
      }, (response) => {
        console.log(this.searchStr)
      })
      this.searchStr = ''
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.searchInput {
    padding: 16px 16px 16px 60px;
    border: none;
    background: #fff;
    box-shadow: inset 0 -2px 1px rgba(0,0,0,0.03);
}

.searchInput input::-webkit-input-placeholder {
	font-style: italic;
	font-weight: 300;
	color: #e6e6e6;
}

.searchInput input::-moz-placeholder {
	font-style: italic;
	font-weight: 300;
	color: #e6e6e6;
}

.searchInput input::input-placeholder {
	font-style: italic;
	font-weight: 300;
	color: #e6e6e6;
}

.searchInput, .edit {
    position: relative;
    /*margin: 0;*/
    width: 100%;
    font-size: 24px;
    font-family: inherit;
    font-weight: inherit;
    line-height: 1.4em;
    border: 0;
    outline:none;
    font-weight: 300;
    font-style: italic;
    /*padding: 6px;*/
    /*border: 1px solid #999;*/
    box-shadow: inset 0 -2px 1px rgba(0,0,0,0.03);
    box-sizing: border-box;
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
}
</style>
