webpackJsonp([1],{0:function(t,e){},"2EO6":function(t,e){},NHnr:function(t,e,a){"use strict";Object.defineProperty(e,"__esModule",{value:!0});var s=a("/5sW"),r=a("8+8L"),n={name:"App",data:function(){return{haproxies:[],servers:[],databases:[]}},created:function(){this.loadData(),setInterval(function(){this.loadData()}.bind(this),3e3)},methods:{loadData:function(){this.$http.get("/api/dashboard").then(function(t){this.haproxies=t.data?t.data.haproxies:[],this.servers=t.data?t.data.servers:[],this.databases=t.data?t.data.databases:[]})},startServer:function(t){this.$http.post("/api/servers?name="+t).then(function(t){alert("Operation succeed, please wait page to reload data!")},function(t){alert(t.data)})},shutdownServer:function(t){this.$http.delete("/api/servers?name="+t).then(function(t){alert("Operation succeed, please wait page to reload data!")},function(t){alert(t.data)})}}},i={render:function(){var t=this,e=t.$createElement,a=t._self._c||e;return a("div",{attrs:{id:"app"}},[t._m(0),t._v(" "),a("h3",[t._v("HAProxy")]),t._v(" "),a("table",{attrs:{border:"1"}},[t._m(1),t._v(" "),a("tbody",t._l(t.haproxies,function(e){return a("tr",{key:e.name},[a("td",[t._v(t._s(e.name))]),t._v(" "),a("td",[a("a",{attrs:{target:"_blank",href:"http://"+e.address}},[t._v(t._s(e.address))])]),t._v(" "),a("td",[t._v(t._s(e.status))])])}))]),t._v(" "),a("h3",[t._v("Application")]),t._v(" "),a("table",{attrs:{border:"1"}},[t._m(2),t._v(" "),a("tbody",t._l(t.servers,function(e){return a("tr",{key:e.name},[a("td",[t._v(t._s(e.name))]),t._v(" "),a("td",[a("a",{attrs:{target:"_blank",href:"http://"+e.address}},[t._v(t._s(e.address))])]),t._v(" "),a("td",[t._v(t._s(e.status))]),t._v(" "),a("td",["running"==e.status?a("a",{attrs:{href:"#"},on:{click:function(a){t.shutdownServer(e.name)}}},[t._v("Shutdown")]):a("a",{attrs:{href:"#"},on:{click:function(a){t.startServer(e.name)}}},[t._v("Start")])])])}))]),t._v(" "),a("h3",[t._v("Database")]),t._v(" "),a("table",{attrs:{border:"1"}},[t._m(3),t._v(" "),a("tbody",t._l(t.databases,function(e){return a("tr",{key:e.name},[a("td",[t._v(t._s(e.name))]),t._v(" "),a("td",[t._v(t._s(e.address))]),t._v(" "),a("td",[t._v(t._s(e.status))])])}))])])},staticRenderFns:[function(){var t=this.$createElement,e=this._self._c||t;return e("h1",[this._v("Control Panel "),e("sub",[this._v("Project Spartan")])])},function(){var t=this.$createElement,e=this._self._c||t;return e("thead",[e("tr",[e("th",[this._v("Name")]),this._v(" "),e("th",[this._v("Address")]),this._v(" "),e("th",[this._v("Status")])])])},function(){var t=this.$createElement,e=this._self._c||t;return e("thead",[e("tr",[e("th",[this._v("Name")]),this._v(" "),e("th",[this._v("Address")]),this._v(" "),e("th",[this._v("Status")]),this._v(" "),e("th",[this._v("Action")])])])},function(){var t=this.$createElement,e=this._self._c||t;return e("thead",[e("tr",[e("th",[this._v("Name")]),this._v(" "),e("th",[this._v("Address")]),this._v(" "),e("th",[this._v("Status")])])])}]};var _=a("VU/8")(n,i,!1,function(t){a("2EO6")},null,null).exports;s.a.config.productionTip=!1,s.a.use(r.a),new s.a({el:"#app",render:function(t){return t(_)}})}},["NHnr"]);