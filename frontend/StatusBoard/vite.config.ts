import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'
import { statusApi } from './server/statusApi'

// Runs on the host, not in Docker: npm run dev -> http://localhost:5180
export default defineConfig({
  plugins: [react(), tailwindcss(), statusApi()],
  server: { port: 5180, strictPort: true },
  preview: { port: 5180, strictPort: true },
})
