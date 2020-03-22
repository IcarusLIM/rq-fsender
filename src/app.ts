import Koa from 'koa'
import koaBody from 'koa-body'
import koaJson from 'koa-json'
import * as log4j from 'log4js'

import {requestLogger, errorHandler, allowCrossOrigin} from './middlewares'
import router from './router'
import {FILE_STORE, SERVER_PORT} from './config'

const logger = log4j.getLogger()
logger.level = 'debug'

const app = new Koa();

app.use(requestLogger)
app.use(errorHandler)
app.use(allowCrossOrigin)
app.use(koaJson())
app.use(koaBody({ multipart: true, formidable: {uploadDir: FILE_STORE}}))
app.use(router.routes())

app.listen(SERVER_PORT)
