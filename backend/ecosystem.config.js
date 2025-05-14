const dotenv = require('dotenv');
const fs = require('fs');

// 获取当前指定的环境
const currentEnv = process.env.NODE_ENV || 'development';
const envFile = `.env.${currentEnv}`;

// 检查文件存在才加载
if (fs.existsSync(envFile)) {
  dotenv.config({ path: envFile });
  console.log(`✅ Loaded env from ${envFile}`);
} else {
  console.warn(`⚠️ No ${envFile} file found`);
}

console.log("process.env.PORT", process.env.PORT)
/**
 * @description pm2 configuration file.
 * @example
 *  production mode :: pm2 start ecosystem.config.js --only prod
 *  development mode :: pm2 start ecosystem.config.js --only dev
 */
 module.exports = {
  apps: [
    {
      name: 'ucode_review', // pm2 start App name
      script: 'dist/index.js',
      autorestart: true, // auto restart if process crash
      watch: false, // files change automatic restart
      ignore_watch: ['node_modules', 'logs'], // ignore files change
      max_memory_restart: '1G', // restart if process use more than 1G memory
      merge_logs: true, // if true, stdout and stderr will be merged and sent to pm2 log
      output: './logs/access.log', // pm2 log file
      error: './logs/error.log', // pm2 error log file
      env_development: {
        PORT: 9000,
        NODE_ENV: 'development',
        DOMAIN: 'http://lobalhost:9000',
        DB_HOST: "localhost",
        DB_PORT: 3306,
        DB_USER: "mysql",
        DB_PASSWORD: "mysql123456",
        DB_DATABASE: "ucode_review"
      },
      env_production: {
        PORT: 9000,
        NODE_ENV: 'production',
        DOMAIN: 'http://localhost:9000',
        DB_HOST: "mariadb",
        DB_PORT: 3306,
        DB_USER: "root",
        DB_PASSWORD: "mysql123456",
        DB_DATABASE: "ucode_review"
      }
    }
  ]
};
