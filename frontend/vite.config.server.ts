import { defineConfig } from 'vite';
import config from './vite.config';

export default Object.assign(config, defineConfig({
    ssr: {
        format: 'cjs'   
    }
}));