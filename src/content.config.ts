import { defineCollection } from "astro:content";
import { glob } from "astro/loaders";
import { z } from "astro/zod";

const blog = defineCollection({
  // Load Markdown and MDX files in the `src/content/blog/` directory.
  loader: glob({ base: "./src/content/blog", pattern: "**/*.{md,mdx}" }),
  // Type-check frontmatter using a schema
  schema: ({}) =>
    z.object({
      title: z.string(),
      date: z.coerce.date(),
      uri: z.string(),
      author: z.object({
        name: z.string(),
        email: z.email(),
      }),
      description: z.string(),
    }),
});

const history = defineCollection({
  loader: glob({ base: "./src/content/history", pattern: "**/*.{md,mdx}" }),
  schema: z.object({
    message: z.string(),
    date: z.coerce.date().optional(),
  }),
});

export const collections = { blog, history };
