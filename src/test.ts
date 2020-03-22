import redis from './redis'

const f = async () => {
    console.log(111)
    try {
        const res = await redis.lpop("q:r:www.baidu.com::http")
        redis.close()
        console.log(222)
        console.log(res)
    } catch (e) {
        console.log(e)
    }
}
f()