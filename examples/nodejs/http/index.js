'use nodent';
const Koa = require('koa');

const app = new Koa();

app.use(async ctx => {
  ctx.response.type='html';

  ctx.body = `
  <html>
    Hello <a href="http://opspec.io">opspec</a> World!
  </html>
  `;

});

app.listen(3000);
