(window.webpackJsonp=window.webpackJsonp||[]).push([[75],{177:function(e,t,r){"use strict";r.r(t),r.d(t,"frontMatter",(function(){return i})),r.d(t,"metadata",(function(){return o})),r.d(t,"rightToc",(function(){return p})),r.d(t,"default",(function(){return s}));var n=r(1),a=r(6),c=(r(0),r(193)),i={title:"Object"},o={id:"reference/opspec/types/object",title:"Object",description:"Object typed values are a container for string indexed values (referred to as properties).",source:"@site/docs/reference/opspec/types/object.md",permalink:"/docs/reference/opspec/types/object",editUrl:"https://github.com/opctl/opctl/edit/master/docs/docs/reference/opspec/types/object.md",lastUpdatedBy:"Chris Dostert",lastUpdatedAt:1587672399,sidebar:"docs",previous:{title:"Number",permalink:"/docs/reference/opspec/types/number"},next:{title:"Socket",permalink:"/docs/reference/opspec/types/socket"}},p=[{value:"Initialization",id:"initialization",children:[]},{value:"Property Referencing",id:"property-referencing",children:[]},{value:"Coercion",id:"coercion",children:[]}],l={rightToc:p},b="wrapper";function s(e){var t=e.components,r=Object(a.a)(e,["components"]);return Object(c.b)(b,Object(n.a)({},l,r,{components:t,mdxType:"MDXLayout"}),Object(c.b)("p",null,"Object typed values are a container for string indexed values (referred to as properties)."),Object(c.b)("p",null,"Objects..."),Object(c.b)("ul",null,Object(c.b)("li",{parentName:"ul"},"are immutable, i.e. assigning to an object results in a copy of the original object"),Object(c.b)("li",{parentName:"ul"},"can be passed in/out of ops via ",Object(c.b)("a",Object(n.a)({parentName:"li"},{href:"/docs/reference/opspec/op-directory/op/parameter/object"}),"object parameters")),Object(c.b)("li",{parentName:"ul"},"can be initialized via ",Object(c.b)("a",Object(n.a)({parentName:"li"},{href:"#initialization"}),"object initialization")),Object(c.b)("li",{parentName:"ul"},"properties can be referenced via ",Object(c.b)("a",Object(n.a)({parentName:"li"},{href:"#property-referencing"}),"object property referencing")),Object(c.b)("li",{parentName:"ul"},"are coerced according to ",Object(c.b)("a",Object(n.a)({parentName:"li"},{href:"#coercion"}),"object coercion"))),Object(c.b)("h3",{id:"initialization"},"Initialization"),Object(c.b)("p",null,"Object typed values can be constructed from a literal or templated object."),Object(c.b)("p",null,"A templated object is an object which includes one or more value reference.\nAt runtime, each reference gets evaluated and replaced with it's corresponding value."),Object(c.b)("h4",{id:"initialization-example-literal"},"Initialization Example (literal)"),Object(c.b)("pre",null,Object(c.b)("code",Object(n.a)({parentName:"pre"},{className:"language-yaml"}),"myObject:\n    prop1: value\n")),Object(c.b)("h4",{id:"initialization-example-templated"},"Initialization Example (templated)"),Object(c.b)("p",null,"given:"),Object(c.b)("ul",null,Object(c.b)("li",{parentName:"ul"},Object(c.b)("inlineCode",{parentName:"li"},"/someDir/file2.txt")," is embedded in op"),Object(c.b)("li",{parentName:"ul"},Object(c.b)("inlineCode",{parentName:"li"},"prop2Name")," is in scope"),Object(c.b)("li",{parentName:"ul"},Object(c.b)("inlineCode",{parentName:"li"},"someObject"),Object(c.b)("ul",{parentName:"li"},Object(c.b)("li",{parentName:"ul"},"is in scope"),Object(c.b)("li",{parentName:"ul"},"is type coercible to object"),Object(c.b)("li",{parentName:"ul"},"has property ",Object(c.b)("inlineCode",{parentName:"li"},"someProperty")))),Object(c.b)("li",{parentName:"ul"},Object(c.b)("inlineCode",{parentName:"li"},"prop4")," is in scope")),Object(c.b)("pre",null,Object(c.b)("code",Object(n.a)({parentName:"pre"},{className:"language-yaml"}),"# interpolate properties\nmyObject:\n    prop1: string $(/someDir/file2.txt)\n    $(prop2Name): $(someObject.someProperty)\n    prop3: [ sub, array, 2]\n    # Shorthand property name; equivalent to prop4: $(prop4)\n    prop4:\n")),Object(c.b)("h3",{id:"property-referencing"},"Property Referencing"),Object(c.b)("p",null,"Object properties can be referenced via ",Object(c.b)("inlineCode",{parentName:"p"},"$(OBJECT.PROPERTY)")," or ",Object(c.b)("inlineCode",{parentName:"p"},"$(OBJECT[PROPERTY])")," syntax."),Object(c.b)("h4",{id:"property-referencing-example-from-scope"},"Property Referencing Example (from scope)"),Object(c.b)("p",null,"given:"),Object(c.b)("ul",null,Object(c.b)("li",{parentName:"ul"},Object(c.b)("inlineCode",{parentName:"li"},"someObject"),Object(c.b)("ul",{parentName:"li"},Object(c.b)("li",{parentName:"ul"},"is in scope"),Object(c.b)("li",{parentName:"ul"},"is type coercible to object"),Object(c.b)("li",{parentName:"ul"},"contains property ",Object(c.b)("inlineCode",{parentName:"li"},"someProperty"))))),Object(c.b)("pre",null,Object(c.b)("code",Object(n.a)({parentName:"pre"},{className:"language-yaml"}),"$(someObject.someProperty)\n")),Object(c.b)("h3",{id:"coercion"},"Coercion"),Object(c.b)("p",null,"Object typed values are coercible to:"),Object(c.b)("ul",null,Object(c.b)("li",{parentName:"ul"},Object(c.b)("a",Object(n.a)({parentName:"li"},{href:"/docs/reference/opspec/types/boolean"}),"boolean")," (objects which are null or empty coerce to false; all else coerce to true)"),Object(c.b)("li",{parentName:"ul"},Object(c.b)("a",Object(n.a)({parentName:"li"},{href:"/docs/reference/opspec/types/file"}),"file")," (will be serialized to JSON)"),Object(c.b)("li",{parentName:"ul"},Object(c.b)("a",Object(n.a)({parentName:"li"},{href:"/docs/reference/opspec/types/string"}),"string")," (will be serialized to JSON)")),Object(c.b)("h4",{id:"coercion-example-object-to-string"},"Coercion Example (object to string)"),Object(c.b)("pre",null,Object(c.b)("code",Object(n.a)({parentName:"pre"},{className:"language-yaml"}),"name: objAsString\ninputs:\n  obj:\n    object:\n      default:\n        prop1: prop1Value\n        prop2: [ item1 ]\nrun:\n  container:\n    image: { ref: alpine }\n    cmd:\n    - echo\n    - $(obj)\n")))}s.isMDXComponent=!0},193:function(e,t,r){"use strict";r.d(t,"a",(function(){return s})),r.d(t,"b",(function(){return O}));var n=r(0),a=r.n(n);function c(e,t,r){return t in e?Object.defineProperty(e,t,{value:r,enumerable:!0,configurable:!0,writable:!0}):e[t]=r,e}function i(e,t){var r=Object.keys(e);if(Object.getOwnPropertySymbols){var n=Object.getOwnPropertySymbols(e);t&&(n=n.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),r.push.apply(r,n)}return r}function o(e){for(var t=1;t<arguments.length;t++){var r=null!=arguments[t]?arguments[t]:{};t%2?i(Object(r),!0).forEach((function(t){c(e,t,r[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(r)):i(Object(r)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(r,t))}))}return e}function p(e,t){if(null==e)return{};var r,n,a=function(e,t){if(null==e)return{};var r,n,a={},c=Object.keys(e);for(n=0;n<c.length;n++)r=c[n],t.indexOf(r)>=0||(a[r]=e[r]);return a}(e,t);if(Object.getOwnPropertySymbols){var c=Object.getOwnPropertySymbols(e);for(n=0;n<c.length;n++)r=c[n],t.indexOf(r)>=0||Object.prototype.propertyIsEnumerable.call(e,r)&&(a[r]=e[r])}return a}var l=a.a.createContext({}),b=function(e){var t=a.a.useContext(l),r=t;return e&&(r="function"==typeof e?e(t):o({},t,{},e)),r},s=function(e){var t=b(e.components);return(a.a.createElement(l.Provider,{value:t},e.children))},u="mdxType",j={inlineCode:"code",wrapper:function(e){var t=e.children;return a.a.createElement(a.a.Fragment,{},t)}},m=Object(n.forwardRef)((function(e,t){var r=e.components,n=e.mdxType,c=e.originalType,i=e.parentName,l=p(e,["components","mdxType","originalType","parentName"]),s=b(r),u=n,m=s["".concat(i,".").concat(u)]||s[u]||j[u]||c;return r?a.a.createElement(m,o({ref:t},l,{components:r})):a.a.createElement(m,o({ref:t},l))}));function O(e,t){var r=arguments,n=t&&t.mdxType;if("string"==typeof e||n){var c=r.length,i=new Array(c);i[0]=m;var o={};for(var p in t)hasOwnProperty.call(t,p)&&(o[p]=t[p]);o.originalType=e,o[u]="string"==typeof e?e:n,i[1]=o;for(var l=2;l<c;l++)i[l]=r[l];return a.a.createElement.apply(null,i)}return a.a.createElement.apply(null,r)}m.displayName="MDXCreateElement"}}]);