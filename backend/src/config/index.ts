import { config } from 'dotenv';
config({ path: `.env.${process.env.NODE_ENV || 'development'}` });

console.log("current environment is", process.env.NODE_ENV, process.env.DOMAIN)

export const CREDENTIALS = process.env.CREDENTIALS === 'true';
export const { 
    NODE_ENV, PORT, DOMAIN, DB_HOST, DB_PORT, DB_USER="root", 
    DB_PASSWORD, DB_DATABASE="ucode_review", 
    LOG_FORMAT, LOG_DIR, ORIGIN,
    AI_MODEL, AI_API, AI_KEY
} = process.env;