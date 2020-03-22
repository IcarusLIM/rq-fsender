import Koa from 'koa'
import * as log4j from 'log4js'
import redis from '../redis'

const logger = log4j.getLogger('upload.controller')

export default {
    upload: async (ctx: Koa.Context, next: Koa.Next) => {
        const files = ctx.request.files
        if (!files) {
            throw new Error('No file in request form')
        }
        const file = files.file
        logger.info(file)
        const filePath = file.path
        const fid = filePath.substr(filePath.lastIndexOf('/') + 1, filePath.length)
        const file_info = {
            size: file.size,
            fid: fid,
            path: filePath,
            name: file.name,
            type: file.type
        }
        // @ts-ignore
        await redis.sadd('os:rq:fsender:files', JSON.stringify(file_info))
        ctx.body = { res: true, file_info: file_info }
    }
}