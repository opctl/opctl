const koaApp = new require('koa')();
const koaRouter = require('koa-router');

koaApp.use(
  koaRouter()
    .get('/', async koaContext => {
      koaContext.response.body = 'Ok';
    })
);

return koaApp.listen(80);
