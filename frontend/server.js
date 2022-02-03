//@ts-check
const fs = require('fs');
const path = require('path');
const polka = require('polka');
const send = require('@polka/send-type');

const { createServer: createViteServer } = require('vite');
const isTest = process.env.NODE_ENV === 'test' || !!process.env.VITE_TEST_BUILD;
const resolve = (/** @type {string} */ p) => path.resolve(__dirname, p);

async function createServer(
  root = process.cwd(),
  isProd = process.env.NODE_ENV === 'production'
) {
  const indexProd = isProd
    ? fs.readFileSync(resolve('dist/client/index.html'), 'utf-8')
    : '';

  const manifest = isProd
    ? // @ts-ignore
    require('./dist/client/ssr-manifest.json')
    : {};

  /**
   * @type {import('vite').ViteDevServer}
   */
  let vite;

  const app = polka();
  if (!isProd) {
    vite = await createViteServer({
      root,
      logLevel: isTest ? 'error' : 'info',
      server: {
        middlewareMode: 'ssr',
        watch: {
          usePolling: true,
          interval: 100
        }
      }
    });

    app.use(vite.middlewares);
  } else {
    const compression = require('compression')();
    const sirv = require('sirv/build');
    const morgan = require('morgan');
    /** @type {import('sirv').Options}  */
    const assetOptions = { extensions: [], ignores: ['index.html'], single: true };
    const assets = sirv(resolve('dist/client'), assetOptions);
    app.use(compression);
    app.use(assets);
    app.use(morgan('common'));
  }

  let template, render;
  if (isProd) {
    const mod = await import('./dist/server/main-server.js');
    template = indexProd;
    render = mod.render;
  }

  app.get('*', async (req, res) => {
    try {
      const url = req.originalUrl;
      if (!isProd) {
        // always read fresh template in dev
        template = await vite.transformIndexHtml(url, fs.readFileSync(resolve('index.html'), 'utf-8'));
        const mod = await vite.ssrLoadModule('/src/main-server.js');
        render = mod.render;
      }

      const [appHtml, preloadLinks, appHead] = await render(url, manifest);
      const { htmlAttrs, headTags, bodyAttrs } = appHead;
      const html = template
        .replace(`data-html-attrs=""`, htmlAttrs)
        .replace(`<!--head-tags-->`, headTags)
        .replace(`data-body-attrs=""`, bodyAttrs)
        .replace(`<!--preload-links-->`, preloadLinks)
        .replace(`<!--app-html-->`, appHtml);

      send(res, 200, html, { 'Content-Type': 'text/html' });
    } catch (e) {
      vite && vite.ssrFixStacktrace(e);
      console.log(e.stack);
      send(res, 500, e.stack);
    }
  });
  return { app, vite };
}

if (!isTest) {
  createServer().then(({ app }) =>
    app.listen(3000, () => {
      console.log('Serving on: http://localhost:3000');
    })
  );
}