(function(e){function t(t){for(var r,a,i=t[0],u=t[1],l=t[2],p=0,d=[];p<i.length;p++)a=i[p],o[a]&&d.push(o[a][0]),o[a]=0;for(r in u)Object.prototype.hasOwnProperty.call(u,r)&&(e[r]=u[r]);c&&c(t);while(d.length)d.shift()();return s.push.apply(s,l||[]),n()}function n(){for(var e,t=0;t<s.length;t++){for(var n=s[t],r=!0,i=1;i<n.length;i++){var u=n[i];0!==o[u]&&(r=!1)}r&&(s.splice(t--,1),e=a(a.s=n[0]))}return e}var r={},o={app:0},s=[];function a(t){if(r[t])return r[t].exports;var n=r[t]={i:t,l:!1,exports:{}};return e[t].call(n.exports,n,n.exports,a),n.l=!0,n.exports}a.m=e,a.c=r,a.d=function(e,t,n){a.o(e,t)||Object.defineProperty(e,t,{enumerable:!0,get:n})},a.r=function(e){"undefined"!==typeof Symbol&&Symbol.toStringTag&&Object.defineProperty(e,Symbol.toStringTag,{value:"Module"}),Object.defineProperty(e,"__esModule",{value:!0})},a.t=function(e,t){if(1&t&&(e=a(e)),8&t)return e;if(4&t&&"object"===typeof e&&e&&e.__esModule)return e;var n=Object.create(null);if(a.r(n),Object.defineProperty(n,"default",{enumerable:!0,value:e}),2&t&&"string"!=typeof e)for(var r in e)a.d(n,r,function(t){return e[t]}.bind(null,r));return n},a.n=function(e){var t=e&&e.__esModule?function(){return e["default"]}:function(){return e};return a.d(t,"a",t),t},a.o=function(e,t){return Object.prototype.hasOwnProperty.call(e,t)},a.p="/";var i=window["webpackJsonp"]=window["webpackJsonp"]||[],u=i.push.bind(i);i.push=t,i=i.slice();for(var l=0;l<i.length;l++)t(i[l]);var c=u;s.push([0,"chunk-vendors"]),n()})({0:function(e,t,n){e.exports=n("56d7")},"56d7":function(e,t,n){"use strict";n.r(t);n("cadf"),n("551c"),n("f751"),n("097d");var r=n("2b0e"),o=function(){var e=this,t=e.$createElement,n=e._self._c||t;return n("div",{attrs:{id:"app"}},[n("router-view")],1)},s=[],a=(n("7faf"),n("2877")),i={},u=Object(a["a"])(i,o,s,!1,null,null,null),l=u.exports,c=n("8c4f"),p=function(){var e=this,t=e.$createElement,n=e._self._c||t;return n("div",{staticClass:"home"},[n("div",{staticClass:"tools-bar"},[n("div",{staticStyle:{flex:"1"}},[e._v("\n      Filter:\n      "),n("input",{directives:[{name:"model",rawName:"v-model",value:e.keyword,expression:"keyword"}],staticClass:"search",attrs:{placeholder:"输入过滤内容",type:"search"},domProps:{value:e.keyword},on:{input:function(t){t.target.composing||(e.keyword=t.target.value)}}})]),n("div",[e._v("\n      展示最近：\n      "),n("input",{directives:[{name:"model",rawName:"v-model",value:e.maxCount,expression:"maxCount"}],staticClass:"input",attrs:{type:"number",step:"1",min:"1"},domProps:{value:e.maxCount},on:{input:function(t){t.target.composing||(e.maxCount=t.target.value)}}}),e._v("  条\n    ")])]),n("table",{staticClass:"table"},[n("thead",[n("tr",[e._l(e.columns,function(t,r){return n("th",{key:r,on:{click:function(n){return e.toggleSort(t)}}},[e._v("\n          "+e._s(t.label)+"\n          "),e.sortColumn===t?n("span",[e._v("\n            "+e._s("asc"===e.orders?"▲":"▼")+"\n          ")]):e._e()])}),n("th",[e._v("source")])],2)]),n("transition-group",{attrs:{tag:"tbody",name:"list"}},e._l(e.filteredList.length?e.filteredList:e.sortedList,function(t){return n("tr",{key:t["@timestamp"]},[e._l(e.columns,function(r,o){return n("td",{key:o},[e._v("\n          "+e._s(e.pathValue(r.path,t))+"\n        ")])}),n("td",[n("ul",e._l(t.event.properties,function(t,r){return n("div",{key:r,staticClass:"prop-item"},[n("label",[e._v(e._s(r))]),e._v(": "+e._s(t))])}),0)])],2)}),0)],1)])},d=[],f=n("a4bb"),v=n.n(f),m=n("f499"),h=n.n(m),b=n("2ef0"),y=n.n(b),g={name:"home",data:function(){return{messages:[],keyword:"",sortColumn:null,orders:"desc",maxCount:2e3,columns:window["columns"]||[{label:"Time",path:"event.time"},{label:"distinct_id",path:"event.distinct_id"},{label:"type",path:"event.type"},{label:"event",path:"event.event"}]}},created:function(){var e=("https:"===window.location.protocol?"wss://":"ws://")+window.location.host+"/topic",t=new WebSocket(e);t.addEventListener("open",function(e){console.log("Hello Server!")}),t.addEventListener("error",function(e){alert("socket 连接失败")}),t.addEventListener("message",this.onMessage)},computed:{sortedList:function(){return this.sortColumn?y.a.orderBy(this.messages,this.sortColumn.path,this.orders):this.messages},filteredList:function(){var e=this;return this.keyword?this.messages.filter(function(t){return h()(t).indexOf(e.keyword)>-1}):[]}},methods:{onMessage:function(e){var t=JSON.parse(e.data);t.event&&t.event.properties&&(t.event.properties=this.sortByKeys(t.event.properties)),this.messages.unshift(t),this.messages=this.messages.slice(0,this.maxCount)},sortByKeys:function(e){var t=v()(e),n=y.a.sortBy(t);return y.a.fromPairs(y.a.map(n,function(t){return[t,e[t]]}))},toggleSort:function(e){this.sortColumn=e,"desc"===this.orders?this.orders="asc":this.orders="desc"},pathValue:function(e,t){return y.a.get(t,e,"")}}},w=g,_=(n("c219"),Object(a["a"])(w,p,d,!1,null,null,null)),x=_.exports;r["a"].use(c["a"]);var k=new c["a"]({mode:"history",base:"/",routes:[{path:"/",name:"home",component:x}]}),C=n("2f62");r["a"].use(C["a"]);var O=new C["a"].Store({state:{},mutations:{},actions:{}});r["a"].config.productionTip=!1,new r["a"]({router:k,store:O,render:function(e){return e(l)}}).$mount("#app")},"7faf":function(e,t,n){"use strict";var r=n("8fba"),o=n.n(r);o.a},"8fba":function(e,t,n){},c219:function(e,t,n){"use strict";var r=n("e9bb"),o=n.n(r);o.a},e9bb:function(e,t,n){}});
//# sourceMappingURL=app.13736971.js.map