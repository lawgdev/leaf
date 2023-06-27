import { z } from "zod";

const env_schema = z.object({
  API_URL_V1: z.string().default("http://100.105.87.12:8080/v1"),
  TWIG_URL: z.string().default("ws://100.105.87.12:4000"),
});

export const env = env_schema.parse(process.env);
