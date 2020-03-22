import KoaRouter from 'koa-router'
import uploadController from './upload/controller'

const router = new KoaRouter();
router.prefix('/fsender');

router.post('/upload', uploadController.upload)

export default router