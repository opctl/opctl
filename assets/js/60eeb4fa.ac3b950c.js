"use strict";(self.webpackChunkopctl=self.webpackChunkopctl||[]).push([[4539],{3905:function(e,t,n){n.d(t,{Zo:function(){return p},kt:function(){return m}});var r=n(7294);function a(e,t,n){return t in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function o(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);t&&(r=r.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),n.push.apply(n,r)}return n}function i(e){for(var t=1;t<arguments.length;t++){var n=null!=arguments[t]?arguments[t]:{};t%2?o(Object(n),!0).forEach((function(t){a(e,t,n[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):o(Object(n)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))}))}return e}function l(e,t){if(null==e)return{};var n,r,a=function(e,t){if(null==e)return{};var n,r,a={},o=Object.keys(e);for(r=0;r<o.length;r++)n=o[r],t.indexOf(n)>=0||(a[n]=e[n]);return a}(e,t);if(Object.getOwnPropertySymbols){var o=Object.getOwnPropertySymbols(e);for(r=0;r<o.length;r++)n=o[r],t.indexOf(n)>=0||Object.prototype.propertyIsEnumerable.call(e,n)&&(a[n]=e[n])}return a}var s=r.createContext({}),c=function(e){var t=r.useContext(s),n=t;return e&&(n="function"==typeof e?e(t):i(i({},t),e)),n},p=function(e){var t=c(e.components);return r.createElement(s.Provider,{value:t},e.children)},u={inlineCode:"code",wrapper:function(e){var t=e.children;return r.createElement(r.Fragment,{},t)}},d=r.forwardRef((function(e,t){var n=e.components,a=e.mdxType,o=e.originalType,s=e.parentName,p=l(e,["components","mdxType","originalType","parentName"]),d=c(n),m=a,h=d["".concat(s,".").concat(m)]||d[m]||u[m]||o;return n?r.createElement(h,i(i({ref:t},p),{},{components:n})):r.createElement(h,i({ref:t},p))}));function m(e,t){var n=arguments,a=t&&t.mdxType;if("string"==typeof e||a){var o=n.length,i=new Array(o);i[0]=d;var l={};for(var s in t)hasOwnProperty.call(t,s)&&(l[s]=t[s]);l.originalType=e,l.mdxType="string"==typeof e?e:a,i[1]=l;for(var c=2;c<o;c++)i[c]=n[c];return r.createElement.apply(null,i)}return r.createElement.apply(null,n)}d.displayName="MDXCreateElement"},5937:function(e,t,n){n.r(t),n.d(t,{assets:function(){return p},contentTitle:function(){return s},default:function(){return m},frontMatter:function(){return l},metadata:function(){return c},toc:function(){return u}});var r=n(3117),a=n(102),o=(n(7294),n(3905)),i=["components"],l={title:"Run a go service",sidebar_label:"run a go service"},s=void 0,c={unversionedId:"run-a-go-service",id:"run-a-go-service",title:"Run a go service",description:"We'll now look at an op to run a sample Go application. Please see the code at https://github.com/opctl/opctl/tree/main/examples/run-a-go-service",source:"@site/docs/run-a-go-service.md",sourceDirName:".",slug:"/run-a-go-service",permalink:"/docs/run-a-go-service",draft:!1,editUrl:"https://github.com/opctl/opctl/edit/main/website/docs/run-a-go-service.md",tags:[],version:"current",lastUpdatedBy:"=",lastUpdatedAt:1682904716,formattedLastUpdatedAt:"May 1, 2023",frontMatter:{title:"Run a go service",sidebar_label:"run a go service"}},p={},u=[{value:"mysql",id:"mysql",level:4},{value:"dev",id:"dev",level:4}],d={toc:u};function m(e){var t=e.components,n=(0,a.Z)(e,i);return(0,o.kt)("wrapper",(0,r.Z)({},d,n,{components:t,mdxType:"MDXLayout"}),(0,o.kt)("p",null,"We'll now look at an op to run a sample Go application. Please see the code at ",(0,o.kt)("a",{parentName:"p",href:"https://github.com/opctl/opctl/tree/main/examples/run-a-go-service"},"https://github.com/opctl/opctl/tree/main/examples/run-a-go-service"),"\nThe sample application we have requires a mysql database, therefore to run it locally we need to:"),(0,o.kt)("ol",null,(0,o.kt)("li",{parentName:"ol"},"Start a MySQL database with some DB credentials"),(0,o.kt)("li",{parentName:"ol"},"Seed the DB with sample data"),(0,o.kt)("li",{parentName:"ol"},"Start the application and provide the MySQL DB credentials as inputs - our application will read those from environment variables")),(0,o.kt)("p",null,"The ops to make that happen are explained below."),(0,o.kt)("h4",{id:"mysql"},"mysql"),(0,o.kt)("pre",null,(0,o.kt)("code",{parentName:"pre",className:"language-yaml"},"name: mysql\ndescription: runs a mysql container, seeding it with sample data\ninputs:\n  MYSQL_USER:\n    description: username for MySQL user to create\n    string:\n      default: testuser\n  MYSQL_PASSWORD:\n    description: password for MySQL user to create\n    string:\n      default: testpassword\n  MYSQL_DATABASE:\n    description: name of mysql database to create\n    string:\n      default: testapp\n  MYSQL_HOST:\n    string:\n      default: run-a-go-service-mysql\nrun:\n  container:\n    dirs:\n      # mount our data seed scripts\n      /docker-entrypoint-initdb.d:\n    image:\n      ref: 'mysql:8'\n    # this sets an opctl overlay network DNS A record so the containers available via this name.\n    name: run-a-go-service-mysql\n    envVars:\n      MYSQL_USER:\n      MYSQL_PASSWORD:\n      MYSQL_DATABASE:\n      MYSQL_RANDOM_ROOT_PASSWORD: \"yes\"\n    ports:\n      3306: 3306\n")),(0,o.kt)("p",null,"This op will run a MySQL database in a Docker container and seeds the MySQL database with a single table and a single row of data."),(0,o.kt)("p",null,"Note that instead of writing the ",(0,o.kt)("inlineCode",{parentName:"p"},"cmd")," inline, we could have also put the run script in a file in the op directory and referenced it.\ne.g."),(0,o.kt)("pre",null,(0,o.kt)("code",{parentName:"pre",className:"language-yaml"},'cmd: ["/run.sh"]\n')),(0,o.kt)("p",null,"This becomes useful for readability if the script we're running in a container is large."),(0,o.kt)("p",null,"Note that while we could have also built a custom docker image that contains that script and ran it, that's not a recommended practice. While the op remains portable with that approach, it is less transparent and harder to maintain, since we would have the extra step of building and pushing that image whenever we make a change to the script. Instead, we usually opt for the leanest Docker image that fulfills our use case, and include any scripts or binaries we need in the op directory."),(0,o.kt)("h4",{id:"dev"},"dev"),(0,o.kt)("pre",null,(0,o.kt)("code",{parentName:"pre",className:"language-yaml"},"name: dev\ndescription: runs go-svc for development\ninputs:\n  MYSQL_USER:\n    description: username for MySQL user to create\n    string:\n      default: testuser\n  MYSQL_PASSWORD:\n    description: password for MySQL user to create\n    string:\n      default: testpassword\n  MYSQL_DATABASE:\n    description: Database to create\n    string:\n      default: testapp\n  MYSQL_HOST:\n    string:\n      default: mysql # the host for mysql is the container name in the mysql op\nrun:\n  parallel:\n    - op:\n        ref: ../mysql # we reference the mysql op using its relative path to this op\n        inputs:\n          # we pass the relevant inputs through to the mysql op\n          MYSQL_USER:\n          MYSQL_PASSWORD:\n          MYSQL_HOST:\n          MYSQL_DATABASE:\n    - container:\n        image:\n          ref: 'golang:1.20'\n        name: go-svc\n        dirs:\n          # mount the source code of our app to the container\n          /src: $(../..) \n        envVars:\n          # let our code know how to connect to the DB\n          MYSQL_USER:\n          MYSQL_PASSWORD:\n          MYSQL_HOST:\n          MYSQL_DATABASE:\n        workDir: /src\n        ports:\n          8080: 8080\n        cmd:\n          - sh\n          - -ce\n          - |\n            go get -u github.com/cespare/reflex # get reflex to watch and hot-reload our main.go\n            sleep 30 # we'll sleep while the MySQL DB starts\n            /go/bin/reflex -g 'main.go' -s -- sh -c 'go run main.go'\n")),(0,o.kt)("p",null,"This op will do the following, in parallel:"),(0,o.kt)("ol",null,(0,o.kt)("li",{parentName:"ol"},"Run the mysql op, passing in the inputs needed to create it"),(0,o.kt)("li",{parentName:"ol"},"run our service in the ",(0,o.kt)("inlineCode",{parentName:"li"},"golang")," container. We'll use ",(0,o.kt)("a",{parentName:"li",href:"https://github.com/cespare/reflex"},"reflex")," to make development easier and avoid having to restart the whole operation after every change to main.go")),(0,o.kt)("p",null,"Now to run the service locally, we only run ",(0,o.kt)("inlineCode",{parentName:"p"},"opctl run dev")),(0,o.kt)("p",null,"Notice how we can run ops in a ",(0,o.kt)("inlineCode",{parentName:"p"},"parallel")," block. Our app runs parallel to the MySQL DB container (via the ",(0,o.kt)("inlineCode",{parentName:"p"},"mysql")," op), because we need mysql running while our app runs."))}m.isMDXComponent=!0}}]);