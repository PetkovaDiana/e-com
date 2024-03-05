from fastapi import APIRouter
from jinja2 import Environment, FileSystemLoader
import pdfkit

from .send_email import send_email
from .api_models import AnswerAPI, OrderAPI


router = APIRouter()
env = Environment(loader=FileSystemLoader('.'))
template = env.get_template("app/index.html")


@router.post("/create_payment_invoice/", response_model=AnswerAPI)
async def get_search(
        order: OrderAPI
):
    '''
    Создание счета на безналичную оплату для юр лич
    '''

    html_template = template.render({'order': order})
    path = 'app/payment_invoice/'
    url = f'https://ufaelectro.ru/media/payment_invoice/order_{order.number}.pdf'
    pdfkit.from_string(html_template, path + f'order_{order.number}.pdf')
    file_name = f'Cчет на оплату заказа №{order.number} от {order.date}г. БЭСМ'
    await send_email(order.customer.email, file_name, path + f'order_{order.number}.pdf')

    return {'status': 200, 'url': url}


