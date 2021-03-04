import { defineConfig } from 'vite';
import reactRefresh from '@vitejs/plugin-react-refresh';

export default defineConfig({
  build: {
    outDir: '../src/html',
  },
  server: {
    proxy: {
      '/api': 'http://localhost:8080',
    },
  },
  plugins: [reactRefresh()],
});
