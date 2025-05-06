import App from './app';
import userRouter from './routes/user.route';
import namespaceRouter from './routes/namespace.route';
import resourceRouter from './routes/resource.route';
import roleRouter from './routes/role.route';
import rolePermissionRouter from './routes/permission.route';
import webhookRouter from './routes/webhook.route';
import ruleRouter from './routes/rule.route';

const app = new App(
    [
        new namespaceRouter(),
        new resourceRouter(),
        new roleRouter(),
        new userRouter(),
        new rolePermissionRouter(),
        new webhookRouter(),
        new ruleRouter()
    ]
);

app.listen();