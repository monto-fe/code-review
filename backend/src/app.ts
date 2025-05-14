import express from 'express';
import swaggerJSDoc from 'swagger-jsdoc';
import swaggerUi from 'swagger-ui-express';
import cors from 'cors';
import fs from 'fs';

import DB from './databases';
import { authenticateJwt } from './middlewares/auth';
import errorHandler from './middlewares/errorHandler';
import { PORT, DOMAIN } from './config';
import { Routes } from './interfaces/routes.interface';
// 初始化ai信息
import aiConfigManager from './services/aiConfigManager';
// 初始化gitlab信息
import { GitlabManagerService } from './services/gitlabManager.service';
 
const corsOptions = {
	origin: "*",
	methods: "GET,HEAD,PUT,PATCH,POST,DELETE",
	optionsSuccessStatus: 204 // some legacy browsers (IE11, various SmartTVs) choke on 204
}

class App {
	app: express.Application;
	port = 3000;
	constructor(routers: Routes[]) {
		this.port = Number(PORT) || this.port;
		this.app = express();
		this.connectToDatabase();
		this.initializeSwagger();
		this.initializeMiddlewares();
		this.initializeRoutes(routers);
	}

	private connectToDatabase() {
		DB.sequelize.sync();
	}

	private initializeMiddlewares() {
		this.app.use(cors(corsOptions));
		
		this.app.use(express.json({limit: '50mb'}));
		this.app.use(express.urlencoded({ extended: true }));
		this.app.use(authenticateJwt);
		this.app.use(errorHandler)
	}

	private initializeRoutes(routes: Routes[]) {
		routes.forEach(route => {
			this.app.use('/v1', route.router);
		});
	}

	private initializeSwagger() {
		const options = {
			swaggerDefinition: {
				openapi: '3.0.0',
				info: {
					title: 'api document',
					version: '1.0.0',
					description:
						'swagger guide：https://swagger.io/docs/specification/basic-structure/',
				},
				servers: [{
					url: DOMAIN,
					description: 'local environment'
				}]
			},
			apis: ['./swagger/*.yaml'],
		};

		const specs = swaggerJSDoc(options);
		// 将 Swagger JSON 对象写入文件
		fs.writeFileSync('./swagger-output.json', JSON.stringify(specs, null, 2), 'utf-8');
		this.app.use('/api-docs', swaggerUi.serve, swaggerUi.setup(specs));
		this.app.get('/swagger.json', (req, res) => {
		res.setHeader('Content-Type', 'application/json');
			res.send(specs);
		});
	}
	public async listen() {
		try {
			await aiConfigManager.loadConfig();
			const gitlabManager = await GitlabManagerService.init();
			const cache = gitlabManager.getCache();
			console.log('GitLab Token 缓存初始化完成:', cache);
			this.app.listen(this.port, () => {
				console.log(`TypeScript with Express
				  http://localhost:${this.port}/`);
			});
		} catch (error) {
			console.error('Error loading AIConfig:', error);
		}
	}
}

export default App;