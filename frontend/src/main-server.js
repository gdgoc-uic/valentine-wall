import { renderToString } from 'vue/server-renderer';
import { createApp } from "./main";
import { basename } from 'path';
import { notiwindSSRShim } from './notify';
import { renderHeadToString } from '@vueuse/head';

export async function render(url, manifest) {
  const { app, router, head, store } = createApp();
  app.use(notiwindSSRShim());
  await router.push(url);
  const to = router.currentRoute.value;
  const requiresAuth = to.matched.some(record => record.meta.requiresAuth);
  if (requiresAuth) {
    await router.replace('/');
  }
  await router.isReady();
  const ctx = {};
  const appHtml = await renderToString(app, ctx);
  const appHead = renderHeadToString(head);
  const preloadLinks = renderPreloadLinks(ctx.modules, manifest);
  const state = JSON.stringify(store.state);
  return [appHtml, preloadLinks, appHead, state];
}

function renderPreloadLinks(modules, manifest) {
  let links = ''
  const seen = new Set()
  modules.forEach((id) => {
    const files = manifest[id]
    if (files) {
      files.forEach((file) => {
        if (!seen.has(file)) {
          seen.add(file)
          const filename = basename(file)
          if (manifest[filename]) {
            for (const depFile of manifest[filename]) {
              links += renderPreloadLink(depFile)
              seen.add(depFile)
            }
          }
          links += renderPreloadLink(file)
        }
      })
    }
  })
  return links
}

function renderPreloadLink(file) {
  if (file.endsWith('.js')) {
    return `<link rel="modulepreload" crossorigin href="${file}">`
  } else if (file.endsWith('.css')) {
    return `<link rel="stylesheet" href="${file}">`
  } else if (file.endsWith('.woff')) {
    return ` <link rel="preload" href="${file}" as="font" type="font/woff" crossorigin>`
  } else if (file.endsWith('.woff2')) {
    return ` <link rel="preload" href="${file}" as="font" type="font/woff2" crossorigin>`
  } else if (file.endsWith('.gif')) {
    return ` <link rel="preload" href="${file}" as="image" type="image/gif">`
  } else if (file.endsWith('.jpg') || file.endsWith('.jpeg')) {
    return ` <link rel="preload" href="${file}" as="image" type="image/jpeg">`
  } else if (file.endsWith('.png')) {
    return ` <link rel="preload" href="${file}" as="image" type="image/png">`
  } else {
    // TODO
    return ''
  }
}