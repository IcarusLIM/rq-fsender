import Redis from 'redis'
import { REDIS_CONF } from './config'

import { promisify } from 'util'


const client = Redis.createClient({ host: REDIS_CONF.host, port: REDIS_CONF.port, password: REDIS_CONF.password })

export default {
    lpush: promisify(client.lpush).bind(client),
    lpop: promisify(client.lpop).bind(client),
    sadd: promisify(client.sadd).bind(client),
    spop: promisify(client.spop).bind(client),
    close: () => {
        client.end(true)
    }
}