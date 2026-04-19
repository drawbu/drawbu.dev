import { getCollection } from "astro:content";
import rss from "@astrojs/rss";

export async function GET(context) {
  const posts = await getCollection("blog");
  return rss({
    title: "drawbu.dev blog",
    site: "https://drawbu.dev",
    description: "software engineering student talking about cool stuff",
    items: posts.map((post) => ({
      title: post.data.title,
      link: `/blog/${post.id}/`,
      description: post.data.description,
      pubDate: post.data.date,
    })),
  });
}
