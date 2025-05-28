# Common build stage
FROM node:18-alpine as common-build-stage

COPY . ./app

WORKDIR /app

RUN npm i -g vite && npm i -g pnpm

# 接收传入的构建参数
ARG VITE_APP_APIHOST

# 生成 .env.production 文件
RUN echo "VITE_APP_APIHOST=$VITE_APP_APIHOST" >> .env.production

RUN pnpm install && pnpm build

RUN ls -al

# EXPOSE 3000

FROM nginx:1.19.1 as production-build-stage

COPY ./nginx.conf /etc/nginx/conf.d/default.conf

COPY --from=common-build-stage /app/dist /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]