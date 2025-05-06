# monto-acl-frontend

## 开发注意

- node: v18+
- 需使用 pnpm

## 账号

- admin
- heng.du

## 笔记

1. 本地调试 mobx 内部代理对象，可以这样：

```js
import { toJS } from 'mobx';

console.log(toJS(user));
```
