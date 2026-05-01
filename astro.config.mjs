// @ts-check

import mdx from "@astrojs/mdx";
import sitemap from "@astrojs/sitemap";
import { defineConfig, fontProviders } from "astro/config";
import tailwindcss from "@tailwindcss/vite";

// https://astro.build/config
export default defineConfig({
  integrations: [mdx(), sitemap()],

  fonts: [
    {
      provider: fontProviders.local(),
      name: "Iosevka Comfy",
      cssVariable: "--font-iosevka",
      fallbacks: ["sans-serif"],
      options: {
        variants: [
          {
            weight: 400,
            style: "normal",
            display: "swap",
            src: [
              "./src/assets/fonts/iosevka-comfy-fixed-normalregularupright.woff2",
            ],
          },
          {
            weight: 700,
            style: "normal",
            display: "swap",
            src: [
              "./src/assets/fonts/iosevka-comfy-fixed-normalboldupright.woff2",
            ],
          },
          {
            weight: 800,
            style: "normal",
            display: "swap",
            src: [
              "./src/assets/fonts/iosevka-comfy-fixed-normalextraboldupright.woff2",
            ],
          },
        ],
      },
    },
  ],

  vite: {
    plugins: [tailwindcss()],
  },
});
