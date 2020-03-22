import Koa from 'koa'
import * as log4j from 'log4js'

const logger = log4j.getLogger('middlewares')

export const errorHandler = async (ctx: Koa.Context, next: Koa.Next) => {
    try {
        await next();
    } catch (e) {
        logger.warn(e)
        ctx.status = e.statusCode || e.status || 500;
        ctx.body = { res: false, err: e.message }
    }
}

export const requestLogger = async (ctx: Koa.Context, next: Koa.Next) => {
    const req = ctx.request
    logger.info(`${req.method} ${req.url}`)
    await next()
    logger.info(`${req.method} ${req.url} ${ctx.status}`)
}

function setCrosHeader(ctx: Koa.Context) {
    ctx.set('Access-Control-Allow-Origin', '*')
    ctx.set('Access-Control-Allow-Headers', '*')
    ctx.set('Access-Control-Allow-Methods', '*')
}

export const allowCrossOrigin = async (ctx: Koa.Context, next: Koa.Next) => {
    if (ctx.request.method === 'OPTIONS') {
        logger.debug('Options req')
        setCrosHeader(ctx)
        ctx.status = 204
    } else {
        await next()
        setCrosHeader(ctx)
    }
}