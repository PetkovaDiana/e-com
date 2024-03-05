from django.core.files.base import ContentFile
from django.db import models
from django.utils.html import format_html
from django.forms.widgets import ClearableFileInput
from django.utils.safestring import mark_safe
import base64


class RequestCall(models.Model):
    phone = models.CharField(verbose_name="Телефон", max_length=25)
    name = models.CharField(verbose_name="Имя", blank=True, null=True, max_length=100)
    created_at = models.DateTimeField(verbose_name="Дата создания")
    user = models.ForeignKey('auth_app.User', verbose_name='Пользователь', null=True, blank=True, on_delete=models.SET_NULL)
    message = models.TextField(verbose_name="Текст сообщения")

    class Meta:
        managed = False
        db_table = 'request_call'
        verbose_name = 'Заявка на обратный звонок'
        verbose_name_plural = 'Заявки на обратные звонки'

    def __str__(self):
        return f'Заявка на обратный звонок - {self.phone}'


class ResponsesVacancy(models.Model):
    surname = models.CharField(verbose_name="Фамилия", max_length=30)
    name = models.CharField(verbose_name="Имя", max_length=30)
    lastname = models.CharField(verbose_name="Отчество", max_length=30, blank=True, null=True)
    phone = models.CharField(verbose_name="Телефон", max_length=30)
    email = models.EmailField(verbose_name="Электронная почта", max_length=30, blank=True, null=True)
    created_at = models.DateTimeField(verbose_name="Дата создания")
    comment = models.TextField(verbose_name="Комментарий", blank=True, null=True)
    vacancy = models.ForeignKey(
        'cms.Vacancies', verbose_name='Вакансия', null=True, blank=True, on_delete=models.SET_NULL)

    class Meta:
        managed = False
        db_table = 'request_vacancy'
        verbose_name = 'Отклик на вакансию'
        verbose_name_plural = 'Отклики на вакансии'

    def __str__(self):
        return f'Отклик на вакансию - {self.vacancy.title}'


class Review(models.Model):
    body = models.TextField(verbose_name="Текст", blank=True, null=True)
    rating = models.BigIntegerField(verbose_name="Рейтинг", blank=True, null=True)
    product_uuid = models.ForeignKey('catalog.Product', models.DO_NOTHING, db_column='product_uuid', verbose_name="Товар")
    user = models.ForeignKey('auth_app.User', models.DO_NOTHING, verbose_name="Юзер")
    recommend = models.BooleanField(blank=True, null=True, verbose_name="Рекомендовали бы данный товар")
    created_at = models.DateTimeField(verbose_name="Дата создания")

    class Meta:
        managed = False
        db_table = 'review'
        verbose_name = 'Отзыв'
        verbose_name_plural = 'Отзывы'

    def __str__(self):
        return f'Отзыв - {self.product_uuid.title}'



class ReviewPhotos(models.Model):
    id = models.BigAutoField(primary_key=True)
    image = models.BinaryField(blank=True, null=True)
    review = models.ForeignKey(Review, models.DO_NOTHING, blank=True, null=True)

    class Meta:
        managed = False
        db_table = 'review_photos'
        verbose_name = 'Фотографии отзывов'
        verbose_name_plural = 'Фотография отзыва'

    def image_tag(self):
        # bytes_str = base64.b64decode(b64_str.encode('ascii'))
        image_base64 = base64.b64encode(self.image).decode('utf-8')
        return mark_safe(f'<img src="data:image/png;base64,{image_base64}" width="150px"/>')
    image_tag.short_description = 'Фотография'

    def __str__(self):
        return f'Фото - {self.id}'